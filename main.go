package main

import (
	"zepp-os-dev-tool/internal/server"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Package main launches the ZeppOsDevTool application.
//
// ZeppOsDevTool is a development helper for ZeppOS.
func main() {
	// Config

	logger, _ := zap.NewDevelopment()
	cfg := zap.NewDevelopmentConfig()
	cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	defer logger.Sync()

	server.RunServer()
}
