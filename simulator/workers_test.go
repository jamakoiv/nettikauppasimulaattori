package nettikauppasimulaattori

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)


func TestCSVRowToWorker(t *testing.T) {
    test_row := [7]string{"123", "Joku", "Tyyppi", "{1;2;3}", 
        "{10;11;12;13;14}", "4", "15"}

    target := Worker{123, "Joku", "Tyyppi",
        []time.Weekday{time.Monday, time.Tuesday, time.Wednesday},
        []int{10, 11, 12, 13, 14},
        4, 15,}

    result, _ := CSVRowToWorker(test_row[:])

    if !reflect.DeepEqual(result, target) {
        t.Fatalf(fmt.Sprintf("Wanted %v, got %v", target, result))
    }
}
