package main

import (
	"flag"
	"fmt"
	"github.com/ycrxun/onion/registry"
	"github.com/ycrxun/onion/services/account/server"
	"github.com/ycrxun/onion/services/account/storage"
	"github.com/ycrxun/onion/tracing"
	"os"
	"github.com/ycrxun/onion/metrics"
	"github.com/ycrxun/onion/metrics/prome"
	"github.com/prometheus/client_golang/prometheus"
)

func run(port int, consul *registry.Client, jaegeraddr string, storage storage.Storage) error {

	tracer, err := tracing.Init("account", jaegeraddr)
	if err != nil {
		return fmt.Errorf("tracing init error: %v", err)
	}

	id, err := consul.Register("account", port)
	defer consul.Deregister(id)
	if err != nil {
		return fmt.Errorf("failed to register service: %v", err)
	}

	fieldKeys := []string{"method"}

	m := metrics.Metrics{
		Counter: prome.NewCounterFrom(prometheus.CounterOpts{
			Namespace: "api",
			Subsystem: "account_service",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, fieldKeys),
		Histogram: prome.NewHistogramFrom(prometheus.HistogramOpts{
			Namespace: "api",
			Subsystem: "account_service",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys),
	}

	srv := server.NewServer(tracer, m, storage)

	return srv.Run(port)
}

func main() {
	var (
		port       = flag.Int("port", 9527, "The server port")
		jaegeraddr = flag.String("jaeger_addr", "jaeger:6831", "Jaeger address")
		consuladdr = flag.String("consul_addr", "consul:8500", "Consul address")
	)
	flag.Parse()
	if len(os.Args) < 2 {
		os.Exit(1)
	}

	consul, err := registry.NewClient(*consuladdr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	storage := storage.NewMemoryStorage()

	if err := run(*port, consul, *jaegeraddr, storage); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
