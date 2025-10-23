package metrics

import (
	"net/http"
	"time"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequestsTotal *prometheus.CounterVec
	httpDurationHist  *prometheus.HistogramVec
	inited            bool
)

func Init() {
	if inited {
		return
	}

	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests labeled by method, route, and status.",
		},
		[]string{"method", "route", "status"},
	)

	httpDurationHist = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latency in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "route"},
	)

	inited = true
}

func Handler() fiber.Handler {
	return adaptor.HTTPHandler(promhttp.Handler())
}

func Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()

		if inited {
			method := string(c.Method())
			route := c.Route().Path
			status := c.Response().StatusCode()
			httpRequestsTotal.WithLabelValues(method, route, http.StatusText(status)).Inc()
			httpDurationHist.WithLabelValues(method, route).Observe(time.Since(start).Seconds())
		}

		return err
	}
}
