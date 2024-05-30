package env

import (
	"errors"
	"github.com/semho/chat-microservices/chat-server/internal/config"
	"net"
	"os"
)

const (
	jaegerHostEnvName = "JAEGER_HOST"
	jaegerPortEnvName = "JAEGER_PORT"
)

type jaegerConfig struct {
	host string
	port string
}

func NewJaegerConfig() (config.ClientConfig, error) {
	host := os.Getenv(jaegerHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("jaeger host not found")
	}
	port := os.Getenv(jaegerPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("jaeger port not found")
	}

	return &jaegerConfig{
		host: host,
		port: port,
	}, nil
}

func (cfg *jaegerConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}
