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
    work_days []int
    work_hours []int
    salary_per_hour int
}

var Workers = []Worker{}

func (w *Worker) GetDailySalary() int {
    return len(w.work_hours) * w.salary_per_hour 
}

func (w *Worker) CheckIfWorking(t time.Time) bool {
    a := slices.Contains(w.work_hours, t.Hour())
    b := slices.Contains(w.work_days, t.Hour())

    return a && b
}

func (w *Worker) Work(ctx context.Context, client *bigquery.Client) error {
    if !w.CheckIfWorking(time.Now())  {
        return nil
    }
    
    orders, err := GetOpenOrders(ctx, client)
    if err != nil { return err }

    err = UpdateOrder(rand.Intn(len(orders)), ctx, client)
    if err != nil { return err }
    
    return nil
}

func GetOpenOrders(ctx context.Context, client *bigquery.Client) ([]int, error) {
    // TODO: Move ids to config file somewhere.
    project_id := "nettikauppasimulaattori"
    dataset_id := "store_operational"
    table_id := "orders"

    sql := fmt.Sprintf("SELECT id FROM `%s.%s.%s` WHERE status = %d",
        project_id, dataset_id, table_id, ORDER_PENDING)
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

func UpdateOrder(order_id int, ctx context.Context, client *bigquery.Client) error {
    project_id := "nettikauppasimulaattori"
    dataset_id := "store_operational"
    table_id := "orders"

    now, _ := nowInTimezone("Europe/Helsinki")

    sql := fmt.Sprintf("UPDATE `%s.%s.%s` SET status = %d, shipping_date = \"%s\", tracking_number = %d WHERE id = %d",
        project_id, 
        dataset_id,
        table_id, 
        ORDER_SHIPPED,
        Time2SQLDate(now), 
        rand.Int(),
        order_id)
    slog.Debug(sql)

    q := client.Query(sql)
    job, err := q.Run(ctx)
    if err != nil { return err }

    status, err := job.Wait(ctx)
    if err != nil { return err }
    if status.Err() != nil { return status.Err() }

    return nil
}
