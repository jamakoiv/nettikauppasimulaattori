package nettikauppasimulaattori

import (
    "context"
    "cloud.google.com/go/bigquery"
    "fmt"
    "strings"
    "math/rand"
    "time"

    "golang.org/x/exp/slog"
    "google.golang.org/api/iterator"
)

type Database interface {
    SendOrder(Order) error
    GetOpenOrders() (Orders, error)
    UpdateOrder(Order) error
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
    slog.Info(fmt.Sprint("Opening connection to BigQuery."))

    db.project = project
    db.dataset = dataset
    db.orders_table = orders_table
    db.order_items_table = order_items_table
    db.timezone = timezone
    db.ctx = ctx

    client, err := bigquery.NewClient(db.ctx, db.project)
    if err != nil { 
        slog.Error("Error creating BigQuery-client.") 
        return err
    }

    db.client = client
    return nil
}

func (db *DatabaseBigQuery) Close() {
    db.client.Close()
}

func (db *DatabaseBigQuery) SendOrder(order Order) error {
    slog.Info(fmt.Sprintf("Sending order %d to BigQuery.", order.id))

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
            slog.Error(fmt.Sprint("Error submitting query."))
            return err
        }

        status, err := job.Wait(db.ctx)
        if err != nil {
            slog.Error(fmt.Sprint("Error waiting for query to finish."))
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

func (db *DatabaseBigQuery) GetOpenOrdersSQLquery() string {
    return fmt.Sprintf("SELECT id, customer_id, delivery_type, status FROM `%s.%s.%s` WHERE status = %d", 
        db.project, db.dataset, db.orders_table, ORDER_PENDING)
}

func (db *DatabaseBigQuery) GetOpenOrders() (Orders, error) {
    var orders Orders

    sql := db.GetOpenOrdersSQLquery()

    q := db.client.Query(sql)
    job, err := q.Run(db.ctx)
    if err != nil { 
        slog.Error(fmt.Sprint("Error running query in GetOpenOrder: ", err))
    return orders, err }

    status, err := job.Wait(db.ctx)
    if err != nil { 
        slog.Error(fmt.Sprint("Error waiting in GetOpenOrder: ", err))
    return orders, err 
    }

    if status.Err() != nil { 
        slog.Error(fmt.Sprint("Error returned by bigquery: ", err))
    return orders, status.Err()
    }

    it, err := job.Read(db.ctx)
    if err != nil { 
        slog.Error(fmt.Sprint("Error reading data returned by bigquery: ", err))
    return orders, err
    }

    for {
        var order OrderReceiver
        if it.Next(&order) == iterator.Done { break }
        orders.Append(ConvertOrderReceiverToOrder(order))
    }

    if len(orders) == 0 {
        return orders, ErrorEmptyOrdersList
    } else {
        return orders, nil
    }
}

func (db *DatabaseBigQuery) GetUpdateOrderSQLquery(order Order, t time.Time) string {
    sql := fmt.Sprintf("UPDATE `%s.%s.%s` SET status = %d, shipping_date = \"%s\", last_modified = \"%s\", tracking_number = %d WHERE id = %d",
        db.project, db.dataset, db.orders_table, 
        ORDER_SHIPPED,
        Time2SQLDate(t), Time2SQLDatetime(t), 
        rand.Int(),
        order.id)

    return sql
}

func (db *DatabaseBigQuery) UpdateOrder(order Order) error {
    now, _ := nowInTimezone(db.timezone)

    sql := db.GetUpdateOrderSQLquery(order, now)

    q := db.client.Query(sql)
    job, err := q.Run(db.ctx)
    if err != nil { 
        slog.Error(fmt.Sprint("Error running query in UpdateOrder: ", err))
        return err 
    }
    status, err := job.Wait(db.ctx)
    if err != nil { 
        slog.Error(fmt.Sprint("Error waiting query in UpdateOrder: ", err))
        return err 
    }
    if status.Err() != nil { 
        slog.Error(fmt.Sprint("Received error from bigquery: ", err))
        return status.Err() 
    }

    return nil
}


type DatabaseBigQueryDummy struct {
    db DatabaseBigQuery
}

func (dummy *DatabaseBigQueryDummy) Init(ctx context.Context, project string, dataset string, 
    orders_table string, order_items_table string, timezone string) error {
    slog.Info(fmt.Sprint("Opening DummyDatabase-connection."))

    dummy.db.project = project
    dummy.db.dataset = dataset
    dummy.db.orders_table = orders_table
    dummy.db.order_items_table = order_items_table
    dummy.db.timezone = timezone
    dummy.db.ctx = ctx

    return nil
}

func (dummy *DatabaseBigQueryDummy) Close() {
    slog.Info(fmt.Sprint("Closing DummyDatabase-connection."))
    return
}

func (dummy *DatabaseBigQueryDummy) SendOrder(order Order) error {
    slog.Info(fmt.Sprintf("Sending order %d to DummyDatabase.", order.id))
    slog.Debug(fmt.Sprintf("Order query: %s", dummy.db.GetInsertOrderSQLquery(order)))
    slog.Debug(fmt.Sprintf("Order-items query: %s", dummy.db.GetInsertOrderItemsSQLquery(order)))

    return nil
}

func (dummy *DatabaseBigQueryDummy) GetOpenOrders() (Orders, error) {
    slog.Info(fmt.Sprintf("Getting open orders from DummyDatabase."))
    slog.Debug(fmt.Sprintf("Query: %s", dummy.db.GetOpenOrdersSQLquery()))
    
    return Orders{}, nil
}

func (dummy *DatabaseBigQueryDummy) UpdateOrder(order Order) error {
    now, _ := nowInTimezone(dummy.db.timezone)

    slog.Info(fmt.Sprintf("Updating order %d in DummyDatabase.", order.id))
    slog.Debug(fmt.Sprintf("Query: %s", dummy.db.GetUpdateOrderSQLquery(order, now)))

    return nil
}
