@echo off
REM 清理旧构建
echo 清理旧构建...
if exist build rmdir /s /q build
mkdir build

REM 编译Windows版本
echo 编译Windows版本...
set GOOS=windows
set GOARCH=amd64
go build -o build\ai-cli.exe
powershell Compress-Archive -Path build\ai-cli.exe,README.md,LICENSE,config.yaml -DestinationPath build\ai-cli-windows-amd64.zip -Force

REM 编译Linux版本
echo 编译Linux版本...
set GOOS=linux
set GOARCH=amd64
go build -o build\ai-cli-linux
echo 请手动打包Linux版本: build\ai-cli-linux + README.md + LICENSE + config.yaml

REM 编译MacOS版本
echo 编译MacOS版本...
set GOOS=darwin
set GOARCH=amd64
go build -o build\ai-cli-darwin
echo 请手动打包MacOS版本: build\ai-cli-darwin + README.md + LICENSE + config.yaml

echo 构建完成!
echo Windows版本已打包为: build\ai-cli-windows-amd64.zip
