package nettikauppasimulaattori

import (
	"context"
	"slices"
	"time"

	"cloud.google.com/go/bigquery"
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

func (w *Worker) Work(ctx context.Context, client *bigquery.Client) error {
    slog.Debug("Entering work function.")

    if !w.CheckIfWorking(time.Now())  {
        slog.Debug("Worker not working at this hour.")
        return nil
    }
    
    orders, err := GetOpenOrders(ctx, client)
    if err != nil { 
        slog.Debug("GetOpenOrders failed!")
        return err
    }

    order := orders[0]
    for i := 0; i < w.orders_per_hour; i++ {
        err = UpdateOrder(order, ctx, client)
        if err != nil {
            slog.Debug("UpdateOrder failed!")
            return err
        }

        if len(orders) >= 2 {
            orders = orders[1:]
            order = orders[0]
        } else {
            break
        }
    }
    
    return nil
}
