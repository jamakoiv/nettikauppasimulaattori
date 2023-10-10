package nettikauppasimulaattori

// TODO: Break into separate files:
// customers.go, orders.go, products.go, run.go etc.


import (
	"errors"
	"math"
	"math/rand"
	"time"
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


var Customers = []Customer{
    {10,  "Erkki",    "Nyrhinen"   , 6,  150, 0.15 },
    {11,  "Jaana",    "Lahtinen"   , 7,  150, 0.25 },   
    {12,  "Toni",     "Kuusisto"   , 8,  150, 0.10 },    
    {13,  "Tero",     "Teronen"    , 9,  150, 0.20 },    
    {14,  "Liisa",    "Peronen"    , 9,  150, 0.25 },     
    {22,  "Laura",    "Mukka"      , 12, 150, 0.10 },       
    {24,  "Sakari",   "Herkola"    , 12, 150, 0.15 },
    {31,  "Kalevi",   "Sorsa"      , 12, 150, 0.20 },     
    {33,  "Mauno",    "Koivisto"   , 14, 150, 0.05 },  
    {34,  "Tarja",    "Kekkonen"   , 14, 150, 0.30 },   
    {120, "Hertta",   "Kuusisto"   , 14, 150, 0.15 },  
    {121, "Sari",     "Jokunen"    , 14, 150, 0.20 },      
    {122, "Kaarina",  "Erkylä"     , 17, 150, 0.10 },    
    {123, "Pasi",     "Sarasti"    , 17, 150, 0.20 },
    {200, "Matti",    "Välimäki"   , 17, 150, 0.10 },
    {201, "Matias",   "Honkamaa"   , 17, 150, 0.30 },
    {202, "Mirva",    "Holma"      , 18, 150, 0.20 },
    {203, "Sari",     "Karjalainen", 18, 150, 0.20 },
    {204, "Teija",    "Laakso"     , 18, 150, 0.30 },
    {205, "Mika",     "Rampa"      , 20, 150, 0.05 },
    {206, "Antti",    "Vettenranta", 20, 150, 0.20 },
    {207, "Anri",     "Lindström"  , 20, 150, 0.10 },
    {208, "Taina",    "Vilkuna"    , 20, 150, 0.20 },
    {209, "Sami",     "Turunen"    , 21, 150, 0.10 },
    {210, "Marjo",    "Tiirikka"   , 21, 150, 0.20 },
    {211, "Jirina",   "Alanko"     , 21, 150, 0.20 },
    {212, "Kasper",   "Kukkonen"   , 21, 150, 0.05 },
    {213, "Karina",   "Tiihonen"   , 22, 150, 0.10 },
    {214, "Pauliina", "Kampuri"    , 22, 150, 0.20 },
    {215, "Nelli",    "Numminen"   , 22, 150, 0.25 },
    {216, "Anna",     "Schroderus" , 22, 150, 0.20 },
    {217, "Sabrina",  "Bqain"      , 23, 150, 0.10 },  
    {218, "Tara",     "Junker"     , 23, 150, 0.10 },
}


var Products = []Product{   
    {1001, "Pirkka olut 24-pak.", 10, 25, 0.24},
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

    // Check if customer wants to shop at this time.
    if !(rand.Float64() < customer.ChanceToShop()) {
        return order, errors.New("Order empty.")
    }

    // Customer picks randomly how many and which products to buy.
    n := rand.Intn(10)+1 
    for i := 0; i < n; i++ {
        order.AddItem(products[rand.Intn(len(products))])
    }
    order.status = ORDER_PENDING
    order.customer_id = customer.id

    return order, nil
}
