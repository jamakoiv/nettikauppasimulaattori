package nettikauppasimulaattori

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gorhill/cronexpr"
	"golang.org/x/exp/slog"
)

type Worker struct {
	id              int
	first_name      string
	last_name       string
	cron_line       string
	orders_per_hour int
	salary_per_hour int
}

var (
	RUN_FREQUENCY         string = "1h"
	VALID_CSV_WORKER_SIZE int    = 6
)

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
	if err != nil {
		return res, err
	}

	reader := csv.NewReader(f)
	reader.LazyQuotes = false
	rows, err := reader.ReadAll()
	if err != nil {
		return res, err
	}

	for _, row := range rows {
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
	if err != nil {
		return res, err
	}

	res.first_name = row[1]
	res.last_name = row[2]
	res.cron_line = row[3]

	res.orders_per_hour, err = strconv.Atoi(row[4])
	if err != nil {
		return res, err
	}

	res.salary_per_hour, err = strconv.Atoi(row[5])
	if err != nil {
		return res, err
	}

	return res, nil
}

func (w *Worker) GetDailySalary() int {
	slog.Warn(fmt.Sprintf("Unimplemented at this time."))
	return 0
}

func (w *Worker) CheckIfWorking(t time.Time) bool {
	run_frequency, err := time.ParseDuration(RUN_FREQUENCY)
	if err != nil {
		slog.Error(fmt.Sprintf("Error parsing duration string %s", RUN_FREQUENCY))
	}

	slog.Debug(w.cron_line)
	cron_expr := cronexpr.MustParse(w.cron_line)
	next_run := cron_expr.Next(t)

	until_next_run := next_run.Sub(t)
	if until_next_run <= run_frequency {
		return true
	} else {
		return false
	}
}

func (w *Worker) Work(db Database) error {
	slog.Debug("Entering work function.")

	if !w.CheckIfWorking(time.Now()) {
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
