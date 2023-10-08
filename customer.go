package nettikauppasimulaattori

// TODO: Break into separate files:
// customers.go, orders.go, products.go, run.go etc.


import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
	"golang.org/x/exp/slog"
)

type MessagePublishedData struct {
	Message PubSubMessage
}
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// Define google cloud functions entry-point.
func init() {
	functions.CloudEvent("Run", Run)
}

func Run(ctx context.Context, ev event.Event) error {
    // TODO: Project_id hardcoded in two different places :(
    project_id := "nettikauppasimulaattori"

    slog.Info(fmt.Sprintf("Program started at %v", time.Now()))

    client, err := bigquery.NewClient(ctx, project_id)
    if err != nil { 
        slog.Error("Error creating BigQuery-client.") 
        return err
    }
    defer client.Close()

    orders_in_this_run := false
    for _, customer := range Customers {
        order, err := customer.Shop(Products)
        if err != nil { continue } // If order is empty.

        slog.Debug(fmt.Sprint(order))
        slog.Info(fmt.Sprintf("Sending order %d to BigQuery.", order.id))
        err = order.Send(ctx, client)
        if err != nil { 
            slog.Error(fmt.Sprintf("Error in sending order: %v", err))
        }

        orders_in_this_run = true
    }
    if !orders_in_this_run {
        slog.Info("No orders placed this time.")
    }
    return nil
}


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


var Customers = []Customer{
    {10,  "Erkki",    "Nyrhinen"   , 6,  150, 0.10 },
    {11,  "Jaana",    "Lahtinen"   , 7,  150, 0.20 },   
    {12,  "Toni",     "Kuusisto"   , 8,  150, 0.10 },    
    {13,  "Tero",     "Teronen"    , 9,  150, 0.20 },    
    {14,  "Liisa",    "Peronen"    , 9,  150, 0.20 },     
    {22,  "Laura",    "Mukka"      , 12, 150, 0.10 },       
    {24,  "Sakari",   "Herkola"    , 12, 150, 0.10 },
    {31,  "Kalevi",   "Sorsa"      , 12, 150, 0.20 },     
    {33,  "Mauno",    "Koivisto"   , 14, 150, 0.10 },  
    {34,  "Tarja",    "Kekkonen"   , 14, 150, 0.20 },   
    {120, "Hertta",   "Kuusisto"   , 14, 150, 0.10 },  
    {121, "Sari",     "Jokunen"    , 14, 150, 0.20 },      
    {122, "Kaarina",  "Erkylä"     , 17, 150, 0.10 },    
    {123, "Pasi",     "Sarasti"    , 17, 150, 0.20 },
    {200, "Matti",    "Välimäki"   , 17, 150, 0.10 },
    {201, "Matias",   "Honkamaa"   , 17, 150, 0.30 },
    {202, "Mirva",    "Holma"      , 18, 150, 0.20 },
    {203, "Sari",     "Karjalainen", 18, 150, 0.20 },
    {204, "Teija",    "Laakso"     , 18, 150, 0.30 },
    {205, "Mika",     "Rampa"      , 20, 150, 0.10 },
    {206, "Antti",    "Vettenranta", 20, 150, 0.20 },
    {207, "Anri",     "Lindström"  , 20, 150, 0.10 },
    {208, "Taina",    "Vilkuna"    , 20, 150, 0.20 },
    {209, "Sami",     "Turunen"    , 21, 150, 0.10 },
    {210, "Marjo",    "Tiirikka"   , 21, 150, 0.20 },
    {211, "Jirina",   "Alanko"     , 21, 150, 0.20 },
    {212, "Kasper",   "Kukkonen"   , 21, 150, 0.10 },
    {213, "Karina",   "Tiihonen"   , 22, 150, 0.10 },
    {214, "Pauliina", "Kampuri"    , 22, 150, 0.20 },
    {215, "Nelli",    "Numminen"   , 22, 150, 0.20 },
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


func (order *Order) init() {
    order.id = uint64(rand.Uint32())  // Foolishly hope we don't get two same order IDs.
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

func Now2SQLDatetime() string {
    // Return current time as SQL Datetime.
    t := time.Now()
    return fmt.Sprintf("%d-%d-%d %d:%d:%d",
        t.Year(), t.Month(), t.Day(),
        t.Hour(), t.Minute(), t.Second())
}

func (order *Order) Send(ctx context.Context, client *bigquery.Client) error {
    // TODO: Break creating the SQL-queries into separate functions.
    // TODO: Store project_id etc in separate config-file.

    project_id := "nettikauppasimulaattori"
    dataset_id := "store_operational"
    orders_table_id := "orders"
    order_items_table_id := "order_items"


    // TODO: guard against malicious inputs.
    order_sql := fmt.Sprintf("INSERT INTO `%s.%s.%s` VALUES ", 
        project_id, dataset_id, orders_table_id)
    order_sql = fmt.Sprintf("%s (%d, %d, %d, %d, \"%s\")", 
        order_sql, order.id, order.customer_id, 
        order.delivery_type, order.status, Now2SQLDatetime())

    items_sql := fmt.Sprintf("INSERT INTO `%s.%s.%s` VALUES ", 
        project_id, dataset_id, order_items_table_id)

    for _, item := range order.items {
        items_sql = fmt.Sprintf("%s (%d, %d),", items_sql, order.id, item.id)
    }
    items_sql = strings.TrimSuffix(items_sql, ",")

    // slog.Debug(order_sql)
    // slog.Debug(items_sql)

    queries := [2]string{order_sql, items_sql}
    for _, sql := range queries {
        q := client.Query(sql)
        // q.WriteDisposition = "WRITE_APPEND" // Error with "INSERT INTO..." statement.

        job, err := q.Run(ctx)
        if err != nil { return err }

        status, err := job.Wait(ctx)
        if err != nil { return err }
        if status.Err() != nil { return status.Err() }
    }

    return nil
}
