package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type Product struct {
	ID    string `json:"id"`
	Name  string
	Price float64
}

func (p Product) IncreasePrice(amount float64) *Product {
	if amount < 0 {
		return nil
	}
	return &p
}

func main() {
	product := Product{
		ID:    "1",
		Name:  "Product 1",
		Price: 10.0,
	}
	jsonProduct, err := json.Marshal(product)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(jsonProduct))

	var product2 Product
	json.Unmarshal(jsonProduct, &product2)

	fmt.Println(product2)
}
