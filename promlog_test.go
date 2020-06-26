package promlog_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/practo/klog/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/practo/promlog"
)

const (
	addr     string = ":8080"
	endpoint string = "/metrics"
	testSeverityLevel string = "INFO"
	testPrefix string = "promlog"
)

func getSeverityForLine(line string) string {
	severityChar := strings.Split(line, "")[0]
	switch severityChar {
	case strings.Split(klog.InfoSeverityLevel, "")[0]:
		return klog.InfoSeverityLevel
	case strings.Split(klog.WarningSeverityLevel, "")[0]:
		return klog.WarningSeverityLevel
	case strings.Split(klog.ErrorSeverityLevel, "")[0]:
		return klog.ErrorSeverityLevel
	}

	return ""
}

func TestExposeAndQueryKlogCounters(t *testing.T) {
	// Create Prometheus hook and configure klog to use it:
	hook := promlog.MustNewPrometheusHook(testPrefix, testSeverityLevel)
	klog.AddHook(hook)

	server := httpServePrometheusMetrics(t)

	lines := httpGetMetrics(t)
	validateResult(t, 0, countFor(t, klog.InfoSeverityLevel, lines))
	validateResult(t, 0, countFor(t, klog.WarningSeverityLevel, lines))
	validateResult(t, 0, countFor(t, klog.ErrorSeverityLevel, lines))

	klog.Info("this is at info level!")
	lines = httpGetMetrics(t)
	validateResult(t, 1, countFor(t, klog.InfoSeverityLevel, lines))
	validateResult(t, 0, countFor(t, klog.WarningSeverityLevel, lines))
	validateResult(t, 0, countFor(t, klog.ErrorSeverityLevel, lines))

	klog.Warning("this is at warning level!")
	lines = httpGetMetrics(t)
	validateResult(t, 1, countFor(t, klog.InfoSeverityLevel, lines))
	validateResult(t, 1, countFor(t, klog.WarningSeverityLevel, lines))
	validateResult(t, 0, countFor(t, klog.ErrorSeverityLevel, lines))

	klog.Error("this is at error level!")
	lines = httpGetMetrics(t)
	validateResult(t, 1, countFor(t, klog.InfoSeverityLevel, lines))
	validateResult(t, 1, countFor(t, klog.WarningSeverityLevel, lines))
	validateResult(t, 1, countFor(t, klog.ErrorSeverityLevel, lines))

	server.Close()
}

func TestInvalidSeverityLevel(t *testing.T) {
	_, err := promlog.NewPrometheusHook(testPrefix, "PANIC")
	if err == nil {
		t.Error("expected invalid severity error")
	}
}

// httpServePrometheusMetrics exposes the Prometheus metrics
// over HTTP, in a different go routine.
func httpServePrometheusMetrics(t *testing.T) *http.Server {
	server := &http.Server{
		Addr:    addr,
		Handler: promhttp.Handler(),
	}
	go server.ListenAndServe()
	return server
}

// httpGetMetrics queries the local HTTP server for the
// exposed metrics and parses the response.
func httpGetMetrics(t *testing.T) []string {
	resp, err := http.Get(fmt.Sprintf("http://localhost%v%v", addr, endpoint))
	if err != nil {
		t.Error(err)
		return []string{}
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return []string{}
	}
	lines := strings.Split(string(body), "\n")
	if len(lines) == 0 {
		t.Error("httpGetMetrics returned empty response")
		return []string{}
	}
	return lines
}

// countFor is a helper function to get
// the counter's value for the provided level.
func countFor(t *testing.T, severity string, lines []string) int {
	// Metrics are exposed as per the below example:
	//   # HELP test_debug Number of log statements at debug level.
	//   # TYPE test_debug counter
	//   test_debug 0
	metric := fmt.Sprintf(
		testPrefix + "log_messages_total{severity=\"%v\"}", severity)
	for _, line := range lines {
		items := strings.Split(line, " ")
		if len(items) != 2 { // e.g. {"test_debug", "0"}
			continue
		}
		if items[0] == metric {
			count, err := strconv.ParseInt(items[1], 10, 32)
			if err != nil {
				t.Errorf("error parsing line: %v\n", err)
			}
			return int(count)
		}
	}
	panic(fmt.Sprintf("Could not find %v in %v", metric, lines))
}

func validateResult(t *testing.T, wanted int, got int) {
	if wanted == got {
		return
	}
	t.Errorf("unexpected result: \n got:\t%d \nwant:\t%d", wanted, got)
}
