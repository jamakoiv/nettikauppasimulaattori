package main

import (
	"context"
	"os"

	"github.com/cloudevents/sdk-go/v2/event"
	"nettikauppasimulaattori.piste"

	"golang.org/x/exp/slog"
)


func main() {

	var logLevel = new(slog.LevelVar)
	logger := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: logLevel})
	slog.SetDefault(slog.New(logger))
	logLevel.Set(slog.LevelDebug)

	ctx := context.Background()
	ev := event.Event{}

	nettikauppasimulaattori.Run(ctx, ev)
}
