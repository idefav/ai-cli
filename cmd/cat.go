package cmd

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type CatOptions struct {
	ShowAll         bool
	NumberNonblank  bool
	ShowEnds        bool
	Number          bool
	SqueezeBlank    bool
	ShowTabs        bool
	ShowNonprinting bool
	Help            bool
	Version         bool
}

func parseCatArgs(args []string) (*CatOptions, []string, error) {
	options := &CatOptions{}
	flagSet := flag.NewFlagSet("cat", flag.ContinueOnError)

	flagSet.BoolVar(&options.ShowAll, "A", false, "equivalent to -vET")
	flagSet.BoolVar(&options.ShowAll, "show-all", false, "equivalent to -vET")
	flagSet.BoolVar(&options.NumberNonblank, "b", false, "number nonempty output lines")
	flagSet.BoolVar(&options.NumberNonblank, "number-nonblank", false, "number nonempty output lines")
	flagSet.BoolVar(&options.ShowEnds, "E", false, "display $ at end of each line")
	flagSet.BoolVar(&options.ShowEnds, "show-ends", false, "display $ at end of each line")
	flagSet.BoolVar(&options.Number, "n", false, "number all output lines")
	flagSet.BoolVar(&options.Number, "number", false, "number all output lines")
	flagSet.BoolVar(&options.SqueezeBlank, "s", false, "suppress repeated empty output lines")
	flagSet.BoolVar(&options.SqueezeBlank, "squeeze-blank", false, "suppress repeated empty output lines")
	flagSet.BoolVar(&options.ShowTabs, "T", false, "display TAB characters as ^I")
	flagSet.BoolVar(&options.ShowTabs, "show-tabs", false, "display TAB characters as ^I")
	flagSet.BoolVar(&options.ShowNonprinting, "v", false, "use ^ and M- notation")
	flagSet.BoolVar(&options.ShowNonprinting, "show-nonprinting", false, "use ^ and M- notation")
	flagSet.BoolVar(&options.Help, "help", false, "display help")
	flagSet.BoolVar(&options.Version, "version", false, "display version")

	err := flagSet.Parse(args)
	if err != nil {
		return nil, nil, err
	}

	return options, flagSet.Args(), nil
}

func HandleCat(prompt string) {
	args := strings.Fields(prompt)[1:] // Skip "cat"
	options, files, err := parseCatArgs(args)
	if err != nil {
		fmt.Printf("Error parsing arguments: %v\n", err)
		return
	}

	if options.Help {
		printCatHelp()
		return
	}

	if options.Version {
		fmt.Println("cat (ai-cli) 1.0")
		return
	}

	if len(files) == 0 {
		// Read from stdin
		catFile(os.Stdin, "", options)
		return
	}

	for _, file := range files {
		if file == "-" {
			catFile(os.Stdin, "", options)
			continue
		}

		f, err := os.Open(file)
		if err != nil {
			fmt.Printf("cat: %s: %v\n", file, err)
			continue
		}
		defer f.Close()

		catFile(f, file, options)
	}
}

func catFile(f *os.File, filename string, options *CatOptions) {
	reader := bufio.NewReader(f)
	lineNum := 1
	lastLineEmpty := false

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				fmt.Printf("cat: %s: %v\n", filename, err)
			}
			break
		}

		// Handle -s (squeeze-blank)
		if options.SqueezeBlank {
			if strings.TrimSpace(line) == "" {
				if lastLineEmpty {
					continue
				}
				lastLineEmpty = true
			} else {
				lastLineEmpty = false
			}
		}

		// Handle -n (number) and -b (number-nonblank)
		if options.Number || (options.NumberNonblank && strings.TrimSpace(line) != "") {
			fmt.Printf("%6d\t", lineNum)
			lineNum++
		}

		// Process line content based on options
		processed := line
		if options.ShowNonprinting || options.ShowAll {
			processed = showNonprinting(processed)
		}
		if options.ShowEnds || options.ShowAll {
			processed = strings.TrimRight(processed, "\n") + "$\n"
		}
		if options.ShowTabs || options.ShowAll {
			processed = strings.ReplaceAll(processed, "\t", "^I")
		}

		fmt.Print(processed)
	}
}

func showNonprinting(s string) string {
	var result strings.Builder
	for _, r := range s {
		switch {
		case r == '\t' || r == '\n':
			result.WriteRune(r)
		case r < 32:
			result.WriteString(fmt.Sprintf("^%c", r+64))
		case r == 127:
			result.WriteString("^?")
		case r > 127:
			result.WriteString(fmt.Sprintf("M-%c", r-128))
		default:
			result.WriteRune(r)
		}
	}
	return result.String()
}

func printCatHelp() {
	fmt.Println(`Usage: cat [OPTION]... [FILE]...
Concatenate FILE(s) to standard output.

With no FILE, or when FILE is -, read standard input.

  -A, --show-all           equivalent to -vET
  -b, --number-nonblank    number nonempty output lines, overrides -n
  -E, --show-ends          display $ at end of each line
  -n, --number             number all output lines
  -s, --squeeze-blank      suppress repeated empty output lines
  -T, --show-tabs          display TAB characters as ^I
  -v, --show-nonprinting   use ^ and M- notation, except for LFD and TAB
      --help        display this help and exit
      --version     output version information and exit

Examples:
  cat f - g  Output f's contents, then standard input, then g's contents.
  cat        Copy standard input to standard output.`)
}
