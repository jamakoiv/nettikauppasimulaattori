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

var settings = Settings{
    "nettikauppasimulaattori",
    "store_operational",
    "orders",
    "order_items",

    "Europe/Helsinki",
}
    

/*
    Boilerplate for registering the function for the Eventarc Pub/Sub framework.
*/
type MessagePublishedData struct {
	Message PubSubMessage
}
type PubSubMessage struct {
	Data []byte `json:"data"`
}

/*
    Choose the execution entry-point e.g. start from this function.
    Has to match the entry-point selected in google functions.
*/
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

    /*
        Check all customers and send orders if any.
    */
    orders_in_this_run := false
    for _, customer := range Customers {
        order, err := customer.Shop(Products)
        if err != nil { continue } // If order is empty.

        slog.Debug(fmt.Sprint(order))
        slog.Info(fmt.Sprintf("Sending order %d to BigQuery.", order.id))
        err = order.Send(settings, ctx, client)
        if err != nil { 
            slog.Error(fmt.Sprintf("Error in sending order: %v", err))
        }

        orders_in_this_run = true
    }
    if !orders_in_this_run {
        slog.Info("No orders placed this time.")
    }

    /*
        Check all workers and update orders they complete.
    */
    order_ids, err := GetOpenOrders(settings, ctx, client)
    if err != nil { fmt.Println(err) }
    fmt.Println(order_ids)
    err = UpdateOrder(order_ids[0], settings, ctx, client)
    if err != nil { fmt.Println(err) }

    return nil
}
