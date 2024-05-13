package env

import (
	"errors"
	"github.com/semho/chat-microservices/chat-server/internal/config"
	"net"
	"os"
)

const (
	authGrpcHostEnvName = "AUTH_GRPC_HOST"
	authGrpcPortEnvName = "AUTH_GRPC_PORT"
)

type clientConfig struct {
	authHost string
	authPort string
}

func NewClientConfig() (config.ClientConfig, error) {
	authHost := os.Getenv(authGrpcHostEnvName)
	if len(authHost) == 0 {
		return nil, errors.New("auth grpc host not found")
	}
	authPort := os.Getenv(authGrpcPortEnvName)
	if len(authPort) == 0 {
		return nil, errors.New("auth grpc port not found")
	}

	return &clientConfig{
		authHost: authHost,
		authPort: authPort,
	}, nil
}

func (cfg *clientConfig) Address() string {
	return net.JoinHostPort(cfg.authHost, cfg.authPort)
}
