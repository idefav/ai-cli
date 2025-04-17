package cmd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// processQuery 处理AI查询请求
func processQuery(apiKey, model, basePath string, stream bool) func(string, bool) {
	return func(prompt string, isSummary bool) {
		config := openai.DefaultConfig(apiKey)
		if basePath != "" {
			config.BaseURL = basePath
		}
		client := openai.NewClientWithConfig(config)

		if stream {
			req := openai.ChatCompletionRequest{
				Model: model,
				Messages: []openai.ChatCompletionMessage{
					{
						Role:    openai.ChatMessageRoleUser,
						Content: prompt,
					},
				},
				Stream: true,
			}

			stream, err := client.CreateChatCompletionStream(context.Background(), req)
			if err != nil {
				fmt.Printf("API调用失败: %v\n", err)
				os.Exit(1)
			}
			defer stream.Close()

			fmt.Println("AI回复:")
			for {
				response, err := stream.Recv()
				if err != nil {
					if errors.Is(err, io.EOF) {
						break
					}
					fmt.Printf("\n流式接收错误: %v\n", err)
					break
				}
				fmt.Print(response.Choices[0].Delta.Content)
			}
			fmt.Println()
		} else {
			resp, err := client.CreateChatCompletion(
				context.Background(),
				openai.ChatCompletionRequest{
					Model: model,
					Messages: []openai.ChatCompletionMessage{
						{
							Role:    openai.ChatMessageRoleUser,
							Content: prompt,
						},
					},
				},
			)

			if err != nil {
				fmt.Printf("API调用失败: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("\rAI回复: %s\n", resp.Choices[0].Message.Content)
		}
	}
}

var rootCmd = &cobra.Command{
	Use:   "ai-cli [问题]",
	Short: "AI命令行工具",
	Long: `AI命令行工具，提供LLM交互功能
不带参数运行时进入交互模式`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		apiKey := viper.GetString("ai.apiKey")
		model := viper.GetString("ai.model")
		basePath := viper.GetString("ai.basePath")
		stream := viper.GetBool("ai.stream")

		if apiKey == "" {
			fmt.Println("请在config.yaml中配置API密钥")
			os.Exit(1)
		}

		queryProcessor := processQuery(apiKey, model, basePath, stream)

		// 交互模式
		if len(args) == 0 {
			fmt.Println("ai-cli> 你好，请问有什么帮助么？(输入exit或quit退出)")

			// Set up interrupt handling
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, syscall.SIGINT)
			defer signal.Stop(sigChan)

			scanner := bufio.NewScanner(os.Stdin)
			for {
				fmt.Print("ai-cli> ")

				// Read input with scanner
				if !scanner.Scan() {
					return
				}
				input := scanner.Text()

				// Skip empty input
				if input == "" {
					continue
				}

				// Handle tab completion by checking if input ends with space+tab
				if strings.HasSuffix(input, " \t") {
					// Get the partial filename before tab
					lastSpace := strings.LastIndex(input[:len(input)-2], " ")
					prefix := ""
					toComplete := input
					if lastSpace >= 0 {
						prefix = input[:lastSpace+1]
						toComplete = input[lastSpace+1 : len(input)-2]
					}

					// Get matching files/dirs
					matches, _ := filepath.Glob(toComplete + "*")
					if len(matches) > 0 {
						// If single match, complete it
						if len(matches) == 1 {
							completed := matches[0]
							if isDir(completed) {
								completed += "/"
							}
							input = prefix + completed
						} else {
							// Show all matches
							fmt.Printf("\n%s\n", strings.Join(matches, "  "))
						}
					}

					// Reprint prompt with completed input
					fmt.Printf("\rai-cli> %s", input)
					continue
				}

				// Handle Ctrl+C interrupt
				select {
				case <-sigChan:
					// Clear current line and move cursor to start
					fmt.Print("\033[2K\r")
					// Only exit if we're not in a subcommand
					if !strings.HasPrefix(input, " ") {
						fmt.Println("感谢使用 AI-CLI, 欢迎再次使用!")
						return
					}
					fmt.Println("(按下Ctrl+C不会退出程序，输入exit或quit退出)")
					continue
				default:
					// Process the input
					if input == "exit" || input == "quit" {
						fmt.Println("感谢使用 AI-CLI, 欢迎再次使用!")
						return
					}
					if input == "clear" {
						HandleClear()
						continue
					}
					if strings.HasPrefix(input, "cat ") {
						HandleCat(input)
						continue
					}
					if strings.HasPrefix(input, "ls") || strings.HasPrefix(input, "ll") {
						HandleLs(input, queryProcessor)
						continue
					}
					if strings.HasPrefix(input, "curl ") {
						HandleCurl(input, queryProcessor)
						continue
					}
					if strings.HasPrefix(input, "wget ") {
						HandleWget(input, queryProcessor)
						continue
					}
					queryProcessor(input, false)
				}
			}
		}

		// 直接提问模式
		queryProcessor(args[0], false)
	},
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func init() {
	// 不再需要prompt标志参数
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
