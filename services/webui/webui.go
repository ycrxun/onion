package main

import (
	"flag"
	"os"
	"fmt"
	"github.com/ycrxun/onion/registry"
	"github.com/ycrxun/onion/tracing"
	"github.com/ycrxun/onion/discovery"
	"github.com/ycrxun/onion/services/webui/api"
	"github.com/ycrxun/onion/services/account/proto"
)

func main() {
	var (
		port       = flag.Int("port", 9528, "The server port")
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

	tracer, err := tracing.Init("webui", *jaegeraddr)
	if err != nil {
		fmt.Errorf("tracing init error: %v", err)
		os.Exit(1)
	}

	ac, err := discovery.Discovery(
		"account",
		discovery.WithTracer(tracer),
		discovery.WithBalancer(consul.Client),
	)

	if err != nil {
		fmt.Errorf("discovery account error: %v", err)
		os.Exit(1)
	}

	srv := api.NewServer(tracer, account.NewAccountServiceClient(ac))

	srv.Run(*port)
}
