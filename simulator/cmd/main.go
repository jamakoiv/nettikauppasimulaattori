package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/cloudevents/sdk-go/v2/event"
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
		logLevel.Set(slog.LevelError)
	}
	logger := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: logLevel})
	slog.SetDefault(slog.New(logger))

	// Run selected target.
	if *targetPtr == "prod" {
		ctx := context.Background()
		ev := event.Event{}

		nettikauppasimulaattori.Run(ctx, ev)
	} else if *targetPtr == "test" {
		// Implement test-run.
		slog.Error(fmt.Sprint("Run target 'test' not implemented yet."))
	} else if *targetPtr == "" {
		slog.Error(fmt.Sprint("Empty run-target. Aborting..."))
	} else {
		slog.Error(fmt.Sprintf("Run target '%s' not valid.", *targetPtr))
	}
}
