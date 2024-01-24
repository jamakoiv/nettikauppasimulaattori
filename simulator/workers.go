package nettikauppasimulaattori

import (
	"errors"
	"fmt"
	"slices"
	"time"

	"golang.org/x/exp/slog"
        "os"
        "encoding/csv"
        "strconv"
)

type Worker struct {
    id int
    first_name string
    last_name string
    work_days []time.Weekday
    work_hours []int
    orders_per_hour int
    salary_per_hour int
}

var VALID_CSV_WORKER_SIZE int = 7

type WorkerCsvError struct {
    length int
}
func (e *WorkerCsvError) Error() string {
    return fmt.Sprintf("Received CSV-row with %v elements. Should be %v elements.",
        e.length, VALID_CSV_WORKER_SIZE)
}

// TODO: Duplication of code in different CSV-readers.
func ReadWorkersCSV(file string) ([]Worker, error) {
    var res []Worker
    
    f, err := os.Open(file)
    if err != nil { return res, err }

    reader := csv.NewReader(f)
    rows, err := reader.ReadAll()
    if err != nil { return res, err }

    for _, row := range rows {
        row = CSVRemoveWhitespace(row)
        worker, err := CSVRowToWorker(row)
        if err != nil { 
            slog.Error(fmt.Sprintf("Error parsing CSV input to Customer: %v", err)) 
            continue
        }
        res = append(res, worker)
    }

    return res, nil
}

func CSVRowToWorker(row []string) (Worker, error) {
    var res Worker
    var err error

    if len(row) != VALID_CSV_WORKER_SIZE {
        return res, &WorkerCsvError{len(row)}
    }

    res.id, err = strconv.Atoi(row[0])
    if err != nil { return res, err }

    res.first_name = row[1]
    res.last_name = row[2]

    res.orders_per_hour, err = strconv.Atoi(row[5])
    if err != nil { return res, err }

    res.salary_per_hour, err = strconv.Atoi(row[6])
    if err != nil { return res, err }


    tmp, err := CSVSplitSerial2Int(row[3], "{}", ";")
    if err != nil { return res, err }
    res.work_days, err = IntToWeekday(tmp)
    if err != nil { return res, err }


    res.work_hours, err = CSVSplitSerial2Int(row[4], "{}", ";")
    if err != nil { return res, err }


    return res, nil
}

func (w *Worker) GetDailySalary() int {
    return len(w.work_hours) * w.salary_per_hour 
}

func (w *Worker) CheckIfWorking(t time.Time) bool {
    a := slices.Contains(w.work_hours, t.Hour())
    b := slices.Contains(w.work_days, t.Weekday())

    return a && b
}

func (w *Worker) Work(db Database) error {
    slog.Debug("Entering work function.")

    if !w.CheckIfWorking(time.Now())  {
        slog.Debug("Worker not working at this hour.")
        return nil
    }
    
    orders, err := db.GetOpenOrders()
    slog.Debug(fmt.Sprintf("Received %d open orders.", len(orders)))
    if errors.Is(err, ErrorEmptyOrdersList) {
        slog.Debug(fmt.Sprint("GetOpenOrders did not return any orders: ", err))
        return err
    } else if err != nil { 
        slog.Debug(fmt.Sprint("GetOpenOrders failed!", err))
        return err
    }

    for i := 0; i < w.orders_per_hour; i++ {
        order, err := orders.Pop()
        slog.Debug(fmt.Sprint(order))
        if errors.Is(err, ErrorEmptyOrdersList) {
            break
        }

        err = db.UpdateOrder(order)
        if err != nil {
            slog.Debug("UpdateOrder failed!")
            return err
        }
    }
    
    return nil
}
