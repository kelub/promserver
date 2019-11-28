package stats

var (
	grpcStatas = NewPromVec("grpc").
		Gauge("wait_total", []string{"service", "service_id", "method"}).
		Counter("wait_state", []string{"service", "service_id", "method", "state"}).
		Histogram("handled", []string{"service", "service_id", "method", "state"}, []float64{.01, .05, .1, 1, 5, 10})
)
