package main

import (
	"flag"
	"fmt"
	"my-chart-app/internal/gui"
	"my-chart-app/internal/server"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// generate random data for line chart

func main() {
	// Config

	var wFlag = flag.Int("height", 0, fmt.Sprintf("Window height (default %d)", 1800))
	var hFlag = flag.Int("width", 0, fmt.Sprintf("Window width (default %d)", 900))
	var disableGUI = flag.Bool("disableGUI", false, fmt.Sprintf("don't show GUI (By default, GUI is shown)"))

	flag.Parse()

	logger, _ := zap.NewDevelopment()
	cfg := zap.NewDevelopmentConfig()
	cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	defer logger.Sync()

	//run server

	if disableGUI != nil && *disableGUI {
		gui.DisableGui()
		server.RunServer()

	} else {
		go server.RunServer()
		gui.RunGui(wFlag, hFlag)

	}

}
