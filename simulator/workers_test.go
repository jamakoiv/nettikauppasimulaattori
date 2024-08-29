package nettikauppasimulaattori

import (
	"fmt"
	"reflect"
	"testing"
)

func TestCSVRowToWorker(t *testing.T) {
	test_row := [6]string{"123", "Joku", "Tyyppi", "0 12-16 * * 1-5", "4", "15"}

	target := Worker{
		123, "Joku", "Tyyppi",
		"0 12-16 * * 1-5",
		4, 15,
	}

	result, _ := CSVRowToWorker(test_row[:])

	if !reflect.DeepEqual(result, target) {
		t.Fatalf(fmt.Sprintf("Wanted %v, got %v", target, result))
	}
}
