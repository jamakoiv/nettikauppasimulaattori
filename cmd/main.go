package main

import (
	"nettikauppasimulaattori"
	"os"

	"golang.org/x/exp/slog"
)


func main() {

	var logLevel = new(slog.LevelVar)
	logger := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: logLevel})
	slog.SetDefault(slog.New(logger))
	logLevel.Set(slog.LevelDebug)

	nettikauppasimulaattori.Run()
}
