#!/bin/bash

# 清理旧构建
echo "清理旧构建..."
rm -rf build/
mkdir -p build

# 编译Windows版本
echo "编译Windows版本..."
GOOS=windows GOARCH=amd64 go build -o build/ai-cli.exe
zip -j build/ai-cli-windows-amd64.zip build/ai-cli.exe README.md LICENSE config.yaml

# 编译Linux版本
echo "编译Linux版本..."
GOOS=linux GOARCH=amd64 go build -o build/ai-cli
tar -czvf build/ai-cli-linux-amd64.tar.gz -C build ai-cli README.md LICENSE config.yaml

# 编译MacOS版本
echo "编译MacOS版本..."
GOOS=darwin GOARCH=amd64 go build -o build/ai-cli
tar -czvf build/ai-cli-darwin-amd64.tar.gz -C build ai-cli README.md LICENSE config.yaml

echo "构建完成！"
echo "构建结果在build目录下:"
ls -lh build/
