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

    order_ids, err := GetOpenOrders(ctx, client)
    if err != nil { fmt.Println(err) }
    fmt.Println(order_ids)
    err = UpdateOrder(order_ids[0], ctx, client)
    if err != nil { fmt.Println(err) }

    return nil
}
