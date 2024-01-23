package nettikauppasimulaattori

import (
	"errors"
	"fmt"
	"slices"
	"time"

	"golang.org/x/exp/slog"
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

var Workers = []Worker{
    {
    111, "Vesa", "Sisättö", 
    []time.Weekday{time.Monday, 
        time.Tuesday, 
        time.Wednesday, 
        time.Thursday,
        time.Friday},
    []int{8, 9, 10, 11, 13, 14, 15},
    4, 15},

    {222, "Seppo", "Lahnakainen", 
    []time.Weekday{time.Monday, 
        time.Tuesday, 
        time.Wednesday, 
        time.Thursday,
        time.Friday},
    []int{8, 9, 10, 11, 13, 14, 15},
    4, 15},

    {333, "Janne", "Virtanen", 
    []time.Weekday{time.Monday, 
        time.Tuesday, 
        time.Wednesday, 
        time.Thursday},
    []int{12, 13, 14, 15, 16},
    4, 15},

    {444, "Erkki", "Kolehmainen", 
    []time.Weekday{time.Monday, 
        time.Tuesday, 
        time.Wednesday, 
        time.Thursday,
        time.Friday},
    []int{9, 10, 11, 12, 13, 14, 15, 16},
    3, 15},

    {555, "Laura", "Kolehmainen", 
    []time.Weekday{time.Monday, 
        time.Tuesday, 
        time.Wednesday, 
        time.Thursday,
        time.Friday},
    []int{9, 10, 11, 12, 13, 14, 15, 16},
    3, 15},

    {555, "Laura", "Kolehmainen", 
    []time.Weekday{time.Monday, 
        time.Tuesday, 
        time.Wednesday, 
        time.Thursday,
        time.Friday},
    []int{9, 10, 11, 12, 13, 14, 15, 16},
    3, 15},

    {666, "Sanna", "Sörppö",
    []time.Weekday{time.Tuesday, 
        time.Wednesday, 
        time.Thursday,
        time.Friday,
        time.Saturday},
    []int{12, 13, 14, 15, 16, 17, 18, 19, 20},
    4, 15},

    {777, "Ville", "Korhonen",
    []time.Weekday{time.Monday, 
        time.Tuesday, 
        time.Wednesday, 
        time.Thursday,
        time.Friday},
    []int{12, 13, 14, 15, 16, 17, 18, 19, 20},
    3, 15},

    {888, "Kiire", "Apulainen",
    []time.Weekday{time.Saturday,
        time.Sunday},
    []int{12, 13, 14, 15, 16, 17, 18},
    3, 15},
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
