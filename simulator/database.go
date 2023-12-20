package nettikauppasimulaattori

import (
    "context"
    "cloud.google.com/go/bigquery"
    "fmt"
    "strings"

    "golang.org/x/exp/slog"
)

type Database interface {
    SendOrder(Order)
    GetOpenOrders() Orders
    UpdateOrder(Order)
    Close()
}


type DatabaseBigQuery struct {
    project string
    dataset string
    orders_table string
    order_items_table string

    timezone string

    ctx context.Context
    client *bigquery.Client
}

func (db *DatabaseBigQuery) Init(ctx context.Context, project string, dataset string, 
    orders_table string, order_items_table string, timezone string) error {

    db.project = project
    db.dataset = dataset
    db.orders_table = orders_table
    db.order_items_table = order_items_table
    db.timezone = timezone

    client, err := bigquery.NewClient(ctx, db.project)
    if err != nil { 
        slog.Error("Error creating BigQuery-client.") 
        return err
    }

    db.client = client
    return nil
}

func (db *DatabaseBigQuery) Close() {
    db.Close()
}

func (db *DatabaseBigQuery) SendOrder(order Order) error {
    order.order_placed, _ = nowInTimezone(db.timezone)

    order_sql := db.GetInsertOrderSQLquery(order)
    items_sql := db.GetInsertOrderItemsSQLquery(order)

    // slog.Debug(order_sql)
    // slog.Debug(items_sql)

    queries := [2]string{order_sql, items_sql}
    for _, sql := range queries {
        q := db.client.Query(sql)
        // q.WriteDisposition = "WRITE_APPEND" // Error with "INSERT INTO..." statement.

        job, err := q.Run(db.ctx)
        if err != nil {
            return err
        }

        status, err := job.Wait(db.ctx)
        if err != nil {
            return err
        }
        if status.Err() != nil {
            return status.Err()
        }
    }

    return nil
}

// Create SQL-query for inserting order to database.
func (db *DatabaseBigQuery) GetInsertOrderSQLquery(order Order) string {
    // TODO: guard against malicious inputs.
    order_sql := fmt.Sprintf("INSERT INTO `%s.%s.%s` VALUES (%d, %d, %d, %d, \"%s\", NULL, NULL, \"%s\")", 
        db.project,
        db.dataset,
        db.orders_table,
        order.id,
        order.customer_id,
        order.delivery_type,
        order.status,
        Time2SQLDatetime(order.order_placed),
        Time2SQLDatetime(order.order_placed))

    return order_sql
}


// Create SQL-query for inserting order items to database.
func (db *DatabaseBigQuery) GetInsertOrderItemsSQLquery(order Order) string {
    var tmp strings.Builder

    tmp.WriteString(fmt.Sprintf("INSERT INTO `%s.%s.%s` VALUES ",
        db.project, db.dataset, db.order_items_table))

    for _, item := range order.items {
        tmp.WriteString(fmt.Sprintf("(%d, %d),", order.id, item.id))
    }

    items_sql := tmp.String()
    items_sql = strings.TrimSuffix(items_sql, ",")

    return items_sql
}
