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
	functions.CloudEvent("Run", Run_gcloud_functions)
}

func RunWorkers(ctx context.Context, client* bigquery.Client) {
    for _, worker := range Workers {
        slog.Info(fmt.Sprintf("Checking worker %d.", worker.id))
        err := worker.Work(ctx, client)
        if err != nil {
            slog.Error(fmt.Sprint(err))
        }
    }
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

// Entry point for running via google cloud functions.
func Run_gcloud_functions(ctx context.Context, ev event.Event) error {
    // TODO: Project_id hardcoded in two different places :(
    project_id := "nettikauppasimulaattori"

    slog.Info(fmt.Sprintf("Program started at %v", time.Now()))

    client, err := bigquery.NewClient(ctx, project_id)
    if err != nil { 
        slog.Error("Error creating BigQuery-client.") 
        return err
    }
    defer client.Close()

    RunCustomers(ctx, client)
    RunWorkers(ctx, client)

    return nil
}


func Run_prod() error {

    return nil
}

func Run_test() error {

    return nil
}
