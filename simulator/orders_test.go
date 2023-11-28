package nettikauppasimulaattori

import (
    "testing"
    "time"
)

var OrdersTestProducts = []Product{   
    {1001, "Pirkka olut 24-pak.", 10, 25, 0.24},
    {2001, "Raspberry Pi 4 4GB", 40, 80, 0.24},
    {3001, "Ruisleipä", 1, 3, 0.14},
    {4000, "Silmarillion, J.R.R Tolkien", 10, 25, 0.10},
}   



func TestTime2SQLDatetime(t *testing.T) {
    var test_time time.Time
    test_time = time.Date(1234, time.Month(6), 7, 8, 9, 10, 0, test_time.Location())

    target := "1234-6-7 8:9:10"
    res := Time2SQLDatetime(test_time)

    if res != target {
        t.Fatalf("Wanted %v, got %v", target, res)
    }
}

func TestTime2SQLDate(t *testing.T) {
    var test_time time.Time
    test_time = time.Date(1234, time.Month(6), 7, 0, 0, 0, 0, test_time.Location())

    target := "1234-6-7"
    res := Time2SQLDate(test_time)

    if res != target {
        t.Fatalf("Wanted %v, got %v", target, res)
    }
}

func TestGetInsertOrderSQLquery(t *testing.T) {
    var test_order Order 
    test_order.id = 1234
    test_order.customer_id = 9876
    test_order.delivery_type = 0
    test_order.status = 0

    test_order.AddItem(OrdersTestProducts[0])
    test_order.AddItem(OrdersTestProducts[1])

    // TODO: Again have to get hack the way around time.Now in testing...
    target := "INSERT INTO `nettikauppasimulaattori.store_operational.orders` VALUES (1234, 9876, 0, 0, \"1234-5-6 7:8:9\", NULL, NULL)"
    res := GetInsertOrderSQLquery(&test_order)

    if res != target {
        t.Fatalf("Wanted '%v', got '%v' instead.", target, res)
    }
}
