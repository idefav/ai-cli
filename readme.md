# AI CLI Tool

[中文文档](readme_cn.md)

A command-line AI interaction tool supporting both streaming and non-streaming responses.

## Features

- Interactive chat mode
- Direct question mode  
- Streaming output support
- Configurable API endpoint
- Multiple AI model support
- Cross-platform support (Windows/Linux/MacOS)
- Easy configuration via YAML file
- Clean and intuitive interface
- Linux-style tab completion (trigger with space+Tab)

## Installation

### Prerequisites
- Go 1.16 or higher
- Git (for cloning the repository)

### Steps
1. Clone the repository:
   ```bash
   git clone git@github.com:idefav/ai-cli.git
   cd ai-cli
   ```
2. Build the project:
   ```bash
   go build
   ```
3. (Optional) Add to PATH for global access

## Usage

### Basic Commands
```bash
# Interactive mode (conversation)
./ai-cli

# Direct question mode
./ai-cli "your question here"

# Help
./ai-cli --help
```

### Streaming Mode
Enable in config.yaml:
```yaml
ai:
  stream: true  # Set to true for streaming responses
```

## Configuration

Configuration files can be placed in either:
1. Current working directory: `./config.yaml` (higher priority)
2. User home directory: `~/.ai-cli/config.yaml`

Example config:
```yaml
ai:
  apiKey: "your-api-key-here"  # Required: Your API key
  model: "default-model"      # Default AI model  
  basePath: ""                # Optional: Custom API endpoint
  stream: false               # Enable streaming response
```

## Examples

### Interactive Session
```bash
$ ./ai-cli
> Hello! How can I assist you today?
> What's the weather like?
The current weather is sunny with a temperature of 22°C.
> exit
```

### Direct Query
```bash
$ ./ai-cli "Explain quantum computing"
Quantum computing is a type of computation...
```

## Build & Distribution

Use the provided build scripts:
```bash
# Linux/MacOS
./build.sh

# Windows
build.bat
```

## Development Notice

This software was entirely generated by AI without any human-written code.

## License

Apache 2.0 License - See [LICENSE](LICENSE) file for details.

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for recent updates and changes.

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Open a pull request
