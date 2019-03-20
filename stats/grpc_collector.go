package stats

import (
	prom "github.com/prometheus/client_golang/prometheus"
)

type GrpcCollector struct {
	namespace string
	*grpcDesc
}

type grpcDesc struct {
	//callingGauge	*prom.GaugeVec	//调用中的数量
	waitGauge    *prom.GaugeVec   //等待中的数量
	startCounter *prom.CounterVec //开始rpc计数
	endCounter   *prom.CounterVec //结束rpc计数
	errorCounter *prom.CounterVec //执行出错rpc计数

	handledHistogram *prom.HistogramVec //执行直方图
}

func NewGrpcCollector() *GrpcCollector {
	namespace := "grpc"
	grpcDesc := &grpcDesc{
		waitGauge: prom.NewGaugeVec(prom.GaugeOpts{
			Namespace: namespace,
			Subsystem: "",
			Name:      "grpc_wait_total",
			Help:      "grpc_wait_total",
		}, []string{"xxx"}),
		startCounter:     prom.NewCounterVec(),
		endCounter:       prom.NewCounterVec(),
		handledHistogram: prom.NewHistogramVec(),
	}
	return &GrpcCollector{
		namespace: namespace,
		grpcDesc:  grpcDesc,
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
