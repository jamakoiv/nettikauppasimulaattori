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

// TODO: GetOpenOrders and updateOrder should rather be in orders.go
// Also define PopOrder for list of orders, ErrorOrderEmpty if list is empty
// and use those rather than 'if len(orders)' in Work.


func GetOpenOrders(ctx context.Context, client *bigquery.Client) ([]Order, error) {
    // TODO: Move ids to config file somewhere.
    project_id := "nettikauppasimulaattori"
    dataset_id := "store_operational"
    table_id := "orders"

    sql := fmt.Sprintf("SELECT id, customer_id, delivery_type, status, order_placed FROM `%s.%s.%s` WHERE status = %d",
        project_id, dataset_id, table_id, ORDER_PENDING)
    slog.Debug(sql)

    slog.Debug("Creating result var.")
    var res []Order

    slog.Debug("Sending query.")
    q := client.Query(sql)
    job, err := q.Run(ctx)
    if err != nil { 
        slog.Error(fmt.Sprint(err))
        return res, err }

    slog.Debug("Wait query.")
    status, err := job.Wait(ctx)
    if err != nil { 
        slog.Error(fmt.Sprint(err))
        return res, err 
    }

    slog.Debug("Check status.")
    if status.Err() != nil { 
        slog.Error(fmt.Sprint(err))
        return res, status.Err()
    }

    slog.Debug("Get iterator.")
    it, err := job.Read(ctx)
    if err != nil { 
        slog.Error(fmt.Sprint(err))
        return res, err
    }

    
    slog.Debug("Parse results.")
    for {
        var tmp OrderReceiver
        if it.Next(&tmp) == iterator.Done { break }
        // fmt.Printf("%d: %T\n", tmp.ID, tmp.ID)
        res = append(res, ConvertOrderReceiverToOrder(tmp))
    }

    // TODO: Add error if res has zero length.

    return res, nil
}

func ConvertOrderReceiverToOrder(o OrderReceiver) Order {
    var res Order

    res.id = o.id
    res.customer_id = o.customer_id
    res.order_placed = o.order_placed
    res.delivery_type = o.delivery_type
    res.status = o.status
        
    return res
}


func UpdateOrder(order Order, ctx context.Context, client *bigquery.Client) error {
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
        order.id)
    slog.Debug(sql)

    q := client.Query(sql)
    job, err := q.Run(ctx)
    if err != nil { return err }

    status, err := job.Wait(ctx)
    if err != nil { return err }
    if status.Err() != nil { return status.Err() }

    return nil
}
