package main

import (
	"log"
	"net/http"

	"go-cart-discount/api"
)

func main() {
	http.HandleFunc("/price", api.Price)
	http.HandleFunc("/price/bulk", api.PriceBulk)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
