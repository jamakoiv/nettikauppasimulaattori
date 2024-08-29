package nettikauppasimulaattori

import (
	"fmt"
	"math"
	"testing"
	"time"
)

func CreateTestCustomers() []Customer {
	var customerA, customerB Customer

	customerA.Id = 13
	customerA.First_name = "Tero"
	customerA.Last_name = "Teronen"
	customerA.Most_active = 9
	customerA.Max_budget = 25
	customerA.Base_purchase_probability = 0.2
	customerA.Product_categories = []int{3}

	customerB.Id = 14
	customerB.First_name = "Liisa"
	customerB.Last_name = "Peronen"
	customerB.Most_active = 9
	customerB.Max_budget = 150
	customerB.Base_purchase_probability = 0.25
	customerB.Product_categories = []int{4}

	var res []Customer
	res = append(res, customerA)
	res = append(res, customerB)

	return res
}

func CheckCustomerEqual(A, B Customer) bool {
	if A.Id != B.Id {
		return false
	} else if A.First_name != B.First_name {
		return false
	} else if A.Last_name != B.Last_name {
		return false
	} else if A.Most_active != B.Most_active {
		return false
	} else if A.Max_budget != B.Max_budget {
		return false
	} else if A.Base_purchase_probability != B.Base_purchase_probability {
		return false
		// } else if !reflect.DeepEqual(A.Product_categories, B.Product_categories) {
		// 	return false
	} else {
		return true
	}
}

func TestDefault_ShoppingWeekVariation(t *testing.T) {
	test_date := time.Date(2024, time.January, 30, 0, 0, 0, 0, time.UTC)
	target := 0.004
	allowed_error := 0.00001

	res := Default_ShoppingWeekVariation(test_date)

	if math.Abs(res-target) > allowed_error {
		t.Fatalf("Wanted %v, got %v, which is not within allowed error of %v",
			res, target, allowed_error)
	}
}

func TestImportCustomers(t *testing.T) {
	var test_customers_file string = "data/test_customers.parquet"

	target := CreateTestCustomers()

	res, err := ImportCustomers(test_customers_file)
	if err != nil {
		t.Fatalf(fmt.Sprintf("Received error %v", err))
	}

	if !CheckCustomerEqual(target[0], res[0]) {
		t.Fatalf(fmt.Sprintf("Received %v, expected %v", target[0], res[0]))
	}
}
