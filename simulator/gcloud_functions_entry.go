package nettikauppasimulaattori

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
	"golang.org/x/exp/slog"
)

type MessagePublishedData struct {
	Message PubSubMessage
}
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// Define google cloud functions entry-point.
func init() {
	functions.CloudEvent("Run", Run)
}

func RunWorkers(ctx context.Context, client* bigquery.Client) {
    for _, worker := range Workers {
        slog.Info(fmt.Sprintf("Checking worker %d.", worker.id))
        err := worker.Work(ctx, client)
        if err != nil {
            slog.Error(fmt.Sprint(err))
        }
    }

    // DebugWorker := Worker{777, "Debug", "Debugger",
    //     []time.Weekday{time.Monday, 
    //         time.Tuesday, 
    //         time.Wednesday, 
    //         time.Thursday,
    //         time.Friday},
    //     []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23},
    //     3, 15}

    // err := DebugWorker.Work(ctx, client)
    // if err != nil {
    //     slog.Error(fmt.Sprint(err))
    // }
}


func RunCustomers(ctx context.Context, client* bigquery.Client) {
    orders_in_this_run := false
    for _, customer := range Customers {
        order, err := customer.Shop(Products)
        if err != nil { continue } // If order is empty.

        slog.Debug(fmt.Sprint(order))
        slog.Info(fmt.Sprintf("Sending order %d to BigQuery.", order.id))
        err = order.Send(ctx, client)
        if err != nil { 
            slog.Error(fmt.Sprintf("Error in sending order: %v", err))
        }

        orders_in_this_run = true
    }
    if !orders_in_this_run {
        slog.Info("No orders placed this time.")
    }
}

func Run(ctx context.Context, ev event.Event) error {
    // TODO: Project_id hardcoded in two different places :(
    project_id := "nettikauppasimulaattori"

    slog.Info(fmt.Sprintf("Program started at %v", time.Now()))

    client, err := bigquery.NewClient(ctx, project_id)
    if err != nil { 
        slog.Error("Error creating BigQuery-client.") 
        return err
    }
    defer client.Close()

    // RunCustomers(ctx, client)
    RunWorkers(ctx, client)

    return nil
}
