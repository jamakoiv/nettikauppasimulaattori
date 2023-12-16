package nettikauppasimulaattori

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"strings"
	"time"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

// Type for creating an order, adding products to it,
// and sendind in to database.
type Order struct {
	id            uint64
	customer_id   int
	items         []Product
	order_placed  time.Time
	delivery_type int
	status        int
}
type Orders []Order

// Type for temporarily storing order-data returned by the database.
type OrderReceiver struct {
	id            uint64
	customer_id   int
	delivery_type int
	status        int
	order_placed  time.Time
}

const ( // Values for Order.status.
	ORDER_PENDING = iota
	ORDER_SHIPPED = iota
	ORDER_EMPTY   = iota
)

const ( // Values for Order.delivery_type.
	SHIP_TO_CUSTOMER   = iota
	COLLECT_FROM_STORE = iota
)

var ErrorEmptyOrdersList = errors.New("List is empty.")


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
	order.id = uint64(rand.Uint32()) // Foolishly hope we don't get two same order IDs.
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
	if order.status == ORDER_EMPTY {
		return ""
	}

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
	order_sql := fmt.Sprintf("INSERT INTO `%s.%s.%s` VALUES (%d, %d, %d, %d, \"%s\", NULL, NULL, \"%s\")", project_id,
		dataset_id,
		orders_table_id,
		order.id,
		order.customer_id,
		order.delivery_type,
		order.status,
		Time2SQLDatetime(order.order_placed),
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
		if err != nil {
			return err
		}

		status, err := job.Wait(ctx)
		if err != nil {
			return err
		}
		if status.Err() != nil {
			return status.Err()
		}
	}

	return nil
}


func (orders *Orders) Append(order Order) {
    *orders = append(*orders, order)
}


func (orders *Orders) Pop() (Order, error) {
	if len(*orders) == 0 {
		return Order{}, ErrorEmptyOrdersList 
	}

	first := (*orders)[0]

	// [1:] panics if length is 1, se we create and return empty list
	// if we pop last element out.
	if len(*orders) == 1 {
		orders = new(Orders) 
	} else {
		(*orders) = (*orders)[1:]
	}

	return first, nil
}


func GetOpenOrders(ctx context.Context, client *bigquery.Client) (Orders, error) {
    // TODO: Move ids to config file somewhere.
    project_id := "nettikauppasimulaattori"
    dataset_id := "store_operational"
    table_id := "orders"

    sql := fmt.Sprintf("SELECT id, customer_id, delivery_type, status, order_placed FROM `%s.%s.%s` WHERE status = %d",
        project_id, dataset_id, table_id, ORDER_PENDING)
    slog.Debug(sql)

    var orders Orders

    // slog.Debug("Sending query.")
    q := client.Query(sql)
    job, err := q.Run(ctx)
    if err != nil { 
		slog.Error(fmt.Sprint("Error running query in GetOpenOrder: ", err))
        return orders, err }

    // slog.Debug("Wait query.")
    status, err := job.Wait(ctx)
    if err != nil { 
		slog.Error(fmt.Sprint("Error waiting in GetOpenOrder: ", err))
        return orders, err 
    }

    // slog.Debug("Check status.")
    if status.Err() != nil { 
		slog.Error(fmt.Sprint("Error returned by bigquery: ", err))
        return orders, status.Err()
    }

    // slog.Debug("Get iterator.")
    it, err := job.Read(ctx)
    if err != nil { 
		slog.Error(fmt.Sprint("Error reading data returned by bigquery: ", err))
        return orders, err
    }

    // slog.Debug("Parse results.")
    for {
        var order OrderReceiver
        if it.Next(&order) == iterator.Done { break }
        // fmt.Printf("%d: %T\n", tmp.ID, tmp.ID)
		orders.Append(ConvertOrderReceiverToOrder(order))
    }
	
	if len(orders) == 0 {
		return orders, ErrorEmptyOrdersList
	} else {
		return orders, nil
	}
}


func ConvertOrderReceiverToOrder(o OrderReceiver) Order {
    var res Order

    res.id = o.id
    res.customer_id = o.customer_id
    res.order_placed = o.order_placed
    res.delivery_type = o.delivery_type
    res.status = o.status
        
    return res
}


func UpdateOrder(order Order, ctx context.Context, client *bigquery.Client) error {
    project_id := "nettikauppasimulaattori"
    dataset_id := "store_operational"
    table_id := "orders"

    now, _ := nowInTimezone("Europe/Helsinki")

    sql := fmt.Sprintf("UPDATE `%s.%s.%s` SET status = %d, shipping_date = \"%s\", last_modified = \"%s\", tracking_number = %d WHERE id = %d",
        project_id, 
        dataset_id,
        table_id, 
        ORDER_SHIPPED,
        Time2SQLDate(now), 
        Time2SQLDatetime(now), 
        rand.Int(),
        order.id)
    slog.Debug(sql)

    q := client.Query(sql)
    job, err := q.Run(ctx)
    if err != nil { 
		slog.Error(fmt.Sprint("Error running query in UpdateOrder: ", err))
		return err 
	}

    status, err := job.Wait(ctx)
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
