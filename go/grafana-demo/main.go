package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// Counter: 计数器，只增不减
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"path", "method", "status"},
	)

	// Gauge: 仪表盘，可增可减
	activeConnections = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_connections",
			Help: "Number of active connections",
		},
	)

	// Histogram: 直方图，观察值的分布
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path"},
	)

	// Summary: 摘要，类似直方图但提供分位数
	responseSize = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "http_response_size_bytes",
			Help:       "Size of HTTP responses in bytes",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"path"},
	)

	// 自定义业务指标
	orderTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "orders_total",
			Help: "Total number of orders processed",
		},
	)

	orderAmount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "order_amount_current",
			Help: "Current order amount in dollars",
		},
	)
)

func init() {
	// 注册所有指标
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(activeConnections)
	prometheus.MustRegister(httpRequestDuration)
	prometheus.MustRegister(responseSize)
	prometheus.MustRegister(orderTotal)
	prometheus.MustRegister(orderAmount)
}

// 模拟业务逻辑，生成随机指标数据
func simulateBusinessMetrics() {
	go func() {
		for {
			// 模拟活跃连接数变化
			activeConnections.Set(float64(rand.Intn(100)))

			// 模拟订单处理
			if rand.Float64() > 0.7 {
				orderTotal.Inc()
				orderAmount.Set(float64(rand.Intn(10000)))
			}

			time.Sleep(2 * time.Second)
		}
	}()
}

// 记录 HTTP 请求的中间件
func metricsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// 调用实际的处理函数
		next(w, r)

		// 记录指标
		duration := time.Since(start).Seconds()
		httpRequestDuration.WithLabelValues(r.URL.Path).Observe(duration)
		httpRequestsTotal.WithLabelValues(r.URL.Path, r.Method, "200").Inc()
		responseSize.WithLabelValues(r.URL.Path).Observe(float64(rand.Intn(1000)))
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html>
<head>
    <title>Grafana Demo Application</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            max-width: 800px;
            margin: 50px auto;
            padding: 20px;
            background-color: #f5f5f5;
        }
        .container {
            background-color: white;
            padding: 30px;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        h1 {
            color: #333;
        }
        .endpoint {
            background-color: #f0f0f0;
            padding: 10px;
            margin: 10px 0;
            border-radius: 4px;
        }
        a {
            color: #0066cc;
            text-decoration: none;
        }
        a:hover {
            text-decoration: underline;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>🚀 Grafana + Prometheus + Go Demo</h1>
        <p>这是一个演示应用，展示如何使用 Go 暴露 Prometheus 指标。</p>
        
        <h2>可用端点：</h2>
        <div class="endpoint">
            <strong>GET /metrics</strong> - Prometheus 指标端点
            <br><a href="/metrics">查看指标</a>
        </div>
        <div class="endpoint">
            <strong>GET /api/data</strong> - 示例 API 端点
            <br><a href="/api/data">调用 API</a>
        </div>
        <div class="endpoint">
            <strong>GET /health</strong> - 健康检查端点
            <br><a href="/health">检查健康状态</a>
        </div>
        
        <h2>访问监控：</h2>
        <ul>
            <li><a href="http://localhost:9090" target="_blank">Prometheus</a> - 时序数据库</li>
            <li><a href="http://localhost:3000" target="_blank">Grafana</a> - 可视化仪表板 (admin/admin)</li>
        </ul>
    </div>
</body>
</html>`
	fmt.Fprint(w, html)
}

func apiDataHandler(w http.ResponseWriter, r *http.Request) {
	// 模拟一些处理延迟
	delay := time.Duration(rand.Intn(500)) * time.Millisecond
	time.Sleep(delay)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status":"success","timestamp":"%s","message":"数据已返回"}`, time.Now().Format(time.RFC3339))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status":"healthy","timestamp":"%s"}`, time.Now().Format(time.RFC3339))
}

func main() {
	// 启动模拟业务指标生成
	simulateBusinessMetrics()

	// 设置路由
	http.HandleFunc("/", metricsMiddleware(homeHandler))
	http.HandleFunc("/api/data", metricsMiddleware(apiDataHandler))
	http.HandleFunc("/health", metricsMiddleware(healthHandler))

	// Prometheus metrics endpoint
	http.Handle("/metrics", promhttp.Handler())

	port := ":8080"
	log.Printf("🚀 服务器启动在 http://localhost%s", port)
	log.Printf("📊 Prometheus 指标: http://localhost%s/metrics", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
