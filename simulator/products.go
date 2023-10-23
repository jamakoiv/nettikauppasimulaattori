package nettikauppasimulaattori

// TODO: Not actually necessary to have price and other info on this side.

type Product struct {
    id int
    name string
    wholesale_price int
    price int
    vat float64
}

const (
    ALCOHOL = 1
    ELECTRONICS = 2
    GROCERIES = 3
    BOOKS = 4
    CLOTHING = 5
)

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


func filterProducts(products []Product, categories []int) []Product {
    var filteredProducts []Product

    for _, category := range categories {
        for _, product := range products {
            // Check if the first digit of the product ID matches the specified category
            if product.id/1000 == category {
                    filteredProducts = append(filteredProducts, product)
            }
        }
    }

    return filteredProducts
}
