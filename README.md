<p align="left">
	<a href="https://github.com/practo/promlog/releases/latest">
		<img src="https://img.shields.io/github/release/practo/promlog.svg"/>
	</a>
	<a href="https://travis-ci.org/practo/promlog">
		<img src="https://img.shields.io/travis/practo/promlog.svg"/>
	</a>
	<a href="https://coveralls.io/github/practo/promlog?branch=master">
		<img src="https://img.shields.io/coveralls/practo/promlog.svg"/>
	</a>
	<a href="https://goreportcard.com/report/github.com/practo/promlog">
		<img src="https://goreportcard.com/badge/github.com/practo/promlog"/>
	</a>
	<a href="LICENSE">
		<img src="https://img.shields.io/badge/license-Apache%202.0-blue.svg"/>
	</a>
</p>

# promlog
klog hook to expose the number of log messages as Prometheus metrics:
```
log_messages{level="INFO"}
log_messages{level="WARNING"}
log_messages{level="ERROR"}
log_messages{level="FATAL"}
```

## Usage

Sample code:
```go
package main

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "k8s.io/klog/v2"
	"github.com/practo/promlog"
)

func main() {
	// Create the Prometheus hook:
	hook := promlog.MustNewPrometheusHook()

	// Configure klog to use the Prometheus hook:
	log.AddHook(hook)

	// Expose Prometheus metrics via HTTP, as you usually would:
	go http.ListenAndServe(":8080", promhttp.Handler())

	// Log with klog, as you usually would.
	// Every time the program generates a log message, a Prometheus counter is incremented for the corresponding level.
	for {
		log.Infof("foo")
		time.Sleep(1 * time.Second)
	}
}
```

Run the above program:
```
$ go run main.go
to fill
```

Scrape the Prometheus metrics exposed by the hook:
```
$ curl -fsS localhost:8080 | grep log_messages
to fill
```

## Compile
```
$ go build
```

## Test
```
$ go test

```
