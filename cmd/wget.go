package cmd

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type WgetOptions struct {
	OutputDocument string
	Continue       bool
	Quiet          bool
	Verbose        bool
	Timeout        time.Duration
	UserAgent      string
	AISummarize    bool
}

func parseWgetArgs(args []string) (*WgetOptions, []string, error) {
	options := &WgetOptions{
		Timeout:   30 * time.Second,
		UserAgent: "Wget/1.21.4",
	}
	var urls []string

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "-O", "--output-document":
			i++
			options.OutputDocument = args[i]
		case "-c", "--continue":
			options.Continue = true
		case "-q", "--quiet":
			options.Quiet = true
		case "-v", "--verbose":
			options.Verbose = true
		case "-T", "--timeout":
			i++
			duration, err := time.ParseDuration(args[i] + "s")
			if err != nil {
				return nil, nil, fmt.Errorf("invalid timeout: %v", err)
			}
			options.Timeout = duration
		case "-U", "--user-agent":
			i++
			options.UserAgent = args[i]
		case "--ai":
			options.AISummarize = true
		default:
			if strings.HasPrefix(arg, "-") {
				return nil, nil, fmt.Errorf("unknown option: %s", arg)
			}
			urls = append(urls, arg)
		}
	}

	if len(urls) == 0 {
		return nil, nil, fmt.Errorf("no URLs specified")
	}

	return options, urls, nil
}

func HandleWget(prompt string, processQuery func(string, bool)) {
	args := strings.Fields(prompt)[1:] // Skip "wget"
	options, urls, err := parseWgetArgs(args)
	if err != nil {
		fmt.Printf("Error parsing arguments: %v\n", err)
		return
	}

	client := &http.Client{
		Timeout: options.Timeout,
	}

	for _, urlStr := range urls {
		startTime := time.Now()
		if !options.Quiet {
			fmt.Printf("--%s--  %s\n", startTime.Format("2006-01-02 15:04:05"), urlStr)
		}

		req, err := http.NewRequest("GET", urlStr, nil)
		if err != nil {
			fmt.Printf("创建请求失败: %v\n", err)
			continue
		}
		req.Header.Set("User-Agent", options.UserAgent)

		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("下载失败: %v\n", err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("HTTP错误: %s\n", resp.Status)
			continue
		}

		outputPath := options.OutputDocument
		if outputPath == "" {
			if u, err := url.Parse(urlStr); err == nil {
				outputPath = filepath.Base(u.Path)
				if outputPath == "" || outputPath == "/" || strings.HasSuffix(outputPath, "/") {
					outputPath = "index.html"
				}
				// Handle URLs without file extensions
				if !strings.Contains(outputPath, ".") {
					contentType := resp.Header.Get("Content-Type")
					if strings.Contains(contentType, "text/html") {
						outputPath += ".html"
					} else if strings.Contains(contentType, "application/json") {
						outputPath += ".json"
					} else if strings.Contains(contentType, "text/plain") {
						outputPath += ".txt"
					} else {
						outputPath += ".download"
					}
				}
			} else {
				outputPath = "downloaded_file"
			}
		}
		// Ensure we have a valid filename
		if outputPath == "" || outputPath == "." {
			outputPath = "index.html"
		}

		var file *os.File
		if options.Continue {
			file, err = os.OpenFile(outputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		} else {
			file, err = os.Create(outputPath)
		}
		if err != nil {
			fmt.Printf("创建文件失败: %v\n", err)
			continue
		}
		defer file.Close()

		_, err = io.Copy(file, resp.Body)
		if err != nil {
			fmt.Printf("写入文件失败: %v\n", err)
			continue
		}

		if !options.Quiet {
			duration := time.Since(startTime)
			fileInfo, _ := os.Stat(outputPath)
			size := fileInfo.Size()
			speed := float64(size) / duration.Seconds() / 1024

			fmt.Printf("Length: %d [%s]\n", size, resp.Header.Get("Content-Type"))
			fmt.Printf("Saving to: '%s'\n", outputPath)
			fmt.Printf("\n%s (%0.1f KB/s) - '%s' saved [%d/%d]\n",
				time.Now().Format("2006-01-02 15:04:05"),
				speed,
				outputPath,
				size,
				size)
		}

		if options.AISummarize {
			file.Seek(0, 0)
			content, err := io.ReadAll(file)
			if err != nil {
				fmt.Printf("读取文件内容失败: %v\n", err)
				continue
			}

			summaryPrompt := fmt.Sprintf("请总结以下下载内容:\n%s\n用中文简洁概括主要内容", string(content))
			fmt.Println("AI总结:")
			processQuery(summaryPrompt, true)
		}
	}
}
