package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/push"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var addr = flag.String("listen-address", ":8081", "listen http req.")

var (
	c = promauto.NewCounter(prometheus.CounterOpts{
		Name: "app_sample_metric",
		Help: "Sample metric",
	})

	h = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "app_sample_histogram",
		Help: "Sample histogram",
	})

	d = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "app_sample_devices",
		Help: "Sample counter opts devices"}, []string{"device"})

	e = promauto.NewCounter(prometheus.CounterOpts{
		Name: "app_push_metric",
		Help: "Sample metric(push)",
	})
)

func main() {

	go func() {
		for {
			rand.Seed(time.Now().UnixNano())
			h.Observe(float64(rand.Intn(100-0+1) + 0))
			d.With(prometheus.Labels{"device": "/dev/sda"}).Inc()
			c.Inc()
			fmt.Print(".")
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		for {
			// Example of metric push
			err := push.New("http://pushgateway:9091", "app_job").Collector(e).Add()
			if err != nil {
				_ = fmt.Errorf("%v", err)
			}
			e.Inc()
			fmt.Print("_")
			time.Sleep(1 * time.Second)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}
