package nettikauppasimulaattori

import (
    "fmt"
    "math/rand"
    "context"
    "time"
    "strings"
    "log/slog"

    "cloud.google.com/go/bigquery"
)


type Order struct {
    id uint64
    customer_id int
    items []Product
    delivery_type int
    status int
}

const (     // Values for Order.status.
    ORDER_PENDING = iota
    ORDER_SHIPPED = iota
    ORDER_EMPTY = iota
)

const (     // Values for Order.delivery_type.
    SHIP_TO_CUSTOMER = iota
    COLLECT_FROM_STORE = iota
)


func Now2SQLDatetime(timezone string) string {
    // Return current time as SQL Datetime.
    var t time.Time
    tz, err := time.LoadLocation(timezone)

    if err != nil {
        err_str := fmt.Sprintf("Error getting timezone 'time.LoadLocation(%s'): %s", timezone, err)
        slog.Error(err_str)
        t = time.Now()
    } else {
        t = time.Now().In(tz)
    }

    return fmt.Sprintf("%d-%d-%d %d:%d:%d",
        t.Year(), t.Month(), t.Day(),
        t.Hour(), t.Minute(), t.Second())
}


func (order *Order) init() {
    order.id = uint64(rand.Uint32())  // Foolishly hope we don't get two same order IDs.
    order.status = ORDER_EMPTY
    order.delivery_type = rand.Intn(2)
}

func (order *Order) AddItem(item Product) {
    order.items = append(order.items, item) 
}

// Satisfy Stringer-interface.
func (order *Order) String() string {
    if order.status == ORDER_EMPTY { return "" }

    var str string = fmt.Sprintf("Order %v\n--------------\n", order.id)

    for _, item := range order.items {
        str = str + fmt.Sprintf("%v: %v\n", order.customer_id, item.name)
    }
    return str
}

func (order *Order) Send(ctx context.Context, client *bigquery.Client) error {
    // TODO: Break creating the SQL-queries into separate functions.
    // TODO: Store project_id etc in separate config-file.

    project_id := "nettikauppasimulaattori"
    dataset_id := "store_operational"
    orders_table_id := "orders"
    order_items_table_id := "order_items"

    log_timezone := "Europe/Helsinki"

    // TODO: guard against malicious inputs.
    order_sql := fmt.Sprintf("INSERT INTO `%s.%s.%s` VALUES ", 
        project_id, dataset_id, orders_table_id)
    order_sql = fmt.Sprintf("%s (%d, %d, %d, %d, \"%s\", NULL, NULL)", 
        order_sql, order.id, order.customer_id, 
        order.delivery_type, order.status, Now2SQLDatetime(log_timezone))

    items_sql := fmt.Sprintf("INSERT INTO `%s.%s.%s` VALUES ", 
        project_id, dataset_id, order_items_table_id)

    for _, item := range order.items {
        items_sql = fmt.Sprintf("%s (%d, %d),", items_sql, order.id, item.id)
    }
    items_sql = strings.TrimSuffix(items_sql, ",")

    // slog.Debug(order_sql)
    // slog.Debug(items_sql)

    queries := [2]string{order_sql, items_sql}
    for _, sql := range queries {
        q := client.Query(sql)
        // q.WriteDisposition = "WRITE_APPEND" // Error with "INSERT INTO..." statement.

        job, err := q.Run(ctx)
        if err != nil { return err }

        status, err := job.Wait(ctx)
        if err != nil { return err }
        if status.Err() != nil { return status.Err() }
    }

    return nil
}
