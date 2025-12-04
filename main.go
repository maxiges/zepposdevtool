package main

import (
	"flag"
	"fmt"
	"zepp-os-dev-tool/internal/gui"
	"zepp-os-dev-tool/internal/server"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Package main launches the ZeppOsDevTool application.
//
// ZeppOsDevTool is a development helper for ZeppOS.
//
// Flags:
//
//	-height int      Window height (default 1800)
//	-width int       Window width (default 900)
//	-disableGUI bool Don't show the GUI (by default the GUI is shown)
func main() {
	// Config

	var wFlag = flag.Int("height", 0, fmt.Sprintf("Window height (default %d)", 1800))
	var hFlag = flag.Int("width", 0, fmt.Sprintf("Window width (default %d)", 900))
	var disableGUI = flag.Bool("disableGUI", false, "don't show GUI (By default, GUI is shown)")

	flag.Parse()

	logger, _ := zap.NewDevelopment()
	cfg := zap.NewDevelopmentConfig()
	cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	defer logger.Sync()

	//run server
	// If the `-disableGUI` flag is set, the GUI is disabled and only the
	// server runs. Otherwise, start the GUI in a separate goroutine and
	// then run the server in the main goroutine so both are active.
	if disableGUI != nil && *disableGUI {
		// Explicitly disable GUI (implementation in `internal/gui`).
		gui.DisableGui()
		// Run the HTTP server (blocks until stopped).
		server.RunServer()

	} else {
		go gui.RunGui(wFlag, hFlag)
		// Run the HTTP server (blocks until stopped).
		server.RunServer()
	}
}
