package server

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"log/slog"
	"strconv"
	"time"
)

var metrics = promauto.NewSummaryVec(prometheus.SummaryOpts{
	Namespace:  "auth",
	Subsystem:  "http",
	Name:       "request",
	Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
},
	[]string{"status"},
)

func observe(d time.Duration, status int) {
	slog.Debug("test")
	metrics.WithLabelValues(strconv.Itoa(status)).Observe(d.Seconds())
}
