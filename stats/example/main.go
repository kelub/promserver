package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"kelub/promserver/stats"
	"net/http"
	"time"
)

var (
	testStats = stats.NewPromVec("test").
		Gauge("run_total", []string{"service", "method"}).
		Counter("run_state", []string{"service", "method", "state"}).
		Histogram("handled", []string{"service", "method", "state"}, []float64{.01, .05, .1, 1, 5, 10})
)

func main() {
	go Run()
	StartPromServer()
}

func Run() {
	gTick := time.NewTicker(2 * time.Second)
	cTick := time.NewTicker(5 * time.Second)
	defer gTick.Stop()
	defer cTick.Stop()
	for {
		select {
		case <-gTick.C:
			testStats.GaugeAdd([]string{"service_a", "method_a"}, 1)
		case <-cTick.C:
			testStats.CounterAdd([]string{"service_a", "method_a", "true"}, 1)
		default:
			now := time.Now()
			time.Sleep(3 * time.Second)
			testStats.HandleTime([]string{"service_a", "method_a", "true"}, now)
		}
	}
}

func StartPromServer() error {
	ListenAddr := ":9102"
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Test Exporter</title></head>
			<body>
			<h1>Test Exporter</h1>
			<p><a href="` + "/metrics" + `">Metrics</a></p>
			</body>
			</html>`))
	})
	logrus.Infof("Listening on %s", ListenAddr)
	if err := http.ListenAndServe(ListenAddr, nil); err != nil {
		logrus.WithError(err).Fatalln("启动失败", err)
		return err
	}
	return nil
}
