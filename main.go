package main

import (
	"ai-cli/cmd"
	"github.com/spf13/viper"
)

func main() {
	// 初始化viper配置
	viper.SetConfigName("config") // 配置文件名 (不带扩展名)
	viper.SetConfigType("yaml")   // 配置文件类型
	viper.AddConfigPath(".")      // 在当前目录查找配置文件
	
	// 读取配置文件
	viper.ReadInConfig() // 忽略错误，配置项会有默认值或运行时检查
	
	cmd.Execute()
}
