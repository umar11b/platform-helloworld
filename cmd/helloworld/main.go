package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/umar/helloworld/internal/handlers"
)

func main() {
	reg := prometheus.NewRegistry()
	reg.MustRegister(
		prometheus.NewGoCollector(),
		prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
	)

	helloHandler := handlers.NewHelloHandler(reg)

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", handlers.Healthz)
	mux.HandleFunc("/readyz", handlers.Readyz)
	mux.Handle("/hello", helloHandler)
	mux.HandleFunc("/fail", helloHandler.HandleFail)
	mux.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	log.Println("starting helloworld on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("server exited: %v", err)
	}
}
