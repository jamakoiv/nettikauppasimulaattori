package main

import (
    "fmt"
    "time"
    "errors"
    "math"
    "math/rand"
    //"cloud.google.com/go/bigquery"
)

type Customer struct {
    id int
    first_name string
    last_name string
    most_active int
    price_point int
    base_purchase_probability float64
}

type Product struct {
    id int
    name string
    wholesale_price int
    price int
    vat float64
}

type Order struct {
    id uint64
    customer_id int
    items []Product
    delivery_type int
    status int
}

const (
    ORDER_PENDING = iota
    ORDER_SHIPPED = iota
    ORDER_EMPTY = iota
)

const (
    SHIP_TO_CUSTOMER = iota
    COLLECT_FROM_STORE = iota
)

var Customers = []Customer{ {10, "Erkki", "Nyrhinen", 18, 20, 0.05},
                            {11, "Jaana", "Lahtinen", 21, 50, 0.05},                     
                            {12, "Toni", "Kuusisto", 22, 30, 0.02},                     
                            {13, "Tero", "Teronen", 17, 100, 0.02},                     
                            {14, "Liisa", "Peronen", 12, 5, 0.10},                     
                            {22, "Laura", "Mukka", 18, 5, 0.10},                     
                            {24, "Sakari", "Herkola", 12, 25, 0.03},                     
                            {31, "Kalevi", "Sorsa", 18, 30, 0.02},                     
                            {33, "Mauno", "Koivisto", 18, 100, 0.02},                     
                            {34, "Tarja", "Kekkonen", 20, 30, 0.03},                     
                            {120,"Hertta", "Kuusisto", 21, 15, 0.07},                     
                            {121,"Sari", "Jokunen", 7, 50, 0.01},                     
                            {122,"Kaarina", "Erkylä", 8, 20, 0.02},                     
                            {123,"Pasi", "Sarasti", 9, 100, 0.04},                     
                          }

var Products = []Product{   {1001, "Pirkka olut 24-pak.", 10, 25, 0.24},
                            {1002, "Pirkka olut 6-pak.", 3, 8, 0.24},
                            {2001, "Raspberry Pi 4 4GB", 40, 80, 0.24},
                            {2002, "Raspberry Pi 4 8GB", 50, 100, 0.24},
                            {2003, "VHS-kasetteja 10-pak", 5, 8, 0.24},
                            {2004, "C-kasetteja 10-pak", 3, 8, 0.24},
                            {2005, "LCD-televisio", 150, 300, 0.24},
                            {2006, "Iso LCD-televisio", 200, 400, 0.24},

                            {3001, "Ruisleipä", 1, 3, 0.14},
                            {3002, "Rasvaton Maito 1L", 1, 2, 0.14},
                            {3003, "Kevytmaito 1L", 1, 2, 0.14},
                            {3004, "Täysmaito 1L", 1, 2, 0.14},
                            {3005, "Kauraleipä", 1, 3, 0.14},
                            {3006, "Mysliä 1kg", 1, 4, 0.14},
                            {3007, "Perunoita 1kg", 1, 2, 0.14},

                            {4000, "Silmarillion, J.R.R Tolkien", 10, 25, 0.10},
                            {4001, "Tabu, Timo Mukka", 5, 15, 0.10},
                            {4002, "Robottien aamunkoitto, Isaac Asimov", 10, 15, 0.10},
                            {4003, "Holmenkollen, Matti Hagelberg", 15, 30, 0.10},
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

func (customer *Customer) ChanceToShop() float64 {
    hour := time.Now().Hour()

    prob := calc_probability(hour, 
        customer.base_purchase_probability,
        customer.most_active,
        rand.Intn(10))

    return prob
}

func (customer *Customer) Shop(products []Product) (*Order, error) {

    order := new(Order)
    order.init()
    order.customer_id = customer.id

    if !(rand.Float64() < customer.ChanceToShop()) {
        return order, errors.New("Order empty.")
    }

    // Customer picks randomly how many and which products to buy.
    n := rand.Intn(10)+1 
    for i := 0; i < n; i++ {
        order.AddItem(products[rand.Intn(len(products))])
    }
    order.status = ORDER_PENDING

    return order, nil
}


func (order *Order) init() {
    order.id = rand.Uint64()  // Foolishly hope we don't get two same order IDs.
    order.status = ORDER_EMPTY
    order.delivery_type = rand.Intn(2)
}

func (order *Order) AddItem(item Product) {
    order.items = append(order.items, item) 
}

// Satisfy Stringer-interface.
func (order *Order) String() string {
    if order.status == ORDER_EMPTY { return "" }

    var str string = fmt.Sprintf("Order %v\n--------------\n", order.id)

    for _, item := range order.items {
        str = str + fmt.Sprintf("%v: %v\n", order.customer_id, item.name)
    }
    return str
}

//func (o *Order) Send(client bigquery.Client) {
func (o *Order) Send() {

}


func main() {
    
    for _, customer := range Customers {
        order, err := customer.Shop(Products)
        if err == nil {
            fmt.Println(order)
            order.Send()
        } else {
            // fmt.Println(err)
            continue
        }
    }
}
