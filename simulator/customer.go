package nettikauppasimulaattori

// TODO: Break into separate files:
// customers.go, orders.go, products.go, run.go etc.

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	"encoding/csv"
	"os"
	"strconv"

	"golang.org/x/exp/slog"
)

type Customer struct {
    id int
    first_name string
    last_name string
    most_active int
    max_budget int
    base_purchase_probability float64
    product_categories []int
}

var VALID_CSV_ROW_SIZE int = 7

type CustomerCsvError struct {
    length int
}

func (e *CustomerCsvError) Error() string {
    return fmt.Sprintf("Received CSV-row with %v elements. Should be %v elements.",
        e.length, VALID_CSV_ROW_SIZE)
}

func ReadCustomersCSV(file string) ([]Customer, error) {
    var res []Customer
    
    f, err := os.Open(file)
    if err != nil { return res, err }

    reader := csv.NewReader(f)
    rows, err := reader.ReadAll()
    if err != nil { return res, err }

    for _, row := range rows {
        customer, err := CSVRowToCustomer(row)
        if err != nil { 
            slog.Error(fmt.Sprintf("Error parsing CSV input to Customer: %v", err)) 
            continue
        }
        res = append(res, customer)
    }

    return res, nil
}

func CSVRowToCustomer(row []string) (Customer, error) {
    var res Customer
    var err error
    
    if len(row) != VALID_CSV_ROW_SIZE {
        return res, &CustomerCsvError{len(row)}
    }

    // Atoi fails if there are any whitespace chars.
    for i := range row {
        row[i] = strings.ReplaceAll(row[i], " ", "")
        row[i] = strings.ReplaceAll(row[i], "\t", "")
    }

    res.id, err = strconv.Atoi(row[0])
    if err != nil { return res, err }

    res.first_name = row[1]
    res.last_name = row[2]

    res.most_active, err = strconv.Atoi(row[3])
    if err != nil { return res, err }

    res.max_budget, err = strconv.Atoi(row[4])
    if err != nil { return res, err }

    res.base_purchase_probability, err = strconv.ParseFloat(row[5], 64)
    if err != nil { return res, err }

    tmp := strings.Trim(row[6], "{}")
    categories := strings.Split(tmp, ";")

    for _, cat := range categories {
        c, err := strconv.Atoi(cat)
        if err != nil { return res, err }

        res.product_categories = append(res.product_categories, c)
    }

    return res, nil
}


func calc_probability(x int, base_probability float64, target int, spread int) float64 {
    // Calculate probability which drops as we get further away from 'target'.
    // When x == target: prob -> base_probability.
    // When x == target +- spread: prob -> 0.

    mu := float64(base_probability)/float64(spread)
    res := base_probability - float64(math.Abs(float64(target-x))*mu)
    // fmt.Printf("res: %e\n", res)

    if res < 0 { res = 0 }

    return res
}

func (customer *Customer) ChanceToShop(t time.Time, day_var float64, week_var float64) float64 {

    base_spread := 5
    rand_spread := rand.Intn(5)

    prob := calc_probability(t.Hour(), 
        customer.base_purchase_probability + day_var + week_var,
        customer.most_active,
        base_spread + rand_spread)

    return prob
}

func (customer *Customer) Shop(products []Product) (Order, error) {

    // order := new(Order)
    var order Order
    order.init()

    now := time.Now()
    day_variation := Default_ShoppingWeekdayVariation(now)
    week_variation := Default_ShoppingWeekVariation(now)

    // Check if customer wants to shop at this time.
    if !(rand.Float64() < customer.ChanceToShop(now, day_variation, week_variation)) {
        return order, errors.New("Order empty.")
    }

    // Filter products by product category and customer budget.
    products = FilterProductsByCategory(products, customer.product_categories)
    products = FilterProductsByPrice(products, customer.max_budget)
    
    // Customer picks randomly how many and which products to buy.
    var money_remaining int
    n := rand.Intn(20)+1 
    for i := 0; i < n; i++ {
        order.AddItem(products[rand.Intn(len(products))])

        money_remaining = customer.max_budget - order.TotalPrice()
        products = FilterProductsByPrice(products, money_remaining)
        if len(products) == 0 { break }
    }
    order.status = ORDER_PENDING
    order.customer_id = customer.id

    return order, nil
}

func Default_ShoppingWeekdayVariation(now time.Time) float64 { 
    variation := map[time.Weekday]float64 {
        time.Monday: -0.01,
        time.Tuesday: -0.01,
        time.Wednesday: 0.02,
        time.Thursday: 0.01,
        time.Friday: 0.03,
        time.Saturday: 0.04,
        time.Sunday: 0.01,
    }

    return variation[now.Weekday()]
}

func Default_ShoppingWeekVariation(now time.Time) float64 {
    t_start := time.Date(2024, time.January, 23, 0, 0, 0, 0, time.UTC) 
    week := int64(604800) // One week in seconds
    ramp := float64(0.004) // How much shopping chance goes up per week.

    return float64((now.Unix() - t_start.Unix()) / week) * ramp 
}
