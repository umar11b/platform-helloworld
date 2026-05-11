package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type HelloHandler struct {
	requestsTotal    *prometheus.CounterVec
	requestDuration  *prometheus.HistogramVec
}

func NewHelloHandler(reg prometheus.Registerer) *HelloHandler {
	requestsTotal := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "helloworld_requests_total",
		Help: "Total number of requests to the hello handler.",
	}, []string{"handler", "status_code"})

	requestDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "helloworld_request_duration_seconds",
		Help:    "Latency of requests to the hello handler in seconds.",
		Buckets: prometheus.DefBuckets,
	}, []string{"handler", "status_code"})

	reg.MustRegister(requestsTotal, requestDuration)

	return &HelloHandler{
		requestsTotal:   requestsTotal,
		requestDuration: requestDuration,
	}
}

func (h *HelloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	secret := os.Getenv("HELLO_SECRET")
	statusCode := http.StatusOK

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"message": secret})

	duration := time.Since(start).Seconds()
	code := strconv.Itoa(statusCode)
	h.requestsTotal.WithLabelValues("hello", code).Inc()
	h.requestDuration.WithLabelValues("hello", code).Observe(duration)
}
