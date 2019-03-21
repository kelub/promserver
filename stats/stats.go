package stats

import (
	"context"
	"github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"net/http"
	"time"
)

type rpcMonitor struct {
	service    string
	name       string
	service_id string
	collector  *GrpcCollector
	startTime  time.Time
}

func NewRpcMonitor(service string, name string, service_id string) *rpcMonitor {
	logrus.Infof("NewRpcMonitor %s", name)
	m := &rpcMonitor{
		service:    service,
		name:       name,
		service_id: service_id,
	}
	m.startTime = time.Now()
	return m
}

func (m *rpcMonitor) WaitGaugeInc() {
	m.collector.waitGaugeInc(m.service, m.service_id, m.name)
}

func (m *rpcMonitor) WaitGaugeDec() {
	m.collector.waitGaugeDec(m.service, m.service_id, m.name)
}

func (m *rpcMonitor) WaitGaugeAdd(value float64) {
	m.collector.waitGaugeAdd(m.service, m.service_id, m.name, value)
}

func (m *rpcMonitor) StartCounterInc() {
	m.collector.startCounterInc(m.service, m.service_id, m.name)
}

func (m *rpcMonitor) EndCounterInc() {
	m.collector.startCounterInc(m.service, m.service_id, m.name)
}

func (m *rpcMonitor) ErrorCounterInc() {
	m.collector.startCounterInc(m.service, m.service_id, m.name)
}

// TODO 具体错误标示
func (m *rpcMonitor) HandledEnd(status string) {
	m.collector.handledTime(m.service, m.service_id, m.name, status, m.startTime)
}

// gRPC 拦截器
func UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	monitor := NewRpcMonitor("service_test", "77", info.FullMethod)
	monitor.WaitGaugeInc()
	monitor.StartCounterInc()
	resp, err = handler(ctx, req)
	monitor.EndCounterInc()
	monitor.WaitGaugeDec()
	if err != nil {
		monitor.ErrorCounterInc()
		monitor.HandledEnd("failed")
	}
	monitor.HandledEnd("succeed")
	return resp, err
}

type nsqMonitor struct {
	service   string
	topic     string
	channel   string
	collector *NsqCollector
	startTime time.Time
	nsqType   string
}

func NewNsqMonitor(service string, topic string, channel string, nsqType string) *nsqMonitor {
	logrus.Infof("NewNsqMonitor %s", topic)
	m := &nsqMonitor{
		service: service,
		topic:   topic,
		channel: channel,
		nsqType: nsqType,
	}
	m.startTime = time.Now()
	return m
}

// TODO 具体错误标示
func (n *nsqMonitor) HandledEnd(status string) {
	n.collector.handledTime(n.service, n.topic, status, n.startTime)
}

// 运行
func StartPromServer() {
	ListenAddr := ":9102"
	prometheus.MustRegister(NewGrpcCollector())
	prometheus.MustRegister(NewNsqCollector())
	http.Handle("/metrics", prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Nsq Exporter</title></head>
			<body>
			<h1>Node Exporter</h1>
			<p><a href="` + "/metrics" + `">Metrics</a></p>
			</body>
			</html>`))
	})
	logrus.Infof("Listening on %s", ListenAddr)
	if err := http.ListenAndServe(ListenAddr, nil); err != nil {
		logrus.WithError(err).Fatalln("启动失败", err)
	}
}
