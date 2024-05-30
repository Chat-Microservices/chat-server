package tests

import "github.com/semho/chat-microservices/chat-server/internal/logger"

func initLogger() {
	if logger.Logger() == nil {
		logger.InitDefault("info")
	}
}
