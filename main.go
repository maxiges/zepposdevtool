package main

import (
	"my-chart-app/internal/server"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// generate random data for line chart

func main() {
	// Config
	logger, _ := zap.NewDevelopment()
	cfg := zap.NewDevelopmentConfig()
	cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	defer logger.Sync()

	//run server

	server.RunServer()

}
