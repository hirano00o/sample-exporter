package main

import (
	"flag"
	"log"
	"net/http"
	"sample-exporter/collector"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	addr = flag.String("listen address", ":9090", "The Address to listen on for HTTP Requests.")
)

func main() {
	flag.Parse()

	c, err := collector.NewSampleCollector()
	if err != nil {
		log.Fatal(err)
	}
	prometheus.MustRegister(c)

	http.Handle("/metrics", promhttp.Handler())

	log.Println("Listening on ", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal(err)
	}
}
