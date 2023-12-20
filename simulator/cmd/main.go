package main

import (
	"flag"
	"fmt"
	"os"

	"nettikauppasimulaattori.piste"

	"golang.org/x/exp/slog"
)

func main() {

	targetPtr := flag.String("target", "", "Run target. Valid inputs: 'prod' or 'test'.")
	verbosePtr := flag.Bool("verbose", false, "Output more debug-prints.")
	flag.Parse()

	// Setup logging.
	var logLevel = new(slog.LevelVar)
	if *verbosePtr == true {
		logLevel.Set(slog.LevelDebug)
	} else {
		logLevel.Set(slog.LevelInfo)
	}
	logger := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: logLevel})
	slog.SetDefault(slog.New(logger))

	// Run selected target.
	if *targetPtr == "prod" {
		nettikauppasimulaattori.Run_prod()
	} else if *targetPtr == "test" {
		nettikauppasimulaattori.Run_test()
	} else if *targetPtr == "" {
		slog.Error(fmt.Sprint("Empty run-target. Aborting..."))
	} else {
		slog.Error(fmt.Sprintf("Run target '%s' not valid.", *targetPtr))
	}
}
