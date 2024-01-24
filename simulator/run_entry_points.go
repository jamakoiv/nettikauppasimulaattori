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
        slog.Debug(fmt.Sprintf("Running worker %d.", worker.id))
        err := worker.Work(db)
        if err != nil {
            slog.Error(fmt.Sprint(err))
        }
    }
}

func RunCustomers(db Database, customers []Customer) {
    orders_in_this_run := false
    for _, customer := range customers {
        slog.Debug(fmt.Sprintf("Running customer %d.", customer.id))
        order, err := customer.Shop(Products)
        if err != nil { continue } // If order is empty.

        slog.Debug(fmt.Sprint(order))
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

    customers, err := ReadCustomersCSV("data/customers.csv")
    if err != nil { slog.Error(fmt.Sprintf("Failed to read customers data from file: %v", err)) }

    RunCustomers(&db, customers)
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
    defer db.Close()

    customers, err := ReadCustomersCSV("data/customers.csv")
    if err != nil { slog.Error(fmt.Sprintf("Failed to read customers data from file: %v", err)) }

    RunCustomers(&db, customers)
    RunWorkers(&db)

    return nil
}

// Entry point for running test-run locally.
func Run_test() error {
    slog.Info(fmt.Sprintf("Dry run started locally at %v", time.Now()))

    var db DatabaseBigQueryDummy
    err := db.Init(context.Background(), 
            "nettikauppasimulaattori",
            "store_operational",
            "orders",
            "order_items",
            "Europe/Helsinki")
    if err != nil { slog.Error("Database init failed.") }
    defer db.Close()

    customers, err := ReadCustomersCSV("data/customers.csv")
    if err != nil { slog.Error(fmt.Sprintf("Failed to read customers data from file: %v", err)) }

    RunCustomers(&db, customers)
    RunWorkers(&db)

    return nil
}
