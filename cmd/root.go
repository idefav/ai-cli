package cmd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

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

		processQuery := func(prompt string) {
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

				fmt.Println("AI回复(流式):")
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

		// 交互模式
		if len(args) == 0 {
			fmt.Println("ai-cli> 你好，请问有什么帮助么？(输入exit或quit退出)")
			scanner := bufio.NewScanner(os.Stdin)
			for {
				fmt.Print("ai-cli> ")
				if !scanner.Scan() {
					break
				}

				prompt := scanner.Text()
				if prompt == "exit" || prompt == "quit" {
					break
				}
				if prompt == "" {
					continue
				}
				processQuery(prompt)
			}
			return
		}

		// 直接提问模式
		processQuery(args[0])
	},
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
