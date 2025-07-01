# PMGO - 现代化的 Go 进程管理器

[English](../README.md) | **中文**

PMGO 是一个用 Go 语言编写的轻量级、现代化进程管理器，专为 Go 应用程序设计。它帮助您保持应用程序永远运行，具有自动重启、实时监控和基于 Web 的管理界面等功能。

## ✨ 核心特性

- 🚀 **现代化架构**: 基于 Go 1.21+ 构建，采用清洁代码和最新 Go 模式
- 🔄 **自动重启**: 智能进程监控，故障时自动重启
- 🌐 **Web 界面**: 精美的 Web 仪表板进行进程管理
- 📊 **实时监控**: CPU、内存使用率和性能指标
- 🔌 **RESTful API**: 完整的 HTTP API 用于编程控制
- 🛡️ **安全性**: 安全的守护进程模式，具有适当的进程隔离
- 📝 **丰富日志**: 结构化日志，支持多种输出格式
- ⚡ **高性能**: 优化的低资源使用

## 📖 目录

- [安装](#安装)
- [快速开始](#快速开始)
- [帮助系统](#帮助系统)
- [配置](#配置)
- [命令参考](#命令参考)
- [API 文档](#api-文档)
- [Web 界面](#web-界面)
- [开发指南](#开发指南)
- [常见问题](#常见问题)

## 🚀 安装

### 使用安装脚本（推荐）

```bash
curl -sSL https://raw.githubusercontent.com/struCoder/pmgo/master/scripts/install.sh | bash
```

### 从源码编译

```bash
git clone https://github.com/struCoder/pmgo.git
cd pmgo
make build
sudo mv pmgo /usr/local/bin/
```

### 使用 Go 命令

```bash
go install github.com/struCoder/pmgo@latest
```

### Docker

```bash
docker pull strucoder/pmgo:latest
docker run -p 9876:9876 -p 8080:8080 strucoder/pmgo:latest
```

## 🎯 快速开始

### 1. 启动守护进程

```bash
# 启动 PMGO 守护进程
pmgo serve

# 后台运行
pmgo serve --daemon
```

### 2. 管理进程

```bash
# 启动 Go 源文件
pmgo start app.go myapp --args="port=8080"

# 启动二进制可执行文件
pmgo start ./myapp myapp --binary --args="config.yaml"

# 使用多个参数启动
pmgo start ./server webserver -b --args="--port=8080" --args="--env=prod"

# 查看所有进程状态和资源使用情况
pmgo list

# 查看进程详细信息
pmgo info myapp

# 重启进程
pmgo restart myapp

# 停止进程（保留在进程列表中）
pmgo stop myapp

# 永久删除进程
pmgo delete myapp

# 显示帮助信息和示例
pmgo help
pmgo help start  # 查看start命令的详细帮助
```

### 3. Web 界面

打开浏览器访问 `http://localhost:8080` 来使用 Web 管理界面。

## 💡 帮助系统

PMGO 提供了完善的帮助系统，每个命令都有详细的说明和使用示例。

### 查看总体帮助

```bash
pmgo help
```

输出示例：
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

### 查看命令详细帮助

```bash
pmgo help start
```

输出示例：
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

### 常用帮助命令

```bash
pmgo help              # 查看所有命令概览
pmgo help start        # 启动进程的详细说明
pmgo help list         # 进程列表的详细说明
pmgo help info         # 进程信息的详细说明
pmgo help restart      # 重启进程的详细说明
pmgo help stop         # 停止进程的详细说明
pmgo help delete       # 删除进程的详细说明
```

## ⚙️ 配置

PMGO 支持 YAML、TOML 和 JSON 格式的配置文件。默认配置文件位置：`~/.pmgo/config.yaml`

### 基本配置

```yaml
# 服务器配置
server:
  host: "localhost"
  port: 9876        # API 端口
  web_port: 8080    # Web 界面端口
  tls_enabled: false

# 日志配置
logging:
  level: "info"     # debug, info, warn, error
  format: "text"    # text, json
  file: ""          # 日志文件路径（空表示输出到控制台）

# 进程管理配置
processes:
  default_restart_policy: "always"    # always, on-failure, no
  max_restart_attempts: 5
  restart_delay: "1s"
  health_check_interval: "30s"
```

### 高级配置

```yaml
# 安全配置
security:
  auth_enabled: true
  jwt_secret: "your-secret-key"
  api_keys:
    - "api-key-1"
    - "api-key-2"

# 通知配置
notifications:
  enabled: true
  webhook_url: "https://your-webhook.com"
  slack_webhook: "https://hooks.slack.com/..."
  email:
    smtp_host: "smtp.gmail.com"
    smtp_port: 587
    username: "your-email@gmail.com"
    password: "your-password"
```

## 📚 命令参考

### 帮助和示例

```bash
# 显示所有命令和使用示例
pmgo help

# 显示特定命令的详细帮助
pmgo help start
pmgo help list
pmgo help info
```

### 守护进程管理

```bash
# 启动守护进程（通常自动启动）
pmgo serve

# 停止守护进程和所有管理的进程
pmgo kill
```

### 进程管理

```bash
# 启动进程（新的优化格式）
pmgo start <source> <name> [flags]
  -b, --binary         将源文件视为二进制可执行文件
      --args=ARGS ...  传递给进程的参数

# 使用示例：
pmgo start app.go myapp --args="port=8080"          # Go 源文件
pmgo start ./mybin myapp --binary                   # 二进制文件
pmgo start app.go myapp -b --args="config.yaml"     # 多个参数

# 进程控制
pmgo list                  # 列出所有进程及资源使用情况
pmgo info myapp           # 查看进程详细信息和统计数据
pmgo restart myapp        # 重启进程
pmgo stop myapp           # 停止进程（保留在进程列表中）
pmgo delete myapp         # 永久删除进程和所有文件
```

### 其他命令

```bash
pmgo save              # 保存当前进程列表到配置文件
pmgo version           # 显示 PMGO 版本信息
```

## 🔌 API 文档

PMGO 提供完整的 RESTful API，支持所有进程管理操作。

### 基础 URL

```
http://localhost:9876/api/v1
```

### 端点

#### 进程管理

```bash
GET    /processes           # 获取所有进程
POST   /processes           # 启动新进程
GET    /processes/{name}    # 获取特定进程信息
DELETE /processes/{name}    # 停止并删除进程
POST   /processes/{name}/restart  # 重启进程
POST   /processes/{name}/stop     # 停止进程
```

#### 系统信息

```bash
GET    /health             # 健康检查
GET    /metrics            # 系统指标
GET    /version            # 版本信息
```

### API 示例

#### 启动进程

```bash
curl -X POST http://localhost:9876/api/v1/processes \
  -H "Content-Type: application/json" \
  -d '{
    "name": "myapp",
    "command": "./myapp",
    "args": ["--port=8080"],
    "restart_policy": "always",
    "max_restarts": 5
  }'
```

#### 获取进程状态

```bash
curl http://localhost:9876/api/v1/processes/myapp
```

## 🌐 Web 界面

PMGO 提供了一个现代化的 Web 管理界面，默认运行在 `http://localhost:8080`。

### 功能特性

- 📊 **实时仪表板**: 系统资源使用情况和进程状态
- 🔄 **进程管理**: 启动、停止、重启、删除进程
- 📈 **监控图表**: CPU、内存使用趋势图
- 📝 **日志查看**: 实时日志流和历史日志
- ⚙️ **配置管理**: 在线编辑配置文件
- 🔔 **通知中心**: 进程事件和系统通知

### 界面截图

```
┌─────────────────────────────────────────────────────────────┐
│ PMGO Process Manager                              [Settings] │
├─────────────────────────────────────────────────────────────┤
│ Dashboard | Processes | Logs | Monitoring | Settings        │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│ System Overview                    Recent Activities        │
│ ┌─────────────────┐               ┌─────────────────────┐    │
│ │ CPU: 15%        │               │ myapp started       │    │
│ │ Memory: 2.1GB   │               │ worker restarted    │    │
│ │ Processes: 3    │               │ config updated      │    │
│ └─────────────────┘               └─────────────────────┘    │
│                                                             │
│ Active Processes                                            │
│ ┌─────────────────────────────────────────────────────────┐ │
│ │ Name     │ Status  │ CPU  │ Memory │ Uptime   │ Actions │ │
│ │ myapp    │ Running │ 5%   │ 128MB  │ 2h 15m   │ [▶][⏸][🗑] │ │
│ │ worker   │ Running │ 3%   │ 64MB   │ 1h 30m   │ [▶][⏸][🗑] │ │
│ │ api      │ Stopped │ 0%   │ 0MB    │ -        │ [▶][⏸][🗑] │ │
│ └─────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

## 💻 开发指南

### 开发环境设置

```bash
# 克隆仓库
git clone https://github.com/struCoder/pmgo.git
cd pmgo

# 安装依赖
make deps

# 运行测试
make test

# 启动开发模式（热重载）
make dev
```

### 项目结构

```
pmgo/
├── cmd/                    # 应用程序入口点
│   └── pmgo/              # 主应用程序
├── internal/              # 私有应用程序代码
│   ├── api/              # HTTP API 处理器
│   ├── config/           # 配置管理
│   ├── daemon/           # 守护进程实现
│   ├── master/           # 进程主管理器
│   ├── process/          # 进程管理
│   └── web/              # Web 界面
├── pkg/                   # 公共库代码
├── configs/              # 配置文件
├── docs/                 # 文档
├── scripts/              # 构建和部署脚本
└── tests/                # 集成测试
```

### 贡献指南

我们欢迎贡献！请查看 [CONTRIBUTING.md](../CONTRIBUTING.md) 了解详细信息。

## 🔧 故障排除

### 常见问题

#### 1. 端口已被占用

```bash
# 检查端口使用情况
lsof -i :9876
netstat -tlnp | grep 9876

# 更改端口
pmgo serve --port 9877
```

#### 2. 权限问题

```bash
# 检查文件权限
ls -la ~/.pmgo/

# 修复权限
chmod 755 ~/.pmgo/
chmod 644 ~/.pmgo/config.yaml
```

#### 3. 进程启动失败

```bash
# 查看详细日志
pmgo logs myapp

# 启用调试模式
pmgo serve --log-level debug
```

#### 4. Web 界面无法访问

```bash
# 检查 Web 服务状态
curl http://localhost:8080/health

# 检查防火墙设置
sudo ufw status
```

### 性能优化

#### 系统资源监控

```bash
# 查看系统资源使用
pmgo info system

# 监控特定进程
pmgo monitor myapp
```

#### 配置优化

```yaml
# 优化配置示例
processes:
  health_check_interval: "10s"    # 减少检查间隔
  output_buffer_size: 2048       # 增加缓冲区大小

logging:
  level: "warn"                   # 减少日志输出
  rotate: true                    # 启用日志轮转
```

## 📈 性能基准

PMGO 设计时考虑了效率：

- **内存占用**: 守护进程 < 10MB
- **CPU 使用**: 正常负载下 < 1%
- **启动时间**: 守护进程 < 100ms
- **进程限制**: 支持 1000+ 并发进程

## 🤝 社区

- **GitHub Issues**: 错误报告和功能请求
- **GitHub Discussions**: 社区支持和想法交流
- **微信群**: 添加微信号 `pmgo-support` 加入讨论群

## 📄 许可证

MIT License - 查看 [LICENSE](../LICENSE) 了解详情。

## 🙏 致谢

- 灵感来自 [PM2](https://pm2.keymetrics.io/)
- 由 Go 社区 ❤️ 构建