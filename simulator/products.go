package nettikauppasimulaattori


// NOTE: First integer of product-id acts as product category.
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

// TODO: Move products and categories to separate file or database.

var Products = []Product{   
    {1001, "Pirkka olut 24-pak.", 10, 25, 0.24},
    {1002, "Pirkka olut 6-pak.", 3, 8, 0.24},
    {1003, "Viskipullo", 20, 60, 0.24},
    {1004, "Leijona-pullo", 10, 40, 0.24},

    {2001, "Raspberry Pi 4 4GB", 40, 80, 0.24},
    {2002, "Raspberry Pi 4 8GB", 50, 100, 0.24},
    {2003, "VHS-kasetteja 10-pak", 5, 8, 0.24},
    {2004, "C-kasetteja 10-pak", 3, 8, 0.24},
    {2005, "LCD-televisio", 150, 300, 0.24},
    {2006, "Iso LCD-televisio", 200, 400, 0.24},
    {2007, "Tietokone", 300, 420, 0.24},
    {2008, "Parempi tietokone", 400, 700, 0.24},

    {3001, "Ruisleipä", 1, 3, 0.14},
    {3002, "Rasvaton Maito 1L", 1, 2, 0.14},
    {3003, "Kevytmaito 1L", 1, 2, 0.14},
    {3004, "Täysmaito 1L", 1, 2, 0.14},
    {3005, "Kauraleipä", 1, 3, 0.14},
    {3006, "Mysliä 1kg", 1, 4, 0.14},
    {3007, "Perunoita 1kg", 1, 2, 0.14},
    {3008, "Sipsejä 1kg", 3, 10, 0.14},
    {3009, "Irtokarkkeja 1kg", 2, 8, 0.14},

    {4000, "Silmarillion, J.R.R Tolkien", 10, 25, 0.10},
    {4001, "Tabu, Timo Mukka", 5, 15, 0.10},
    {4002, "Robottien aamunkoitto, Isaac Asimov", 10, 15, 0.10},
    {4003, "Holmenkollen, Matti Hagelberg", 15, 30, 0.10},
    {4004, "2001, Arthur Clarke", 10, 30, 0.10},
    {4005, "Foucaltin heiluri, Umberto Eco", 20, 45, 0.10},
}   


func FilterProductsByCategory(products []Product, categories []int) []Product {
    // var filteredProducts []Product
    filteredProducts := []Product{}

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

func FilterProductsByPrice(products []Product, max_price int) []Product {
    // var filteredProducts []Product
    filteredProducts := []Product{}

    for _, product := range products {
        if product.price <= max_price {
            filteredProducts = append(filteredProducts, product)
        }
    }

    return filteredProducts
}
