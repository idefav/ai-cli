package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
)

type CurlOptions struct {
	Data        string
	Fail        bool
	Include     bool
	Output      string
	RemoteName  bool
	Silent      bool
	UploadFile  string
	User        string
	UserAgent   string
	Verbose     bool
	AISummarize bool
}

func parseCurlArgs(args []string) (string, *CurlOptions, error) {
	options := &CurlOptions{}
	var url string

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "-d", "--data":
			i++
			options.Data = args[i]
		case "-f", "--fail":
			options.Fail = true
		case "-i", "--include":
			options.Include = true
		case "-o", "--output":
			i++
			options.Output = args[i]
		case "-O", "--remote-name":
			options.RemoteName = true
		case "-s", "--silent":
			options.Silent = true
		case "-T", "--upload-file":
			i++
			options.UploadFile = args[i]
		case "-u", "--user":
			i++
			options.User = args[i]
		case "-A", "--user-agent":
			i++
			options.UserAgent = args[i]
		case "-v", "--verbose":
			options.Verbose = true
		case "--ai":
			options.AISummarize = true
		default:
			if strings.HasPrefix(arg, "-") {
				return "", nil, fmt.Errorf("unknown option: %s", arg)
			}
			url = arg
		}
	}

	if url == "" {
		return "", nil, fmt.Errorf("URL is required")
	}

	return url, options, nil
}

func HandleCurl(prompt string, processQuery func(string, bool)) {
	args := strings.Fields(prompt)[1:] // Skip "curl"
	url, options, err := parseCurlArgs(args)
	if err != nil {
		fmt.Printf("Error parsing arguments: %v\n", err)
		return
	}

	// Add http:// if missing
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}

	var req *http.Request
	var resp *http.Response

	if options.UploadFile != "" {
		file, err := os.Open(options.UploadFile)
		if err != nil {
			fmt.Printf("Error opening file: %v\n", err)
			return
		}
		defer file.Close()

		req, err = http.NewRequest("PUT", url, file)
		if err != nil {
			fmt.Printf("Error creating request: %v\n", err)
			return
		}
	} else if options.Data != "" {
		req, err = http.NewRequest("POST", url, bytes.NewBufferString(options.Data))
		if err != nil {
			fmt.Printf("Error creating request: %v\n", err)
			return
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		req, err = http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Printf("Error creating request: %v\n", err)
			return
		}
	}

	if options.User != "" {
		req.SetBasicAuth(strings.Split(options.User, ":")[0], strings.Split(options.User, ":")[1])
	}

	if options.UserAgent != "" {
		req.Header.Set("User-Agent", options.UserAgent)
	}

	// Create context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	defer signal.Stop(sigChan)

	// Run request in goroutine
	type result struct {
		resp *http.Response
		err  error
	}
	resultChan := make(chan result, 1)

	go func() {
		client := &http.Client{}
		resp, err := client.Do(req.WithContext(ctx))
		resultChan <- result{resp, err}
	}()

	// Wait for either request completion or interrupt
	select {
	case <-sigChan:
		cancel()
		fmt.Println("\n请求已取消")
		return
	case res := <-resultChan:
		if res.err != nil {
			fmt.Printf("Request failed: %v\n", res.err)
			return
		}
		resp = res.resp
		defer resp.Body.Close()
	}

	if options.Fail && resp.StatusCode >= 400 {
		fmt.Printf("Request failed with status: %d\n", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return
	}

	output := ""
	if options.Include {
		output += fmt.Sprintf("HTTP/1.1 %d %s\n", resp.StatusCode, http.StatusText(resp.StatusCode))
		for k, v := range resp.Header {
			output += fmt.Sprintf("%s: %s\n", k, strings.Join(v, ", "))
		}
		output += "\n"
	}
	output += string(body)

	if options.Output != "" {
		err := os.WriteFile(options.Output, []byte(output), 0644)
		if err != nil {
			fmt.Printf("Error writing to file: %v\n", err)
			return
		}
	} else if options.RemoteName {
		filename := "index.html" // Default if can't determine from URL
		parts := strings.Split(url, "/")
		if len(parts) > 0 {
			last := parts[len(parts)-1]
			if last != "" {
				filename = last
			}
		}
		err := os.WriteFile(filename, []byte(output), 0644)
		if err != nil {
			fmt.Printf("Error writing to file: %v\n", err)
			return
		}
	} else if !options.Silent {
		fmt.Println(output)
	}

	if options.AISummarize {
		summaryPrompt := fmt.Sprintf("请总结以下内容:\n%s\n用中文简洁概括主要内容", string(body))
		fmt.Println("AI总结:")
		processQuery(summaryPrompt, true)
	}
}
