package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/real-web-world/go-web-api/global"
	"github.com/real-web-world/go-web-api/pkg/gin"
)

var (
	httpReqTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_total",
			Help: "每个接口调用的次数",
		},
		[]string{"path"},
	)
	httpReqInFlight = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_request_in_flight",
			Help: "当前正在处理的请求",
		},
	)
	httpReqDurationMillisecond = promauto.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "http_request_duration_ms",
			Help: "每个请求的耗时",
			Objectives: map[float64]float64{0.5: 0.05, 0.8: 0.01, 0.9: 0.01, 0.95: 0.001,
				0.99: 0.001},
			// Buckets: []float64{1, 5, 10, 20, 50, 100, 200, 500, 1000, 5000,10000},
		},
		[]string{"path"},
	)
)

func Prometheus(c *gin.Context) {
	app := ginApp.GetApp(c)
	startProcTime := app.GetProcBeginTime()
	path := c.Request.URL.Path
	isDevApi := path == global.PrometheusApi ||
		strings.Index(path, global.DebugApiPrefix) == 0
	if !isDevApi {
		go httpReqInFlight.Inc()
		defer httpReqInFlight.Dec()
	}
	c.Next()
	if c.Writer.Status() == http.StatusNotFound || isDevApi {
		return
	}
	httpReqTotal.WithLabelValues(path).Inc()
	httpReqDurationMillisecond.WithLabelValues(path).
		Observe(float64(time.Since(*startProcTime).Milliseconds()))
}
