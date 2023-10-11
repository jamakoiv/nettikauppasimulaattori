package nettikauppasimulaattori

/*
    Various utility functions used here and there.
*/

import (
    "fmt"
    "time"
    "log/slog"
)

type Settings struct {
    project_id string
    dataset_id string
    orders_table_id string
    order_items_table_id string

    timezone string
}

// func LoadSettings() Settings {
//     return nil
// }

func NowInTimezone(timezone string) time.Time {
    tz, err := time.LoadLocation(timezone)

    if err != nil {
        err_str := fmt.Sprintf("Error getting timezone 'time.LoadLocation(%s'): %s", 
            timezone, err)
        slog.Error(err_str)

        return time.Now()
    } 

    return time.Now().In(tz)
}

func Time2SQLDatetime(t time.Time) string {
    res := fmt.Sprintf("%d-%d-%d %d:%d:%d",
        t.Year(), t.Month(), t.Day(),
        t.Hour(), t.Minute(), t.Second())

    return res
}

func Time2SQLDate(t time.Time) string {
    res := fmt.Sprintf("%d-%d-%d", t.Year(), t.Month(), t.Day())

    return res
}
