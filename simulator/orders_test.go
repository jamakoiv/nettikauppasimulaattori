package nettikauppasimulaattori

import (
	"fmt"
	"reflect"
	"testing"
	"time"
    "errors"
)

var OrdersTestProducts = []Product{   
    {1001, "Pirkka olut 24-pak.", 10, 25, 0.24},
    {2001, "Raspberry Pi 4 4GB", 40, 80, 0.24},
    {3001, "Ruisleipä", 1, 3, 0.14},
    {4000, "Silmarillion, J.R.R Tolkien", 10, 25, 0.10},
}   

var test_order = Order{
    1234,
    9876,
    []Product{OrdersTestProducts[0], OrdersTestProducts[1]},
    time.Date(1234, time.Month(5), 6, 7, 8, 9, 0, time.UTC),
    0,
    0 }


func TestTime2SQLDatetime(t *testing.T) {
    target := "1234-5-6 7:8:9"
    res := Time2SQLDatetime(test_order.order_placed)

    if res != target {
        t.Fatalf("Wanted %v, got %v", target, res)
    }
}

func TestTime2SQLDate(t *testing.T) {
    target := "1234-5-6"
    res := Time2SQLDate(test_order.order_placed)

    if res != target {
        t.Fatalf("Wanted %v, got %v", target, res)
    }
}

// func TestGetInsertOrderSQLquery(t *testing.T) {
//     // TODO: Hardcoded table names...
//     target := "INSERT INTO `nettikauppasimulaattori.store_operational.orders` VALUES (1234, 9876, 0, 0, \"1234-5-6 7:8:9\", NULL, NULL, \"1234-5-6 7:8:9\")"
//     res := GetInsertOrderSQLquery(&test_order)
// 
//     if res != target {
//         t.Fatalf("Wanted '%v', got '%v' instead.", target, res)
//     }
// }
// 
// func TestGetInsertOrderItemsSQLquery(t *testing.T) {
//     // TODO: Hardcoded table names...
//     target := "INSERT INTO `nettikauppasimulaattori.store_operational.order_items` VALUES (1234, 1001),(1234, 2001)"
//     res := GetInsertOrderItemsSQLquery(&test_order)
// 
//     if res != target {
//         t.Fatalf("Wanted '%v', got '%v' instead.", target, res)
//     }
// }

func TestTotalPrice(t *testing.T) {
    target := 105
    res := test_order.TotalPrice()

    if res != target {
        t.Fatalf("Wanted %v, got %v instead.", target, res)
    }
}

func TestOrdersAppend(t *testing.T) {
    test_time := time.Date(1, 2, 3, 4, 5, 6, 7, time.UTC)

    order := Order{1111, 2222, []Product{}, test_time, 0, 0}
    target := Orders{order}

    res := Orders{}
    res.Append(order)

    if !reflect.DeepEqual(res, target) {
        t.Fatalf("Append failed: wanted %s, got %s",
            fmt.Sprint(target), fmt.Sprint(res))
    }
}


func TestOrdersPop(t *testing.T) {
    test_time := time.Date(1, 2, 3, 4, 5, 6, 7, time.UTC)

    order_A := Order{1111, 2222, []Product{}, test_time, 0, 0}
    order_B := Order{3333, 4444, []Product{}, test_time, 0, 0}
    orders := Orders{order_A, order_B} 

    target_order := order_A
    target_orders := Orders{order_B} 

    res, _ := orders.Pop()
    if !reflect.DeepEqual(res, target_order) {
        t.Fatalf("Pop failed to return correct order: wanted %s, got %s.", 
            fmt.Sprint(target_order), fmt.Sprint(res))
    }

    if !reflect.DeepEqual(orders, target_orders) {
        t.Fatalf("List of orders incorrect after using Pop.")
    }
}

func TestOrdersPopError(t *testing.T) {
    orders := Orders{}

    _, err := orders.Pop()
    if !errors.Is(err, ErrorEmptyOrdersList) {
        t.Fatalf("Wrong error received when using Pop() on empty Orders-list.")
    }
}


func TestConvertOrderReceiverToOrder(t *testing.T) {
    receiver := OrderReceiver{1234, 555, 5, 5}

    target := Order{1234, 555, []Product{}, time.Time{}, 5, 5}
    res := ConvertOrderReceiverToOrder(receiver)

    // For some reason reflect.DeepEqual doesn't work properly here.
    if res.id != target.id ||
        res.customer_id != target.customer_id ||
        res.status != target.status ||
        res.delivery_type != target.delivery_type {

        t.Fatalf("Failed to convert OrderReceiver to Order: wanted %s, got %s.",
            fmt.Sprint(target), fmt.Sprint(res))
    }
}
