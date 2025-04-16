package main

import (
	"ai-cli/cmd"
	"os"
	"path/filepath"
	"github.com/spf13/viper"
)

func main() {
	// 初始化viper配置
	viper.SetConfigName("config") // 配置文件名 (不带扩展名)
	viper.SetConfigType("yaml")   // 配置文件类型
	
	// 添加多个配置路径，按优先级顺序
	viper.AddConfigPath(".")      // 1. 当前执行目录
	home, err := os.UserHomeDir() // 2. 用户主目录下的.ai-cli目录
	if err == nil {
		viper.AddConfigPath(filepath.Join(home, ".ai-cli"))
	}
	
	// 读取配置文件
	viper.ReadInConfig() // 忽略错误，配置项会有默认值或运行时检查
	
	cmd.Execute()
}
