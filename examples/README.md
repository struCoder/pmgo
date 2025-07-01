# PMGO 测试示例

本目录包含用于测试 PMGO 功能的示例应用程序和脚本。

## 📁 文件结构

```
examples/
├── README.md                 # 本说明文件
├── test-pmgo.sh             # 完整的测试脚本
├── test-config.yaml         # 测试配置文件
├── webserver/
│   └── main.go              # 测试用 Web 服务器
├── simple-app/
│   └── main.go              # 简单的测试应用
├── logs/                    # 日志目录
└── data/                    # 数据目录
```

## 🚀 快速开始

### 1. 构建 PMGO

```bash
# 在项目根目录
make build
```

### 2. 运行完整测试

```bash
# 运行基础测试
./examples/test-pmgo.sh

# 运行完整测试（包含负载测试）
./examples/test-pmgo.sh --full-test

# 查看帮助
./examples/test-pmgo.sh --help
```

### 3. 手动测试

#### 启动 PMGO 守护进程

```bash
./pmgo serve
# 或者使用测试配置
./pmgo serve --config examples/test-config.yaml
```

#### 启动测试应用

```bash
# 启动 Web 服务器
./pmgo start examples/webserver/main.go test-webserver --args="--port=8080"

# 启动简单应用
./pmgo start examples/simple-app/main.go simple-app my-simple-app

# 查看进程列表
./pmgo list

# 查看进程详情
./pmgo info test-webserver
```

## 🌐 测试应用说明

### Web 服务器 (webserver/main.go)

一个功能完整的 HTTP 服务器，包含：

- **健康检查**: `GET /health`
- **服务信息**: `GET /info`
- **主页**: `GET /`
- **工作模拟**: `GET /api/work`
- **请求计数**: `GET /api/counter`

**启动参数**:
- `--port=<端口>` 或第一个参数指定端口
- 默认端口: 8080

**特性**:
- 优雅关闭
- JSON API 响应
- HTML 界面
- 自动上报状态
- 支持负载测试

### 简单应用 (simple-app/main.go)

一个基础的长运行进程，特点：

- 定期输出状态
- 显示 PID 和启动信息
- 模拟工作负载
- 支持自定义名称

## 🧪 测试场景

### 1. 基础功能测试

```bash
# 测试进程启动
./pmgo start examples/simple-app/main.go test1

# 测试进程列表
./pmgo list

# 测试进程重启
./pmgo restart test1

# 测试进程停止
./pmgo stop test1

# 测试进程删除
./pmgo delete test1
```

### 2. Web 服务器测试

```bash
# 启动 Web 服务器
./pmgo start examples/webserver/main.go webserver --args="--port=8080"

# 测试 Web 接口
curl http://localhost:8080/health
curl http://localhost:8080/info
curl http://localhost:8080/api/work

# 在浏览器中查看
open http://localhost:8080
```

### 3. 多进程管理

```bash
# 启动多个进程
./pmgo start examples/simple-app/main.go app1 worker1
./pmgo start examples/simple-app/main.go app2 worker2
./pmgo start examples/webserver/main.go web --args="--port=8081"

# 批量管理
./pmgo list
./pmgo save
```

### 4. API 测试

```bash
# PMGO API 测试
curl http://localhost:9876/health
curl http://localhost:9876/api/v1/processes

# 启动进程 (API)
curl -X POST http://localhost:9876/api/v1/processes \
  -H "Content-Type: application/json" \
  -d '{
    "name": "api-test",
    "command": "examples/simple-app/main.go",
    "args": ["api-worker"]
  }'
```

## 🔧 配置说明

### 测试配置文件 (test-config.yaml)

专为测试优化的配置：

- **调试日志**: 更详细的输出
- **快速重启**: 2秒重启延迟
- **频繁检查**: 10秒健康检查
- **小文件**: 10MB 日志轮转

### 生产配置对比

| 配置项 | 测试环境 | 生产环境 |
|--------|----------|----------|
| 日志级别 | debug | info |
| 重启延迟 | 2s | 5s |
| 健康检查 | 10s | 30s |
| 文件大小 | 10MB | 100MB |

## 📊 性能测试

### 负载测试

```bash
# 运行负载测试
./examples/test-pmgo.sh --full-test

# 手动负载测试
for i in {1..100}; do
  curl -s http://localhost:8080/api/work &
done
wait

# 查看结果
curl http://localhost:8080/api/counter
```

### 压力测试

```bash
# 使用 ab 工具
ab -n 1000 -c 10 http://localhost:8080/api/work

# 使用 wrk 工具
wrk -t4 -c10 -d30s http://localhost:8080/api/work
```

## 🐛 故障排除

### 常见问题

1. **端口被占用**
   ```bash
   lsof -i :8080
   lsof -i :9876
   ```

2. **进程无法启动**
   ```bash
   ./pmgo logs test-webserver
   cat pmgo-daemon.log
   ```

3. **权限问题**
   ```bash
   chmod +x examples/test-pmgo.sh
   chmod +x examples/webserver/main.go
   ```

### 清理环境

```bash
# 停止所有进程
./pmgo list | grep -v "No processes" | tail -n +2 | while read line; do
  name=$(echo $line | awk '{print $1}')
  ./pmgo delete "$name"
done

# 停止守护进程
./pmgo kill

# 或者使用脚本清理
./examples/test-pmgo.sh --cleanup-only --kill-daemon
```

## 📝 日志查看

### 应用日志

```bash
# 查看特定应用日志
./pmgo logs test-webserver

# 查看守护进程日志
tail -f pmgo-daemon.log

# 查看配置的日志文件
tail -f examples/logs/pmgo-test.log
```

### 系统监控

```bash
# 查看进程状态
./pmgo list

# 查看系统资源
./pmgo info test-webserver

# API 监控
curl http://localhost:9876/metrics
```

## 🎯 自定义测试

### 创建自己的测试应用

```go
package main

import (
    "fmt"
    "os"
    "time"
)

func main() {
    fmt.Printf("My Test App (PID: %d) started\n", os.Getpid())

    for i := 0; ; i++ {
        fmt.Printf("Running iteration %d\n", i)
        time.Sleep(5 * time.Second)
    }
}
```

### 测试步骤

```bash
# 1. 保存为 my-test.go
# 2. 用 PMGO 启动
./pmgo start my-test.go my-test

# 3. 监控
./pmgo list
./pmgo info my-test

# 4. 清理
./pmgo delete my-test
```

## 🔗 相关链接

- [PMGO 主文档](../README.md)
- [配置文档](../configs/config.yaml)
- [API 文档](../docs/)
- [贡献指南](../CONTRIBUTING.md)