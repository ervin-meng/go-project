package tracer

import (
	"context"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"google.golang.org/grpc"
)

func Init(serviceName string) {
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{ //采样
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: "127.0.0.1:6831",
		},
		ServiceName: serviceName,
	}

	tracer, _, _ := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))

	//traceCloser = closer

	opentracing.SetGlobalTracer(tracer)
}

func OpenTracingGRPCServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		if info.FullMethod != "/grpc.health.v1.Health/Check" {
			return otgrpc.OpenTracingServerInterceptor(opentracing.GlobalTracer())(ctx, req, info, handler)
		} else {
			return handler(ctx, req)
		}
	}
}
