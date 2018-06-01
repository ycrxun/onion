package discovery

import (
	"fmt"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	consul "github.com/hashicorp/consul/api"
	lb "github.com/olivere/grpc/lb/consul"
	opentracing "github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)

// DiscoveryOption allows optional config for discovery
type DiscoveryOption func(name string) (grpc.DialOption, error)

// WithTracer traces rpc calls
func WithTracer(t opentracing.Tracer) DiscoveryOption {
	return func(name string) (grpc.DialOption, error) {
		return grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(t)), nil
	}
}

// WithBalancer enables client side load balancing
func WithBalancer(cc *consul.Client) DiscoveryOption {
	return func(name string) (grpc.DialOption, error) {
		r, err := lb.NewResolver(cc, name, "")
		if err != nil {
			return nil, err
		}
		return grpc.WithBalancer(grpc.RoundRobin(r)), nil
	}
}

// Discovery returns a load balanced grpc client conn with tracing interceptor
func Discovery(name string, opts ...DiscoveryOption) (*grpc.ClientConn, error) {
	dialopts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	for _, fn := range opts {
		opt, err := fn(name)
		if err != nil {
			return nil, fmt.Errorf("config error: %v", err)
		}
		dialopts = append(dialopts, opt)
	}

	conn, err := grpc.Dial(name, dialopts...)
	if err != nil {
		return nil, fmt.Errorf("failed to dial %s: %v", name, err)
	}

	return conn, nil
}
