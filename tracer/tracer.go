package tracer

import (
	"io"

	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log/zap"
	"github.com/uber/jaeger-lib/metrics/prometheus"
	"go.uber.org/zap"
)

func SetupTracer(service string) (io.Closer, error) {
	cfg, err := jaegercfg.FromEnv()
	if err != nil {
		return nil, err
	}

	logger := jaegerlog.NewLogger(zap.L())
	metricsFactory := prometheus.New()

	// Initialize tracer with a logger and a metrics factory
	closer, err := cfg.InitGlobalTracer(
		service,
		jaegercfg.Logger(logger),
		jaegercfg.Metrics(metricsFactory),
	)
	if err != nil {
		return nil, err
	}

	return closer, nil
}
