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
    if !w.CheckIfWorking(time.Now())  {
        return nil
    }
    
    orders, err := GetOpenOrders(ctx, client)
    if err != nil { return err }

    order_id := orders[0]
    for i := 0; i < w.orders_per_hour; i++ {
        err = UpdateOrder(order_id, ctx, client)
        if err != nil { return err }

        if len(orders) >= 2 {
            orders = orders[1:]
            order_id = orders[0]
        } else {
            break
        }
    }
    
    return nil
}

// TODO: GetOpenOrders and updateOrder should rather be in orders.go
// Also define PopOrder for list of orders, ErrorOrderEmpty if list is empty
// and use those rather than 'if len(orders)' in Work.


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

    // TODO: Add error if res has zero length.

    return res, nil
}

func UpdateOrder(order_id int, ctx context.Context, client *bigquery.Client) error {
    project_id := "nettikauppasimulaattori"
    dataset_id := "store_operational"
    table_id := "orders"

    now, _ := nowInTimezone("Europe/Helsinki")

    sql := fmt.Sprintf("UPDATE `%s.%s.%s` SET status = %d, shipping_date = \"%s\", last_modified = \"%s\", tracking_number = %d WHERE id = %d",
        project_id, 
        dataset_id,
        table_id, 
        ORDER_SHIPPED,
        Time2SQLDate(now), 
        Time2SQLDatetime(now), 
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
