package nettikauppasimulaattori

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)


func TestCSVRemoveWhitespace(t *testing.T) {
    test_str := [2]string{"\t 1234", " abc "}
    target := [2]string{"1234", "abc"}

    result := CSVRemoveWhitespace(test_str[:])

    if reflect.DeepEqual(result, target) {
        t.Fatalf(fmt.Sprintf("Received %v, wanted %v.", result, target)) 
    }
}

func TestCSVSplitSerial2Int(t *testing.T) {
    test_str := "{1;2;3;4}"
    target := [4]int{1,2,3,4} // NOTE: Has to be array, not a slice, or the comparison fails.

    result, _ := CSVSplitSerial2Int(test_str, "{}", ";")

    if reflect.DeepEqual(result, target) {
        t.Fatalf(fmt.Sprintf("Received %v, wanted %v.", result, target)) 
    }
}

func TestIntToWeekday(t *testing.T) {
    test_input := [4]int{1,2,3,4}
    target := [4]time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday}

    result, _ := IntToWeekday(test_input[:])

    if reflect.DeepEqual(result, target) {
        t.Fatalf(fmt.Sprintf("Received %v, wanted %v", result, target))
    }
}
