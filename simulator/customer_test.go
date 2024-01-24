package nettikauppasimulaattori

import (
	"fmt"
	"math"
	"reflect"
	"testing"
	"time"
)


func TestDefault_ShoppingWeekVariation(t *testing.T) {

    test_date := time.Date(2024, time.January, 30, 0, 0, 0, 0, time.UTC)
    target := 0.004
    allowed_error := 0.00001

    res := Default_ShoppingWeekVariation(test_date)
    
    if (math.Abs(res - target) > allowed_error) {
        t.Fatalf("Wanted %v, got %v, which is not within allowed error of %v",
            res, target, allowed_error)
    }
}


func TestCSVRowToCustomer(t *testing.T) {

    test_row := [7]string{"1", "Jaska", "Jokunen", "10", "300", "0.25", "{1;2;3}"}
    target := Customer{1, "Jaska", "Jokunen", 10, 300, 0.25, []int{1,2,3}}
    
    res, _ := CSVRowToCustomer(test_row[:])

    if !reflect.DeepEqual(res, target) {
        t.Fatalf("Input %v created object %v, which is incorrect.", test_row, res)
    }
}

func TestReadCustomerCSV(t * testing.T) {

    test_file := "test_data/customers_test.csv"
    target := []Customer{
        {10,  "Erkki",    "Nyrhinen"   , 6,  500, 0.15, []int{1,2} },
        {11,  "Jaana",    "Lahtinen"   , 7,  250, 0.25, []int{2} },   
        {12,  "Toni",     "Kuusisto"   , 8,  50 , 0.10, []int{1} },
    }

    res, err := ReadCustomersCSV(test_file)

    if err != nil {
        t.Fatalf(fmt.Sprint(err))
    }

    if !reflect.DeepEqual(res, target) {
        t.Fatalf("Reading from file %v created object %v, which is incorrect.", test_file, res)
    }
}
