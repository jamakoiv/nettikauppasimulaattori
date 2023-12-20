package nettikauppasimulaattori

import (
	"context"
	"fmt"
	"time"

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

func RunWorkers(db Database) {
    for _, worker := range Workers {
        slog.Info(fmt.Sprintf("Checking worker %d.", worker.id))
        err := worker.Work(db)
        if err != nil {
            slog.Error(fmt.Sprint(err))
        }
    }
}

func RunCustomers(db Database) {
    orders_in_this_run := false
    for _, customer := range Customers {
        order, err := customer.Shop(Products)
        if err != nil { continue } // If order is empty.

        slog.Debug(fmt.Sprint(order))
        slog.Info(fmt.Sprintf("Sending order %d to BigQuery.", order.id))
        err = db.SendOrder(order)
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
    slog.Info(fmt.Sprintf("Program started with Run_gcloud_functions at %v", time.Now()))

    var db DatabaseBigQuery
    err := db.Init(ctx, 
            "nettikauppasimulaattori",
            "store_operational",
            "orders",
            "order_items",
            "Europe/Helsinki")
    if err != nil { slog.Error("Database init failed.") }
    defer db.Close()

    RunCustomers(&db)
    RunWorkers(&db)

    return nil
}

// Entry point for running locally.
func Run_prod() error {
    slog.Info(fmt.Sprintf("Program started locally at %v", time.Now()))

    var db DatabaseBigQuery
    err := db.Init(context.Background(), 
            "nettikauppasimulaattori",
            "store_operational",
            "orders",
            "order_items",
            "Europe/Helsinki")
    if err != nil { slog.Error("Database init failed.") }
    // defer db.Close()

    RunCustomers(&db)
    RunWorkers(&db)

    return nil
}

// Entry point for running test-run locally.
func Run_test() error {
    slog.Info(fmt.Sprintf("Dry run started locally at %v", time.Now()))

    slog.Error(fmt.Sprint("Run target 'test' not implemented yet."))

    return nil
}
