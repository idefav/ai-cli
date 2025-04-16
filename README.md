# AI CLI 工具

一个基于命令行的AI交互工具，支持流式和非流式响应。

## 功能特点

- 交互式聊天模式
- 直接提问模式
- 支持流式输出
- 可配置API端点
- 支持多种AI模型

## 安装

1. 确保已安装Go (1.16+)
2. 克隆项目：
   ```bash
   git clone https://github.com/your-repo/ai-cli.git
   cd ai-cli
   ```
3. 构建项目：
   ```bash
   go build
   ```

## 使用方法

### 交互模式
```bash
./ai-cli
```

### 直接提问模式
```bash
./ai-cli "你的问题"
```

### 流式输出
在config.yaml中设置：
```yaml
ai:
  stream: true
```

## 配置

创建config.yaml文件：
```yaml
ai:
  apiKey: "your-api-key"
  model: "deepseek-chat" # 或其他支持的模型
  basePath: "https://api.deepseek.com" # 可选，自定义API地址
  stream: false # 是否启用流式输出
```

## 示例

```bash
# 交互模式
$ ./ai-cli
ai-cli> 你好，请问有什么帮助么？(输入exit或quit退出)
ai-cli> 你好
AI回复: 你好！有什么我可以帮助你的吗？

# 直接提问模式
$ ./ai-cli "介绍一下你自己"
AI回复: 我是一个AI助手...
```

## 许可证

Apache 2.0 - 详见 [LICENSE](LICENSE) 文件
