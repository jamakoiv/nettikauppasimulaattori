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
    {10,  "Erkki",    "Nyrhinen"   , 6,  500, 0.15, []int{1,2} },
    {11,  "Jaana",    "Lahtinen"   , 7,  250, 0.25, []int{2} },   
    {12,  "Toni",     "Kuusisto"   , 8,  50 , 0.10, []int{1} },    
    {13,  "Tero",     "Teronen"    , 9,  25 , 0.20, []int{3} },    
    {14,  "Liisa",    "Peronen"    , 9,  150, 0.25, []int{4} },     
    {22,  "Laura",    "Mukka"      , 12, 150, 0.10, []int{3} },       
    {24,  "Sakari",   "Herkola"    , 12, 250, 0.15, []int{1} },
    {31,  "Kalevi",   "Sorsa"      , 12, 150, 0.20, []int{1} },     
    {33,  "Mauno",    "Koivisto"   , 14, 700, 0.05, []int{2} },  
    {34,  "Tarja",    "Kekkonen"   , 14, 350, 0.30, []int{2} },   
    {120, "Hertta",   "Kuusisto"   , 14, 150, 0.15, []int{4,3} },  
    {121, "Sari",     "Jokunen"    , 14, 250, 0.20, []int{2,1} },      
    {122, "Kaarina",  "Erkylä"     , 17, 150, 0.10, []int{2} },    
    {123, "Pasi",     "Sarasti"    , 17, 150, 0.20, []int{1} },
    {200, "Matti",    "Välimäki"   , 17, 100, 0.10, []int{4} },
    {201, "Matias",   "Honkamaa"   , 17, 100, 0.30, []int{3} },
    {202, "Mirva",    "Holma"      , 18, 150, 0.20, []int{3} },
    {203, "Sari",     "Karjalainen", 18, 150, 0.20, []int{4} },
    {204, "Teija",    "Laakso"     , 18, 500, 0.30, []int{2} },
    {205, "Mika",     "Rampa"      , 20, 500, 0.05, []int{2} },
    {206, "Antti",    "Vettenranta", 20, 50 , 0.20, []int{1} },
    {207, "Anri",     "Lindström"  , 20, 50 , 0.10, []int{1,2} },
    {208, "Taina",    "Vilkuna"    , 20, 150, 0.20, []int{1} },
    {209, "Sami",     "Turunen"    , 21, 750, 0.10, []int{2} },
    {210, "Marjo",    "Tiirikka"   , 21, 25 , 0.20, []int{3} },
    {211, "Jirina",   "Alanko"     , 21, 150, 0.20, []int{4,3} },
    {212, "Kasper",   "Kukkonen"   , 21, 150, 0.05, []int{4} },
    {213, "Karina",   "Tiihonen"   , 22, 150, 0.10, []int{2} },
    {214, "Pauliina", "Kampuri"    , 22, 50 , 0.20, []int{1,2} },
    {215, "Nelli",    "Numminen"   , 22, 25 , 0.25, []int{2} },
    {216, "Anna",     "Schroderus" , 22, 150, 0.20, []int{1} },
    {217, "Sabrina",  "Bqain"      , 23, 300, 0.10, []int{4,2} },  
    {218, "Tara",     "Junker"     , 23, 150, 0.10, []int{4} },

    {219, "Milan",    "Kundera"    , 17, 30 , 0.25, []int{1} },
    {220, "John",     "Kelleher"   , 18, 200, 0.25, []int{2} },
    {221, "Brendan",  "Tierney"    , 18, 250, 0.30, []int{2} },
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
    {232, "Jouko",    "Pukkila"    , 21, 100, 0.35, []int{1} },
    {233, "Ilari",    "Männistö"   , 3 , 30 , 0.15, []int{3,2} },
    {234, "Iiri",     "Salomaa"    , 18, 50 , 0.25, []int{4} },
    {235, "Mikko",    "Akkola"     , 17, 200, 0.25, []int{2} },
    {236, "Maijastiina", "Vilenius", 9 , 50 , 0.15, []int{1} },
    {237, "Pasi",     "Degerström" , 7 , 300, 0.25, []int{2} },
    {238, "Sippo",    "Mentunen"   , 5 , 30 , 0.15, []int{3} },
    {239, "Katimaria","Mustajärvi" , 18, 25 , 0.20, []int{4} },
    {240, "Petteri",  "Oja"        , 18, 50 , 0.30, []int{1} },
    {241, "Jouko",    "Pukkila"    , 17, 100, 0.20, []int{2} },
    {242, "Timo",     "Ronkainen"  , 15, 75 , 0.25, []int{1} },
    {243, "Ossi",     "Hiekkala"   , 21, 100, 0.25, []int{2} },
    {244, "Yrjänä",   "Ermala"     , 18, 20 , 0.15, []int{3} },

    {245, "Ville",    "Similä"     , 15, 20 , 0.10, []int{3} },
    {246, "Mervi",    "Vuorela"    , 18, 75 , 0.10, []int{4} },
    {247, "Viljami",  "Puustinen"  , 19, 25 , 0.15, []int{1} },
    {248, "Linda",    "Fredrikson" , 20, 50 , 0.15, []int{1} },
    {249, "Charles",  "Mingus"     , 20, 400, 0.10, []int{2} },
    {250, "John Lee", "Hooker"     , 22, 100, 0.15, []int{1} },
    {251, "Billy",    "Gibbons"    , 22, 300, 0.10, []int{2} },
    {252, "Frank",    "Bread"      , 15, 30 , 0.10, []int{3} },
    {253, "Haruki",   "Murakami"   , 16, 50 , 0.15, []int{4} },
    {254, "Ian",      "Winwood"    , 16, 30 , 0.10, []int{2} },
    {255, "Paul",     "Branningan" , 17, 20 , 0.15, []int{1} },

    // ChatGTP to the rescue
    {256, "Alice",    "Johnson",    8,  40, 0.15, []int{3}},
    {257, "David",    "Smith",      18, 600,0.10, []int{2}},
    {258, "Eva",      "Andersson",  17, 35, 0.10, []int{1}},
    {259, "Michael",  "Williams",   18, 70, 0.15, []int{4}},
    {260, "Sophia",   "Brown",      14, 55, 0.15, []int{4}},
    {261, "Oliver",   "Davis",      15, 25, 0.10, []int{1}},
    {262, "Emma",     "Garcia",     17, 100,0.10, []int{2}},
    {263, "James",    "Rodriguez",  18, 80, 0.15, []int{3}},
    {264, "Olivia",   "Martinez",   20, 70, 0.15, []int{4}},
    {265, "Liam",     "Lopez",      22, 55, 0.10, []int{1}},
    {266, "Sophie",   "Clark",      9,  30, 0.15, []int{3}},
    {267, "William",  "Moore",      12, 65, 0.10, []int{2}},
    {268, "Ava",      "Taylor",     16, 40, 0.10, []int{1}},
    {269, "Benjamin", "Johnson",    19, 75, 0.15, []int{4}},
    {270, "Charlotte", "Harris",    21, 50, 0.15, []int{4}},
    {271, "Henry",    "Davis",      8,  100,0.10, []int{1}},
    {272, "Olivia",   "Miller",     10, 50, 0.10, []int{2}},
    {273, "Emily",    "Wilson",     12, 20, 0.15, []int{3}},
    {274, "Michael",  "Hernandez",  14, 70, 0.15, []int{4}},
    {275, "Ella",     "Lewis",      16, 30, 0.10, []int{1}},
    {276, "Alexander", "Jackson",   18, 400,0.10, []int{2}},
    {277, "Elizabeth", "White",     20, 45, 0.15, []int{3}},
    {278, "Daniel",   "Young",      22, 60, 0.15, []int{4}},
    {279, "Sophia",   "Scott",      9,  35, 0.10, []int{1}},
    {280, "William",  "Harris",     11, 650,0.10, []int{2}},
    {281, "Amelia",   "King",       13, 20, 0.15, []int{3}},
    {282, "James",    "Lee",        15, 70, 0.15, []int{4}},
    {283, "Mia",      "Martin",     17, 50, 0.10, []int{1}},
    {284, "David",    "Walker",     19, 45, 0.10, []int{2}},
    {285, "Emma",     "Gonzalez",   21, 75, 0.15, []int{3}},
    {286, "Henry",    "Perez",      8,  100,0.15, []int{4}},
    {287, "Olivia",   "Hall",       10, 10, 0.10, []int{1}},
    {288, "Benjamin", "Lewis",      18, 100,0.10, []int{2}},
    {289, "Charlotte", "Collins",   14, 20, 0.15, []int{3}},
    {290, "Liam",     "Adams",      16, 50, 0.15, []int{4}},
    {291, "Emily",    "Russell",    18, 45, 0.10, []int{1}},
    {292, "Alexander", "Price",     20, 60, 0.10, []int{2}},
    {293, "Mia",      "Bennett",    22, 75, 0.15, []int{3}},
    {294, "Ava",      "Brooks",     9,  30, 0.15, []int{4}},
    {295, "Sophie",   "Morgan",     11, 65, 0.10, []int{1}},
    {296, "William",  "Hughes",     13, 40, 0.10, []int{2}},
    {297, "James",    "Kelly",      15, 70, 0.15, []int{3}},
    {298, "Ella",     "Parker",     17, 50, 0.15, []int{4}},
    {299, "Oliver",   "Simmons",    19, 45, 0.10, []int{1}},
    {300, "Emma",     "Foster",     21, 200,0.10, []int{2}},
    {301, "Daniel",   "Cooper",     22, 75, 0.15, []int{3}},
    {302, "Mia",      "Barnes",     8,  150,0.15, []int{4}},
    {303, "Liam",     "Butler",     10, 65, 0.10, []int{1}},
    {304, "Ava",      "Simmons",    11, 40, 0.10, []int{2}},
    {305, "Sophia",   "Ross",       13, 70, 0.15, []int{3}},
    {306, "William",  "Jenkins",    18, 50, 0.15, []int{4}},
    {307, "Olivia",   "Ward",       17, 45, 0.10, []int{1}},
    {308, "James",    "Griffin",    19, 60, 0.10, []int{2}},
    {309, "Ella",     "West",       21, 75, 0.15, []int{4}},
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

    base_spread := 5
    rand_spread := rand.Intn(5)

    prob := calc_probability(hour, 
        customer.base_purchase_probability,
        customer.most_active,
        base_spread + rand_spread)

    return prob
}

func (customer *Customer) Shop(products []Product) (Order, error) {

    // order := new(Order)
    var order Order
    order.init()

    // Check if customer wants to shop at this time.
    if !(rand.Float64() < customer.ChanceToShop()) {
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
