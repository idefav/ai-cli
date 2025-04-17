package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// HandleClear 处理clear命令，清空终端屏幕（跨平台支持）
func HandleClear() {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	case "linux", "darwin": // darwin是MacOS
		cmd = exec.Command("clear")
	default:
		fmt.Println("不支持的操作系统")
		return
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}
