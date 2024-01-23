package nettikauppasimulaattori

import (
    "testing"
    "time"
    "math"
)


func TestDefault_ShoppingWeekVariation(t *testing.T) {

    test_date := time.Date(2024, time.January, 30, 0, 0, 0, 0, time.UTC)
    target := 0.002
    allowed_error := 0.00001

    res := Default_ShoppingWeekVariation(test_date)
    
    if (math.Abs(res - target) > allowed_error) {
        t.Fatalf("Wanted %v, got %v, which is not within allowed error of %v",
            res, target, allowed_error)
    }
}
