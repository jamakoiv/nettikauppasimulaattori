package nettikauppasimulaattori

import (
	"fmt"
	"strconv"
	"strings"
        "time"
)


func CSVRemoveWhitespace(row []string) []string {
    res := make([]string, len(row))

    for i := 0; i < len(row); i++ {
        s := strings.ReplaceAll(row[i], " ", "")
        res[i] = strings.ReplaceAll(s, "\t", "")
    }

    return res
}

func CSVSplitSerial2Int(s string, trim string, delim string) ([]int, error) {
    var res []int

    s = strings.Trim(s, trim)
    tmp := strings.Split(s, delim)

    for _, val_str := range tmp {
        val, err := strconv.Atoi(val_str)
        if err != nil { return res, err }

        res = append(res, val) 
    }

    return res, nil
}

type WeekdayOutOfRangeError struct {
    input int
}
func (e *WeekdayOutOfRangeError) Error() string {
    return fmt.Sprintf("Input value %v not in accepted range of 0-6.", e.input)
}

func IntToWeekday(input []int) ([]time.Weekday, error) {
    var res []time.Weekday

    for _, val := range input {
        if !(val >= 0 && val <= 7) {
            return res, &WeekdayOutOfRangeError{val}
        }

        res = append(res, time.Weekday(val))
    }

    return res, nil
}
