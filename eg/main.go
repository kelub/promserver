package main

import (
	"errors"
	"github.com/Sirupsen/logrus"
	"math/rand"
	"promserver/stats"
	"time"
)

func main() {

	go stats.StartPromServer()
	time.Sleep(1 * time.Second)
	endtime := time.NewTimer(30 * time.Millisecond)
	timed := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-timed.C:
			Interceptor()
		case <-endtime.C:
			break
		}
	}
}

func nsqtest() {
	logrus.Infof("nsqtest  in")
	monitor := stats.NewNsqMonitor(stats.DefaultNsqCollector, "service_test", "77", "test", "Sub")
	logrus.Infof("monitor  %v", monitor)
	err := dothing()
	if err != nil {
		monitor.HandledEnd("failed")
	}
	monitor.HandledEnd("succeed")
	logrus.Infof("nsqtest  out")
}

// grpc
func Interceptor() {
	logrus.Infof("Interceptor  in")
	monitor := stats.NewRpcMonitor(stats.DefaultGrpcCollector, "service_test", "77", "test")
	logrus.Infof("monitor  %v", monitor)
	monitor.WaitGaugeInc()
	monitor.StartCounterInc()
	err := dothing()
	monitor.EndCounterInc()
	monitor.WaitGaugeDec()
	if err != nil {
		monitor.ErrorCounterInc()
		monitor.HandledEnd("failed")
	}
	monitor.HandledEnd("succeed")
	logrus.Infof("Interceptor  out")
}

func dothing() error {
	i := rand.Intn(5)
	logrus.Infof("Sleep %d", i)
	time.Sleep(time.Duration(i) * time.Second)
	if i == 1 {
		return errors.New("error test")
	}
	return nil
}
