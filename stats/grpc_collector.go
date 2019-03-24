package stats

import (
	prom "github.com/prometheus/client_golang/prometheus"
	"time"
)

type GrpcCollector struct {
	namespace string
	//callingGauge	*prom.GaugeVec	//调用中的数量
	waitGauge    *prom.GaugeVec   //等待中的数量 rpc
	startCounter *prom.CounterVec //开始rpc计数
	endCounter   *prom.CounterVec //结束rpc计数
	errorCounter *prom.CounterVec //执行出错rpc计数

	handledHistogram *prom.HistogramVec //rpc处理耗时直方图
}

func NewGrpcCollector() *GrpcCollector {
	namespace := "grpc"
	return &GrpcCollector{
		namespace: namespace,
		waitGauge: prom.NewGaugeVec(prom.GaugeOpts{
			Namespace: namespace,
			Subsystem: "",
			Name:      "grpc_wait_total",
			Help:      "grpc_wait_total",
		}, []string{"service", "service_id", "func_name"}),
		startCounter: prom.NewCounterVec(prom.CounterOpts{
			Namespace: namespace,
			Subsystem: "",
			Name:      "grpc_start_counter",
			Help:      "grpc_start_counter",
		}, []string{"service", "service_id", "func_name"}),
		endCounter: prom.NewCounterVec(prom.CounterOpts{
			Namespace: namespace,
			Subsystem: "",
			Name:      "grpc_end_counter",
			Help:      "grpc_end_counter",
		}, []string{"service", "service_id", "func_name"}),
		errorCounter: prom.NewCounterVec(prom.CounterOpts{
			Namespace: namespace,
			Subsystem: "",
			Name:      "grpc_error_counter",
			Help:      "grpc_error_counter",
		}, []string{"service", "service_id", "func_name"}),
		handledHistogram: prom.NewHistogramVec(prom.HistogramOpts{
			Namespace: namespace,
			Subsystem: "",
			Name:      "grpc_handled_histogram",
			Help:      "grpc_handled_histogram handled time",
			//Buckets:   prom.DefBuckets,
			//TODO Buckets 优化
			Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
			//Buckets: prom.LinearBuckets(0.1, .1, 20),
		}, []string{"service", "service_id", "func_name", "status"}),
	}
}

func (g *GrpcCollector) Describe(ch chan<- *prom.Desc) {
	g.waitGauge.Describe(ch)
	g.startCounter.Describe(ch)
	g.endCounter.Describe(ch)
	g.errorCounter.Describe(ch)
	g.handledHistogram.Describe(ch)
}

func (g *GrpcCollector) Collect(ch chan<- prom.Metric) {
	g.waitGauge.Collect(ch)
	g.startCounter.Collect(ch)
	g.endCounter.Collect(ch)
	g.errorCounter.Collect(ch)
	g.handledHistogram.Collect(ch)
}

func (g *GrpcCollector) waitGaugeAdd(service string, service_id string, func_name string, counter float64) {
	g.waitGauge.WithLabelValues(service, service_id, func_name).Add(counter)
}

func (g *GrpcCollector) waitGaugeInc(service string, service_id string, func_name string) {
	g.waitGauge.WithLabelValues(service, service_id, func_name).Inc()
}

func (g *GrpcCollector) waitGaugeDec(service string, service_id string, func_name string) {
	g.waitGauge.WithLabelValues(service, service_id, func_name).Dec()
}

func (g *GrpcCollector) startCounterInc(service string, service_id string, func_name string) {
	g.startCounter.WithLabelValues(service, service_id, func_name).Inc()
}

func (g *GrpcCollector) startCounterAdd(service string, service_id string, func_name string, counter float64) {
	g.startCounter.WithLabelValues(service, service_id, func_name).Add(counter)
}

func (g *GrpcCollector) endCounterInc(service string, service_id string, func_name string) {
	g.endCounter.WithLabelValues(service, service_id, func_name).Inc()
}

func (g *GrpcCollector) endCounterAdd(service string, service_id string, func_name string, counter float64) {
	g.endCounter.WithLabelValues(service, service_id, func_name).Add(counter)
}
func (g *GrpcCollector) errorCounterInc(service string, service_id string, func_name string) {
	g.errorCounter.WithLabelValues(service, service_id, func_name).Inc()
}

