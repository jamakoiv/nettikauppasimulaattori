package nettikauppasimulaattori

import (
	"errors"
	"math"
	"math/rand"
	"time"

	"github.com/parquet-go/parquet-go"
)

type Customer struct {
	Id                        int     `parquet:"id"`
	First_name                string  `parquet:"first_name"`
	Last_name                 string  `parquet:"last_name"`
	Most_active               int     `parquet:"most_active"`
	Max_budget                int     `parquet:"max_budget"`
	Base_purchase_probability float64 `parquet:"purchase_probability"`
	Product_categories        []int   `parquet:"product_categories"`
}

func ImportCustomers(file string) ([]Customer, error) {
	res, err := parquet.ReadFile[Customer](file)
	if err != nil {
		return res, err
	}

	return res, nil
}

func calc_probability(x int, base_probability float64, target int, spread int) float64 {
	// Calculate probability which drops as we get further away from 'target'.
	// When x == target: prob -> base_probability.
	// When x == target +- spread: prob -> 0.

	mu := float64(base_probability) / float64(spread)
	res := base_probability - float64(math.Abs(float64(target-x))*mu)
	// fmt.Printf("res: %e\n", res)

	if res < 0 {
		res = 0
	}

	return res
}

func (customer *Customer) ChanceToShop(t time.Time, day_var float64, week_var float64) float64 {
	base_spread := 5
	rand_spread := rand.Intn(5)

	prob := calc_probability(t.Hour(),
		customer.Base_purchase_probability+day_var+week_var,
		customer.Most_active,
		base_spread+rand_spread)

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
		return order, errors.New("Order empty")
	}

	// Filter products by product category and customer budget.
	products = FilterProductsByCategory(products, customer.Product_categories)
	products = FilterProductsByPrice(products, customer.Max_budget)

	// Customer picks randomly how many and which products to buy.
	var money_remaining int
	n := rand.Intn(20) + 1
	for i := 0; i < n; i++ {
		order.AddItem(products[rand.Intn(len(products))])

		money_remaining = customer.Max_budget - order.TotalPrice()
		products = FilterProductsByPrice(products, money_remaining)
		if len(products) == 0 {
			break
		}
	}
	order.status = ORDER_PENDING
	order.customer_id = customer.Id

	return order, nil
}

func Default_ShoppingWeekdayVariation(now time.Time) float64 {
	variation := map[time.Weekday]float64{
		time.Monday:    -0.01,
		time.Tuesday:   -0.01,
		time.Wednesday: 0.02,
		time.Thursday:  0.01,
		time.Friday:    0.03,
		time.Saturday:  0.04,
		time.Sunday:    0.01,
	}

	return variation[now.Weekday()]
}

func Default_ShoppingWeekVariation(now time.Time) float64 {
	t_start := time.Date(2024, time.January, 23, 0, 0, 0, 0, time.UTC)
	week := int64(604800)  // One week in seconds
	ramp := float64(0.004) // How much shopping chance goes up per week.

	return float64((now.Unix()-t_start.Unix())/week) * ramp
}
