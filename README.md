<div align="center">
<a>
   <img src="https://i.loli.net/2018/12/06/5c08b9a294c29.png">
</a>
<br/>
<b>PMGO</b>
<br/><br/>

<a href="https://goreportcard.com/report/github.com/struCoder/pmgo">
  <img src="https://goreportcard.com/badge/github.com/struCoder/pmgo" alt="Go Report Card" />
</a>

<a href="https://godoc.org/github.com/struCoder/pmgo">
  <img src="https://godoc.org/github.com/struCoder/pmgo?status.svg" alt="GoDoc" />
</a>

<a href="https://github.com/struCoder/pmgo/releases">
  <img src="https://img.shields.io/github/v/release/struCoder/pmgo" alt="Release" />
</a>

<a href="LICENSE">
  <img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="License" />
</a>
<br/><br/>
</div>

# PMGO

[English](#english) | [中文](#中文)

## English

PMGO is a lightweight, modern process manager written in Go for Go applications. It helps you keep your applications alive forever, with features like auto-restart, real-time monitoring, and web-based management interface.

### ✨ Key Features

- 🚀 **Modern Architecture**: Built with Go 1.21+, featuring clean code and latest Go patterns
- 🔄 **Auto-restart**: Intelligent process monitoring with automatic restart on failure
- 🌐 **Web Interface**: Beautiful web dashboard for process management
- 📊 **Real-time Monitoring**: CPU, memory usage, and performance metrics
- 🔌 **RESTful API**: Complete HTTP API for programmatic control
- 🛡️ **Security**: Secure daemon mode with proper process isolation
- 📝 **Rich Logging**: Structured logging with multiple output formats
- ⚡ **High Performance**: Optimized for low resource usage

### 🚀 Quick Start

#### Installation

```bash
# Install from source
go install github.com/struCoder/pmgo@latest

# Or download pre-built binary
curl -sSL https://github.com/struCoder/pmgo/releases/latest/download/pmgo-linux-amd64 -o pmgo
chmod +x pmgo && sudo mv pmgo /usr/local/bin/
```

#### Basic Usage

```bash
# Start daemon (auto-starts when needed)
pmgo serve

# Start a Go source file
pmgo start app.go myapp --args="port=8080"

# Start a binary executable
pmgo start ./myapp myapp --binary --args="config.yaml"

# Start with multiple arguments
pmgo start ./server webserver -b --args="--port=8080" --args="--env=prod"

# List all processes with detailed status
pmgo list

# View detailed process information
pmgo info myapp

# Restart a process
pmgo restart myapp

# Stop a process (keeps in process list)
pmgo stop myapp

# Delete a process permanently
pmgo delete myapp

# Show help with examples
pmgo help
pmgo help start  # Command-specific help
```

### 💡 Help System

PMGO provides comprehensive help with detailed examples for every command.

#### General Help

```bash
pmgo help
```

Sample output:
```
PMGO - A lightweight process manager for Go applications.

Examples:
  pmgo start app.go myapp                    # Start Go source file
  pmgo start ./mybin myapp --binary          # Start binary executable
  pmgo start app.go myapp --args="port=8080" # Start with arguments
  pmgo list                                  # List all processes

Commands:
  start [<flags>] <source> <name>
    Start and monitor a Go application or binary executable.

  list
    List all managed processes with status and resource usage.

  info <name>
    Show detailed information and statistics about a process.
```

#### Command-Specific Help

```bash
pmgo help start
```

Sample output:
```
usage: pmgo start [<flags>] <source> <name>

Start and monitor a Go application or binary executable.

Examples:
  pmgo start app.go myapp --args="port=8080"   # Go source
  pmgo start ./mybin myapp --binary            # Binary file
  pmgo start app.go myapp -b --args="config"   # Multiple args

Flags:
  -b, --binary         Treat source as binary executable (not Go source)
      --args=ARGS ...  Arguments to pass to the process

Args:
  <source>  Go source file (.go) or binary executable path
  <name>    Unique process name for management
```

#### Available Help Commands

```bash
pmgo help              # Show all commands overview
pmgo help start        # Detailed start command help
pmgo help list         # Detailed list command help
pmgo help info         # Detailed info command help
pmgo help restart      # Detailed restart command help
pmgo help stop         # Detailed stop command help
pmgo help delete       # Detailed delete command help
```

### 📖 Advanced Usage

#### Configuration File

PMGO supports YAML, TOML, and JSON configuration files:

```yaml
# ~/.pmgo/config.yaml
server:
  host: "localhost"
  port: 9876
  web_port: 8080

logging:
  level: "info"
  format: "json"
  file: "/var/log/pmgo.log"

processes:
  default_restart_policy: "always"
  max_restart_attempts: 5
  restart_delay: "1s"
```

#### Command Reference

```bash
# Help and Examples
pmgo help                           # Show all commands with examples
pmgo help start                     # Show detailed start command help

# Process Management
pmgo start <source> <name> [flags]  # Start and monitor an application
  -b, --binary                      # Treat source as binary executable
      --args=ARGS ...               # Arguments to pass to the process

# Examples:
pmgo start app.go myapp --args="port=8080"          # Go source file
pmgo start ./mybin myapp --binary                   # Binary executable
pmgo start app.go myapp -b --args="config.yaml"     # Multiple args

# Process Control
pmgo list                           # List processes with resource usage
pmgo info myapp                     # Show detailed process information
pmgo restart myapp                  # Restart a managed process
pmgo stop myapp                     # Stop process (keeps in list)
pmgo delete myapp                   # Delete process permanently

# Daemon Management
pmgo serve                          # Start PMGO daemon
pmgo kill                           # Stop daemon and all processes
```

#### API Examples

```bash
# Start process via API
curl -X POST http://localhost:9876/api/v1/processes \
  -H "Content-Type: application/json" \
  -d '{
    "name": "myapp",
    "command": "./myapp",
    "args": ["--port=8080"],
    "restart_policy": "always"
  }'

# Get process status
curl http://localhost:9876/api/v1/processes/myapp

# Stop process
curl -X DELETE http://localhost:9876/api/v1/processes/myapp
```

### 🏗️ Architecture

PMGO follows a modern, microservice-inspired architecture:

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   CLI Client    │    │   Web Interface │    │   HTTP API      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
                    ┌─────────────────┐
                    │   PMGO Daemon   │
                    └─────────────────┘
                             │
                    ┌─────────────────┐
                    │  Process Pool   │
                    └─────────────────┘
```

---

## 中文

PMGO 是一个用 Go 语言编写的轻量级、现代化进程管理器，专为 Go 应用程序设计。它帮助您保持应用程序永远运行，具有自动重启、实时监控和基于 Web 的管理界面等功能。

### ✨ 核心特性

- 🚀 **现代化架构**: 基于 Go 1.21+ 构建，采用清洁代码和最新 Go 模式
- 🔄 **自动重启**: 智能进程监控，故障时自动重启
- 🌐 **Web 界面**: 精美的 Web 仪表板进行进程管理
- 📊 **实时监控**: CPU、内存使用率和性能指标
- 🔌 **RESTful API**: 完整的 HTTP API 用于编程控制
- 🛡️ **安全性**: 安全的守护进程模式，具有适当的进程隔离
- 📝 **丰富日志**: 结构化日志，支持多种输出格式
- ⚡ **高性能**: 优化的低资源使用

### 🚀 快速开始

#### 安装

```bash
# 从源码安装
go install github.com/struCoder/pmgo@latest

# 或下载预编译二进制文件
curl -sSL https://github.com/struCoder/pmgo/releases/latest/download/pmgo-linux-amd64 -o pmgo
chmod +x pmgo && sudo mv pmgo /usr/local/bin/
```

#### 基本使用

```bash
# 启动守护进程
pmgo serve

# 启动 Go 应用程序
pmgo start /path/to/your/app app-name

# 列出所有进程
pmgo list

# Web 界面（默认：http://localhost:8080）
pmgo web

# 停止进程
pmgo stop app-name

# 重启进程
pmgo restart app-name

# 查看进程详情
pmgo info app-name

# 删除进程
pmgo delete app-name
```

### 🔧 Development

#### Building from Source

```bash
git clone https://github.com/struCoder/pmgo.git
cd pmgo
make build
```

#### Running Tests

```bash
make test
make test-coverage
```

#### Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details.

### 📊 Performance

PMGO is designed for efficiency:

- **Memory footprint**: < 10MB for daemon
- **CPU usage**: < 1% under normal load
- **Startup time**: < 100ms for daemon
- **Process limit**: 1000+ concurrent processes

### 🤝 Community

- **GitHub Issues**: Bug reports and feature requests
- **Discussions**: Community support and ideas
- **Discord**: Real-time chat and support

### 📄 License

MIT License - see [LICENSE](LICENSE) for details.

### 🙏 Acknowledgments

- Inspired by [PM2](https://pm2.keymetrics.io/)
- Built with ❤️ by the Go community
