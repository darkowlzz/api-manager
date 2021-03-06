package sharedvolume

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

// LatencyMetric observes latency.
type LatencyMetric interface {
	Observe(latency time.Duration)
}

var (
	// ReconcileDuration is the latency metric that measures the duration of the
	// shared volume reconcile loop.
	ReconcileDuration LatencyMetric = &latencyAdapter{m: reconcileLatencyHistogram}

	// registerMetricsOnce keeps track of metrics registration.
	registerMetricsOnce sync.Once
)

var (
	reconcileLatencyHistogram = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "storageos_shared_volume_reconcile_duration_seconds",
			Help:    "Distribution of the length of time to reconcile all shared volumes.",
			Buckets: prometheus.DefBuckets,
		},
	)
)

// RegisterMetrics ensures that the package metrics are registered.
func RegisterMetrics() {
	registerMetricsOnce.Do(func() {
		metrics.Registry.MustRegister(reconcileLatencyHistogram)
	})
}

type latencyAdapter struct {
	m prometheus.Histogram
}

func (l *latencyAdapter) Observe(latency time.Duration) {
	l.m.Observe(latency.Seconds())
}
