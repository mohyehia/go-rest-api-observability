package core

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

type ApplicationMetrics struct {
	GetPostsTotal     prometheus.Counter
	GetPostByIDTotal  prometheus.Counter
	HttpRequestsTotal *prometheus.CounterVec
}

func NewApplicationMetrics(registry *prometheus.Registry) *ApplicationMetrics {
	// 1. Automatically append the standard system/runtime metrics
	registry.MustRegister(collectors.NewGoCollector(), collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	// 2. Initialize your custom domain metric
	appMetrics := &ApplicationMetrics{
		GetPostsTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "api_get_posts_total",
			Help: "The total number of get posts APIs",
		}),
		GetPostByIDTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "api_get_post_by_id_total",
			Help: "The total number of get post by id APIs",
		}),
		HttpRequestsTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "api_http_requests_total",
			Help: "The total number of http requests APIs processed by status and path",
		}, []string{"path", "status"}),
	}
	registry.MustRegister(appMetrics.GetPostsTotal, appMetrics.GetPostByIDTotal, appMetrics.HttpRequestsTotal)
	return appMetrics
}
