// Standard-library only: no external modules, so the container build downloads
// nothing (no go mod tidy, no proxy, no checksum db). It hand-writes the
// Prometheus text exposition format, including real go_* runtime metrics so it
// is still obviously a Go app when scraped.
package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"sync/atomic"
	"time"
)

const lang = "go"

var (
	requests int64
	inflight int64
)

func main() {
	// Simulated workload so the graphs move during the demo.
	go func() {
		for {
			atomic.AddInt64(&requests, int64(rand.Intn(5)+1))
			atomic.StoreInt64(&inflight, int64(rand.Intn(20)))
			time.Sleep(2 * time.Second)
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintf(w, "Demo app: %s\n", lang)
	})

	http.HandleFunc("/metrics", func(w http.ResponseWriter, _ *http.Request) {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")

		// Identity + simulated workload. Named "demo_app" (no "_info" suffix) to
		// stay consistent with the Spring Boot app, whose Prometheus client strips
		// the reserved "_info" suffix from gauges.
		fmt.Fprint(w, "# HELP demo_app Demo app identity\n# TYPE demo_app gauge\n")
		fmt.Fprintf(w, "demo_app{language=%q,app=\"app-go\"} 1\n", lang)

		fmt.Fprint(w, "# HELP demo_requests_total Simulated processed requests\n# TYPE demo_requests_total counter\n")
		fmt.Fprintf(w, "demo_requests_total{language=%q} %d\n", lang, atomic.LoadInt64(&requests))

		fmt.Fprint(w, "# HELP demo_inflight_requests Simulated in-flight requests\n# TYPE demo_inflight_requests gauge\n")
		fmt.Fprintf(w, "demo_inflight_requests{language=%q} %d\n", lang, atomic.LoadInt64(&inflight))

		// Go runtime metrics (same names client_golang emits) -> obviously Go.
		fmt.Fprint(w, "# HELP go_info Information about the Go environment.\n# TYPE go_info gauge\n")
		fmt.Fprintf(w, "go_info{version=%q} 1\n", runtime.Version())

		fmt.Fprint(w, "# HELP go_goroutines Number of goroutines that currently exist.\n# TYPE go_goroutines gauge\n")
		fmt.Fprintf(w, "go_goroutines %d\n", runtime.NumGoroutine())

		fmt.Fprint(w, "# HELP go_memstats_alloc_bytes Number of bytes allocated and still in use.\n# TYPE go_memstats_alloc_bytes gauge\n")
		fmt.Fprintf(w, "go_memstats_alloc_bytes %d\n", m.Alloc)

		fmt.Fprint(w, "# HELP go_memstats_heap_objects Number of allocated objects.\n# TYPE go_memstats_heap_objects gauge\n")
		fmt.Fprintf(w, "go_memstats_heap_objects %d\n", m.HeapObjects)
	})

	_ = http.ListenAndServe(":2112", nil)
}
