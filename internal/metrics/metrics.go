package metrics

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	HTTP *HTTPMetrics
}

type HTTPMetrics struct {
	Requests *prometheus.CounterVec
	Duration *prometheus.HistogramVec
}

func New(reg prometheus.Registerer) *Metrics {
	m := &Metrics{
		HTTP: newHTTPMetrics(),
	}

	reg.MustRegister(
		m.HTTP.Requests,
		m.HTTP.Duration,
	)

	return m
}

func newHTTPMetrics() *HTTPMetrics {
	return &HTTPMetrics{
		Requests: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "nexus",
				Subsystem: "http",
				Name:      "requests_total",
				Help:      "Total number of HTTP requests.",
			},
			[]string{
				"method",
				"path",
				"status",
			},
		),

		Duration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: "nexus",
				Subsystem: "http",
				Name:      "request_duration_seconds",
				Help:      "HTTP request latency in seconds.",

				// Good defaults for HTTP APIs
				Buckets: []float64{
					0.005,
					0.01,
					0.025,
					0.05,
					0.1,
					0.25,
					0.5,
					1,
					2.5,
					5,
					10,
					30,
				},
			},
			[]string{
				"method",
				"path",
			},
		),
	}
}
