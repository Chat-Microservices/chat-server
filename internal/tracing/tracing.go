package tracing

import (
	"github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"
)

func Init(logger *zap.Logger, serviceName, address string) {
	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LocalAgentHostPort: address,
		},
	}

	_, err := cfg.InitGlobalTracer(serviceName)
	if err != nil {
		logger.Fatal("failed to init tracing", zap.Error(err))
	} else {
		logger.Info("init tracing", zap.String("service", serviceName))
	}
}
