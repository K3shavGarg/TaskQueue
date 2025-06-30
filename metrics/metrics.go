package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	JobsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "jobs_total",
			Help: "Total number of jobs submitted",
		},
	)
	JobsProcessed = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "jobs_processed_total",
			Help: "Total number of jobs processed",
		},
	)
	JobsFailed = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "jobs_failed_total",
			Help: "Total number of jobs failed",
		},
	)
	JobDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "job_processing_duration_seconds",
			Help:    "Duration of job processing in seconds",
			Buckets: prometheus.DefBuckets,
		},
	)
)

func Init() {
	prometheus.MustRegister(JobsTotal, JobsProcessed, JobsFailed, JobDuration)
}
