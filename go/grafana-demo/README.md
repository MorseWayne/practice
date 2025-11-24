# Grafana + Prometheus + Go 监控实践项目

这是一个完整的监控系统实践项目，展示如何使用 Go 应用暴露 Prometheus 指标，并通过 Grafana 进行可视化。

## 📋 项目简介

本项目包含：

- **Go 应用**：一个带有多个 HTTP 端点的 Web 服务，暴露 Prometheus 格式的指标
- **Prometheus**：收集和存储时序数据
- **Grafana**：可视化仪表板，展示应用指标

## 🏗️ 项目结构

```text
grafana-demo/
├── main.go                          # Go 应用主文件
├── go.mod                           # Go 依赖管理
├── go.sum                           # Go 依赖校验
├── Dockerfile                       # Go 应用容器化配置
├── docker-compose.yml               # Docker Compose 编排文件
├── Makefile                         # 便捷命令
├── README.md                        # 项目文档
├── prometheus/
│   └── prometheus.yml               # Prometheus 配置
└── grafana/
    ├── provisioning/
    │   ├── datasources/
    │   │   └── prometheus.yml       # Grafana 数据源配置
    │   └── dashboards/
    │       └── default.yml          # 仪表板自动加载配置
    └── dashboards/
        └── go-app-dashboard.json    # 预配置的仪表板
```

## 🚀 快速开始

### 前置要求

- Docker 和 Docker Compose
- Go 1.21+ (如果本地运行)

### 方式一：使用 Docker Compose（推荐）

1. **启动所有服务**

   ```bash
   make up
   # 或者
   docker-compose up -d
   ```

2. **访问服务**
   - Go 应用: <http://localhost:8080>
   - Prometheus: <http://localhost:9090>
   - Grafana: <http://localhost:3000> (用户名/密码: admin/admin)

3. **查看日志**

   ```bash
   make logs
   # 或者
   docker-compose logs -f
   ```

4. **停止服务**

   ```bash
   make down
   # 或者
   docker-compose down
   ```

### 方式二：本地运行 Go 应用

1. **安装依赖**

   ```bash
   go mod download
   ```

2. **运行应用**

   ```bash
   make run
   # 或者
   go run main.go
   ```

3. **使用 Docker 启动 Prometheus 和 Grafana**

   ```bash
   docker-compose up -d prometheus grafana
   ```

## 📊 暴露的指标说明

### 1. Counter（计数器）

- `http_requests_total`: HTTP 请求总数
  - 标签: `path`, `method`, `status`
- `orders_total`: 处理的订单总数

### 2. Gauge（仪表盘）

- `active_connections`: 当前活跃连接数
- `order_amount_current`: 当前订单金额

### 3. Histogram（直方图）

- `http_request_duration_seconds`: HTTP 请求持续时间分布

### 4. Summary（摘要）

- `http_response_size_bytes`: HTTP 响应大小分布

## 🎯 API 端点

| 端点 | 方法 | 描述 |
|------|------|------|
| `/` | GET | 应用主页，显示可用端点 |
| `/metrics` | GET | Prometheus 指标端点 |
| `/api/data` | GET | 示例 API 端点（模拟延迟） |
| `/health` | GET | 健康检查端点 |

## 📈 Grafana 仪表板

预配置的仪表板包含以下面板：

1. **HTTP 请求速率** - 显示每秒请求数
2. **活跃连接数** - 仪表盘显示当前连接数
3. **HTTP 请求延迟** - P95 和 P99 延迟
4. **订单总数** - 累计订单数量
5. **当前订单金额** - 实时订单金额变化

仪表板会自动加载，访问 Grafana 后即可看到。

## 🔧 Makefile 命令

```bash
make help          # 显示所有可用命令
make up            # 启动所有服务
make down          # 停止所有服务
make restart       # 重启所有服务
make logs          # 查看服务日志
make build         # 构建 Go 应用
make run           # 本地运行 Go 应用
make test          # 运行测试
make clean         # 清理容器和数据卷
```

## 📚 学习路径

### 1. 理解 Prometheus 指标类型

- 查看 `main.go` 中的指标定义
- 访问 `/metrics` 端点查看原始指标格式

### 2. 探索 Prometheus

- 打开 <http://localhost:9090>
- 尝试 PromQL 查询：

  ```promql
  rate(http_requests_total[1m])
  histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))
  ```

### 3. 自定义 Grafana 仪表板

- 登录 Grafana (admin/admin)
- 编辑现有仪表板或创建新的
- 尝试不同的可视化类型

### 4. 扩展应用

- 添加新的指标
- 创建告警规则
- 实现更复杂的业务场景

## 🛠️ 常见问题

**Q: Grafana 无法连接到 Prometheus？**
A: 确保所有服务都在同一个 Docker 网络中，检查 `docker-compose.yml` 配置。

**Q: 指标没有显示？**
A:

1. 检查 Go 应用是否正常运行
2. 访问 Prometheus 的 Targets 页面查看抓取状态
3. 确保访问了应用端点以生成指标数据

**Q: 如何重置 Grafana 密码？**
A: 删除 Grafana 数据卷后重启：

```bash
docker-compose down -v
docker-compose up -d
```

## 🔗 相关资源

- [Prometheus 官方文档](https://prometheus.io/docs/)
- [Grafana 官方文档](https://grafana.com/docs/)
- [Prometheus Go 客户端](https://github.com/prometheus/client_golang)
- [PromQL 教程](https://prometheus.io/docs/prometheus/latest/querying/basics/)

## 📝 许可证

MIT License

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！