func (g *GrpcCollector) errorCounterAdd(service string, service_id string, func_name string, counter float64) {
	g.errorCounter.WithLabelValues(service, service_id, func_name).Add(counter)
}

//Observe time set Seconds
func (g *GrpcCollector) handledTime(service string, service_id string, func_name string, status string, start time.Time) {
	g.handledHistogram.WithLabelValues(service, service_id, func_name, status).Observe(time.Since(start).Seconds())
}

type NsqCollector struct {
	namespace string
	*nsqVec
}

type nsqVec struct {
	errorPiDefaultGrpcCollectorounter *prom.CounterVec //Ping失败计数
	waitPubGauge                      *prom.GaugeVec   //等待中的发布数量
	waitSubGauge                      *prom.GaugeVec   //等待中的订阅数量

	startSubCounter     *prom.CounterVec
	endSubCounter       *prom.CounterVec
	errorSubCounter     *prom.CounterVec
	handledSubHistogram *prom.HistogramVec //Sub处理耗时直方图
}

func NewNsqCollector() *NsqCollector {
	namespace := "nsq"
	nsqVec := &nsqVec{
		errorPiDefaultGrpcCollectorounter: prom.NewCounterVec(prom.CounterOpts{
			Namespace: namespace,
			Subsystem: "",
			Name:      "nsq_ping_error_counter",
			Help:      "nsq_ping_error_counter",
		}, []string{"service", "topic"}),
		waitPubGauge: prom.NewGaugeVec(prom.GaugeOpts{
			Namespace: namespace,
			Subsystem: "",
			Name:      "nsq_waitPub_gauge",
			Help:      "nsq_waitPub_gauge",
		}, []string{"service", "topic"}),
		waitSubGauge: prom.NewGaugeVec(prom.GaugeOpts{
			Namespace: namespace,
			Subsystem: "",
			Name:      "nsq_waitSub_gauge",
			Help:      "nsq_waitSub_gauge",
		}, []string{"service", "topic"}),
		startSubCounter: prom.NewCounterVec(prom.CounterOpts{
			Namespace: namespace,
			Subsystem: "",
			Name:      "nsq_startSub_counter",
			Help:      "nsq_startSub_counter",
		}, []string{"service", "topic", "channel"}),
		endSubCounter: prom.NewCounterVec(prom.CounterOpts{
			Namespace: namespace,
			Subsystem: "",
			Name:      "nsq_endSub_counter",
			Help:      "nsq_endSub_counter",
		}, []string{"service", "topic", "channel"}),
		errorSubCounter: prom.NewCounterVec(prom.CounterOpts{
			Namespace: namespace,
			Subsystem: "",
			Name:      "nsq_errorSub_counter",
			Help:      "nsq_errorSub_counter",
		}, []string{"service", "topic", "channel"}),
		handledSubHistogram: prom.NewHistogramVec(prom.HistogramOpts{
			Namespace: namespace,
			Subsystem: "",
			Name:      "nsq_handledSub_histogram",
			Help:      "nsq_handledSub_histogram handled time",
			//Buckets:   prom.DefBuckets,
			//TODO Buckets 优化
			Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
			//Buckets: prom.LinearBuckets(0.1, .1, 20),
		}, []string{"service", "topic", "channel", "status"}),
	}
	return &NsqCollector{
		namespace: namespace,
		nsqVec:    nsqVec,
	}
}

func (n *NsqCollector) Describe(ch chan<- *prom.Desc) {
	n.errorPiDefaultGrpcCollectorounter.Describe(ch)
	n.waitPubGauge.Describe(ch)
	n.waitSubGauge.Describe(ch)
	n.startSubCounter.Describe(ch)
	n.endSubCounter.Describe(ch)
	n.errorSubCounter.Describe(ch)
	n.handledSubHistogram.Describe(ch)
}

func (n *NsqCollector) Collect(ch chan<- prom.Metric) {
	n.errorPiDefaultGrpcCollectorounter.Collect(ch)
	n.waitPubGauge.Collect(ch)
	n.waitSubGauge.Collect(ch)
	n.startSubCounter.Collect(ch)
	n.endSubCounter.Collect(ch)
	n.errorSubCounter.Collect(ch)
	n.handledSubHistogram.Collect(ch)
}

//Observe time set Seconds
func (n *NsqCollector) handledTime(service string, topic string, status string, start time.Time) {
	n.handledSubHistogram.WithLabelValues(service, topic, status).Observe(time.Since(start).Seconds())
}
