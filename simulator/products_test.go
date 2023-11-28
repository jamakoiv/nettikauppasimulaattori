package nettikauppasimulaattori

import (
	"fmt"
	"reflect"
	"testing"
        "strings"
)

var TestProducts = []Product{   
    {1001, "Pirkka olut 24-pak.", 10, 25, 0.24},
    {2001, "Raspberry Pi 4 4GB", 40, 80, 0.24},
    {3001, "Ruisleipä", 1, 3, 0.14},
    {4000, "Silmarillion, J.R.R Tolkien", 10, 25, 0.10},
}   

func TestFilterProductsByCategory(t *testing.T) {
    target := []Product{
        {1001, "Pirkka olut 24-pak.", 10, 25, 0.24},
        {2001, "Raspberry Pi 4 4GB", 40, 80, 0.24},
    }
    target_categories := []int{ALCOHOL, ELECTRONICS}

    res := FilterProductsByCategory(TestProducts, target_categories)

    if !(reflect.DeepEqual(res, target)) {
        var res_ids strings.Builder

        for _, prod := range res {
            res_ids.WriteString(fmt.Sprintf("%v, ", prod.id))
        }

        t.Fatalf("Wanted products '1001, 2001', got products '%s'", &res_ids)
    }
}

func TestFilterProductsByPrice(t *testing.T) {
    target := []Product{
        {1001, "Pirkka olut 24-pak.", 10, 25, 0.24},
        {3001, "Ruisleipä", 1, 3, 0.14},
        {4000, "Silmarillion, J.R.R Tolkien", 10, 25, 0.10},
    }
    target_max_price := 50

    res := FilterProductsByPrice(TestProducts, target_max_price)

    if !(reflect.DeepEqual(res, target)) {
        var res_ids strings.Builder

        for _, prod := range res {
            res_ids.WriteString(fmt.Sprintf("%v, ", prod.id))
        }

        t.Fatalf("Wanted products '1001, 3001, 4000', got products '%s'", &res_ids)
    }
}
