package nettikauppasimulaattori

import (
	"context"
	"fmt"
	"math/rand"
	"slices"
	"time"

	"cloud.google.com/go/bigquery"
	"golang.org/x/exp/slog"
	"google.golang.org/api/iterator"
)

type Worker struct {
    id int
    first_name string
    last_name string
    work_days []time.Weekday
    work_hours []int
    salary_per_hour int
}

var Workers = []Worker{
    {10, "Seppo",   "Laukkanen",    
        []time.Weekday{ time.Monday, 
                        time.Tuesday, 
                        time.Wednesday, 
                        time.Thursday},
        []int{8,9,10,11,13,14,15},   
        15},

    {20, "Kari",    "Freedman",    
        []time.Weekday{ time.Monday, 
                        time.Tuesday, 
                        time.Wednesday, 
                        time.Thursday,
                        time.Friday},
        []int{9,10,11,13,14,15},
        14},
}

func (w *Worker) String() string {

    return ""
}

func (w *Worker) GetDailySalary() int {
    return len(w.work_hours) * w.salary_per_hour 
}

func (w *Worker) CheckIfWorking(t time.Time) bool {
    slog.Debug(fmt.Sprintf("Weekday: %v, Hour %v\n", t.Weekday(), t.Hour()))

    a := slices.Contains(w.work_hours, t.Hour())
    b := slices.Contains(w.work_days, t.Weekday())

    return a && b
}

func (w *Worker) Work(  order_id int,
                        settings Settings, 
                        ctx context.Context, 
                        client *bigquery.Client) error {
    if !w.CheckIfWorking(NowInTimezone(settings.timezone)) { 
        fmt.Println("Worker not working at current time.")
        return nil 
    }

    err := UpdateOrder(order_id, settings, ctx, client)
    if err != nil { return err }
    
    return nil
}

func GetOpenOrders( settings Settings, 
                    ctx context.Context, 
                    client *bigquery.Client) ([]int, error) {

    sql := fmt.Sprintf("SELECT id FROM `%s.%s.%s` WHERE status = %d",
        settings.project_id,
        settings.dataset_id,
        settings.orders_table_id,
        ORDER_PENDING)
    slog.Debug(sql)

    var res []int
    q := client.Query(sql)
    job, err := q.Run(ctx)
    if err != nil { return res, err }

    status, err := job.Wait(ctx)
    if err != nil { return res, err }
    if status.Err() != nil { return res, status.Err() }

    type order_id struct { ID int } // Need this extra struct for receiving single int...
    it, err := job.Read(ctx)
    if err != nil { return res, err }
    for {
        var tmp order_id
        if it.Next(&tmp) == iterator.Done { break }
        // fmt.Printf("%d: %T\n", tmp.ID, tmp.ID)
        res = append(res, tmp.ID)
    }

    return res, nil
}

func UpdateOrder(   order_id int, 
                    settings Settings, 
                    ctx context.Context, 
                    client *bigquery.Client) error {

    t := NowInTimezone(settings.timezone)

    sql := fmt.Sprintf("UPDATE `%s.%s.%s` SET status = %d, shipping_date = \"%s\", tracking_number = %d WHERE id = %d",
        settings.project_id, 
        settings.dataset_id,
        settings.orders_table_id, 
        ORDER_SHIPPED,
        Time2SQLDate(t), 
        rand.Int(),
        order_id)
    slog.Debug(sql)

    q := client.Query(sql)
    job, err := q.Run(ctx)
    if err != nil { return err }

    status, err := job.Wait(ctx)
    if err != nil { return err }
    if status.Err() != nil { return status.Err() }

    slog.Info(fmt.Sprintf("Worker packed and shipper order %d", order_id))

    return nil
}
