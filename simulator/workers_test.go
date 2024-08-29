package nettikauppasimulaattori

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

var (
	test_row    [6]string = [6]string{"123", "Joku", "Tyyppi", "0 12-16 * * 1-5", "4", "15"}
	test_worker           = Worker{
		123, "Joku", "Tyyppi",
		"0 12-16 * * 1-5",
		4, 15,
	}
)

func TestCSVRowToWorker(t *testing.T) {
	result, _ := CSVRowToWorker(test_row[:])

	if !reflect.DeepEqual(result, test_worker) {
		t.Fatalf(fmt.Sprintf("Wanted %v, got %v", test_worker, result))
	}
}

func TestCheckIfWorkingSuccess(t *testing.T) {
	time_location, _ := time.LoadLocation("UTC")
	test_time_success := time.Date(2024, time.August, 29, 14, 0, 0, 0, time_location)

	result := test_worker.CheckIfWorking(test_time_success)

	if result != true {
		t.Fatalf(fmt.Sprintf("Wanted result = true."))
	}
}

func TestCheckIfWorkingFail(t *testing.T) {
	time_location, _ := time.LoadLocation("UTC")
	test_time_fail := time.Date(2024, time.August, 29, 21, 0, 0, 0, time_location)
	result := test_worker.CheckIfWorking(test_time_fail)

	if result != false {
		t.Fatalf(fmt.Sprintf("Wanted result = false."))
	}
}
