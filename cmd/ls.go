package cmd

import (
	"fmt"
	"os"
	"strings"
)

// HandleLs 处理ls/ll命令，列出目录内容
func HandleLs(prompt string, processQuery func(string, bool)) {
	isLl := strings.HasPrefix(prompt, "ll")
	showDetails := isLl || strings.Contains(prompt, "-l") || strings.Contains(prompt, "-lh")
	humanReadable := isLl && strings.Contains(prompt, "-h") || strings.Contains(prompt, "-h") || strings.Contains(prompt, "-lh")
	shouldSummarize := strings.ContainsAny(prompt, "s") && (strings.Contains(prompt, "-s") || strings.Contains(prompt, "-ls") ||
		strings.Contains(prompt, "-hs") || strings.Contains(prompt, "-lhs") || strings.Contains(prompt, "-lls"))

	// 列出当前目录
	files, err := os.ReadDir(".")
	if err != nil {
		fmt.Printf("无法读取目录: %v\n", err)
		return
	}

	var output strings.Builder
	for _, file := range files {
		if !showDetails {
			fmt.Fprintln(&output, file.Name())
		} else {
			// 详细格式
			info, err := file.Info()
			if err != nil {
				fmt.Printf("无法获取文件信息: %v\n", err)
				continue
			}

			size := info.Size()
			sizeStr := fmt.Sprintf("%8d", size)
			if humanReadable {
				sizeStr = formatSize(size)
			}

			// 格式化输出：权限 大小 修改时间 文件名
			fmt.Fprintf(&output, "%s %8s %s %s\n",
				info.Mode().String(),
				sizeStr,
				info.ModTime().Format("Jan _2 15:04"),
				file.Name())
		}
	}

	// 总是先打印原始结果
	fmt.Print(output.String())

	// 如果需要总结，发送给AI
	if shouldSummarize {
		summaryPrompt := fmt.Sprintf("请总结以下文件列表:\n%s\n用中文简洁概括目录内容", output.String())
		fmt.Println("AI总结:")
		processQuery(summaryPrompt, true)
	}
}

func formatSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%dB", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%cB", float64(size)/float64(div), "KMGTPE"[exp])
}
