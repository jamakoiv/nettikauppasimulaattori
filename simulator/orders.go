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
    order_placed time.Time
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

func nowInTimezone(timezone string) (time.Time, error) {
    var t time.Time

    tz, err := time.LoadLocation(timezone)

    if err != nil {
        err_str := fmt.Sprintf("Error getting timezone 'time.LoadLocation(%s'): %s", timezone, err)
        slog.Error(err_str)
        t = time.Now()
    } else {
        t = time.Now().In(tz)
    }

    return t, err
}

func Time2SQLDatetime(t time.Time) string {
    res := fmt.Sprintf("%d-%d-%d %d:%d:%d",
        t.Year(), t.Month(), t.Day(),
        t.Hour(), t.Minute(), t.Second())

    return res
}

func Time2SQLDate(t time.Time) string {
    res := fmt.Sprintf("%d-%d-%d", 
        t.Year(), t.Month(), t.Day())

    return res
}

func (order *Order) init() {
    order.id = uint64(rand.Uint32())  // Foolishly hope we don't get two same order IDs.
    order.status = ORDER_EMPTY
    order.delivery_type = rand.Intn(2)
    order.order_placed = time.Now()
}

func (order *Order) AddItem(item Product) {
    order.items = append(order.items, item) 
}

func (order *Order) TotalPrice() int {
    var total int
    for _, item := range order.items {
        total += item.price
    }
    return total
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


func GetInsertOrderSQLquery(order *Order) string {
    // Create SQL-query for inserting order to database.
    project_id := "nettikauppasimulaattori"
    dataset_id := "store_operational"
    orders_table_id := "orders"

    // TODO: guard against malicious inputs.
    order_sql := fmt.Sprintf("INSERT INTO `%s.%s.%s` VALUES (%d, %d, %d, %d, \"%s\", NULL, NULL)",     project_id, 
        dataset_id, 
        orders_table_id,
        order.id, 
        order.customer_id, 
        order.delivery_type, 
        order.status, 
        Time2SQLDatetime(order.order_placed))

    return order_sql
}

func GetInsertOrderItemsSQLquery(order *Order) string {
    // Create SQL-query for inserting order items to database.
    project_id := "nettikauppasimulaattori"
    dataset_id := "store_operational"
    order_items_table_id := "order_items"

    var tmp strings.Builder

    tmp.WriteString(fmt.Sprintf("INSERT INTO `%s.%s.%s` VALUES ",
                                project_id,
                                dataset_id, 
                                order_items_table_id))

    for _, item := range order.items {
        tmp.WriteString(fmt.Sprintf("(%d, %d),", 
                                    order.id,
                                    item.id))
    }

    items_sql := tmp.String()
    items_sql = strings.TrimSuffix(items_sql, ",")

    return items_sql
}


func (order *Order) Send(ctx context.Context, client *bigquery.Client) error {

    timezone := "Europe/Helsinki"
    order.order_placed, _ = nowInTimezone(timezone)

    order_sql := GetInsertOrderSQLquery(order)
    items_sql := GetInsertOrderItemsSQLquery(order)

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
