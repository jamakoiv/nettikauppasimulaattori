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
