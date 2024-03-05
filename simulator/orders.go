package nettikauppasimulaattori

import (
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"time"
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
// NOTE: All variables have to be public or big-query library fails silently
// when receiving data.
type OrderReceiver struct {
	ID            int
	Customer_id   int
	Delivery_type int
	Status        int
	// Order_placed  time.Time
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

var ErrorEmptyOrdersList = errors.New("list is empty")

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

func ConvertOrderReceiverToOrder(o OrderReceiver) Order {
	var res Order

	res.id = uint64(o.ID)
	res.customer_id = o.Customer_id
	// res.order_placed = o.Order_placed
	res.delivery_type = o.Delivery_type
	res.status = o.Status

	return res
}
