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
    max_budget int
    base_purchase_probability float64
    product_categories []int
}

var Customers = []Customer{
    {10,  "Erkki",    "Nyrhinen"   , 6,  150, 0.15, []int{1,2} },
    {11,  "Jaana",    "Lahtinen"   , 7,  150, 0.25, []int{2} },   
    {12,  "Toni",     "Kuusisto"   , 8,  150, 0.10, []int{1} },    
    {13,  "Tero",     "Teronen"    , 9,  150, 0.20, []int{3} },    
    {14,  "Liisa",    "Peronen"    , 9,  150, 0.25, []int{4} },     
    {22,  "Laura",    "Mukka"      , 12, 150, 0.10, []int{3} },       
    {24,  "Sakari",   "Herkola"    , 12, 150, 0.15, []int{1} },
    {31,  "Kalevi",   "Sorsa"      , 12, 150, 0.20, []int{1} },     
    {33,  "Mauno",    "Koivisto"   , 14, 150, 0.05, []int{2} },  
    {34,  "Tarja",    "Kekkonen"   , 14, 150, 0.30, []int{2} },   
    {120, "Hertta",   "Kuusisto"   , 14, 150, 0.15, []int{4,3} },  
    {121, "Sari",     "Jokunen"    , 14, 150, 0.20, []int{2,1} },      
    {122, "Kaarina",  "Erkylä"     , 17, 150, 0.10, []int{2} },    
    {123, "Pasi",     "Sarasti"    , 17, 150, 0.20, []int{1} },
    {200, "Matti",    "Välimäki"   , 17, 150, 0.10, []int{4} },
    {201, "Matias",   "Honkamaa"   , 17, 150, 0.30, []int{3} },
    {202, "Mirva",    "Holma"      , 18, 150, 0.20, []int{3} },
    {203, "Sari",     "Karjalainen", 18, 150, 0.20, []int{4} },
    {204, "Teija",    "Laakso"     , 18, 150, 0.30, []int{2} },
    {205, "Mika",     "Rampa"      , 20, 150, 0.05, []int{2} },
    {206, "Antti",    "Vettenranta", 20, 150, 0.20, []int{1} },
    {207, "Anri",     "Lindström"  , 20, 150, 0.10, []int{1,2} },
    {208, "Taina",    "Vilkuna"    , 20, 150, 0.20, []int{1} },
    {209, "Sami",     "Turunen"    , 21, 150, 0.10, []int{2} },
    {210, "Marjo",    "Tiirikka"   , 21, 150, 0.20, []int{3} },
    {211, "Jirina",   "Alanko"     , 21, 150, 0.20, []int{4,3} },
    {212, "Kasper",   "Kukkonen"   , 21, 150, 0.05, []int{4} },
    {213, "Karina",   "Tiihonen"   , 22, 150, 0.10, []int{2} },
    {214, "Pauliina", "Kampuri"    , 22, 150, 0.20, []int{1,2} },
    {215, "Nelli",    "Numminen"   , 22, 150, 0.25, []int{2} },
    {216, "Anna",     "Schroderus" , 22, 150, 0.20, []int{1} },
    {217, "Sabrina",  "Bqain"      , 23, 150, 0.10, []int{4,2} },  
    {218, "Tara",     "Junker"     , 23, 150, 0.10, []int{4} },

    {219, "Milan",    "Kundera"    , 17, 30 , 0.25, []int{1} },
    {220, "John",     "Kelleher"   , 18, 50 , 0.25, []int{2} },
    {221, "Brendan",  "Tierney"    , 18, 50 , 0.30, []int{2} },
    {222, "Kimmo",    "Pietiläinen", 21, 100, 0.15, []int{4} },
    {223, "Ethem",    "Alpaydin"   , 19, 100, 0.20, []int{1} },
    {224, "Petri",    "Hiltunen"   , 4,  50 , 0.30, []int{1} },
    {225, "Timo",     "Niemi"      , 11, 100, 0.15, []int{3} },
    {226, "Sallamari","Rantala"    , 11, 200, 0.15, []int{2} },
    {227, "Kaisa",    "Bertel"     , 16, 50 , 0.35, []int{2} },
    {228, "Riikka",   "Puumalainen", 16, 25 , 0.15, []int{4} },
    {229, "Kaisa",    "Herrala"    , 21, 50 , 0.15, []int{4} },
    {230, "Jaakko",   "Herrala"    , 19, 150, 0.20, []int{2} },
    {231, "Muura",    "Kaleva"     , 18, 250, 0.25, []int{2} },
    {232, "Jouko",    "Pukkila"    , 21, 150, 0.35, []int{1} },
    {233, "Ilari",    "Männistö"   , 3 , 30 , 0.15, []int{3,2} },
    {234, "Iiri",     "Salomaa"    , 18, 50 , 0.25, []int{4} },
    {235, "Mikko",    "Akkola"     , 17, 200, 0.25, []int{2} },
    {236, "Maijastiina", "Vilenius", 9 , 200, 0.15, []int{1} },
    {237, "Pasi",     "Degerström" , 7 , 300, 0.25, []int{2} },
    {238, "Sippo",    "Mentunen"   , 5 , 30 , 0.15, []int{3} },
    {239, "Katimaria","Mustajärvi" , 18, 25 , 0.20, []int{4} },
    {240, "Petteri",  "Oja"        , 18, 50 , 0.30, []int{1} },
    {241, "Jouko",    "Pukkila"    , 17, 100, 0.20, []int{2} },
    {242, "Timo",     "Ronkainen"  , 15, 200, 0.25, []int{5} },
    {243, "Ossi",     "Hiekkala"   , 21, 100, 0.25, []int{2} },
    {244, "Yrjänä",   "Ermala"     , 18, 100, 0.15, []int{3} },
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

    products = filterProducts(products, customer.product_categories)
    
    // Customer picks randomly how many and which products to buy.
    n := rand.Intn(10)+1 
    for i := 0; i < n; i++ {
        order.AddItem(products[rand.Intn(len(products))])
    }
    order.status = ORDER_PENDING
    order.customer_id = customer.id

    return order, nil
}
