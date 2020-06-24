package main

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/practo/klog/v2"
	"github.com/practo/promlog"
)

func main() {
	// Create the Prometheus hook:
	hook := promlog.MustNewPrometheusHook("")

	// Configure klog to use the Prometheus hook:
	klog.AddHook(hook)

	// Expose Prometheus metrics via HTTP, as you usually would:
	go http.ListenAndServe(":8080", promhttp.Handler())

	// Log with klog, as you usually would.
	// Every time the program generates a log message,
	// a Prometheus counter is incremented for the corresponding level.
	for {
		klog.Infof("foo")
		time.Sleep(1 * time.Second)
	}
}
