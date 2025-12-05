package api

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"

	"go-cart-discount/service"
)

type PriceResponse struct {
	Subtotal float64 `json:"subtotal"`
	Discount float64 `json:"discount"`
	Total    float64 `json:"total"`
	Error    string  `json:"error,omitempty"`
}

func Price(w http.ResponseWriter, r *http.Request) {
	var cart service.Cart
	json.NewDecoder(r.Body).Decode(&cart)

	sub := cart.Subtotal()

	strategy := service.NewStrategy(cart.DiscountType, cart.DiscountValue)
	disc := strategy.Calculate(cart)

	resp := PriceResponse{
		Subtotal: sub,
		Discount: disc,
		Total:    sub - disc,
	}

	json.NewEncoder(w).Encode(resp)
}

func PriceBulk(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	type cartData struct {
		service.Cart
	}

	var carts []cartData
	json.NewDecoder(r.Body).Decode(&carts)

	maxGoroutines := 3
	sem := make(chan struct{}, maxGoroutines)

	wg := sync.WaitGroup{}
	responses := make([]PriceResponse, len(carts))

	for i := range carts {
		wg.Add(1)

		sem <- struct{}{} // semaphore

		go func(i int) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				responses[i].Error = "canceled"
				<-sem
				return
			default:
				sub := carts[i].Subtotal()
				strategy := service.NewStrategy(carts[i].DiscountType, carts[i].DiscountValue)
				disc := strategy.Calculate(carts[i].Cart)

				responses[i] = PriceResponse{
					Subtotal: sub,
					Discount: disc,
					Total:    sub - disc,
				}
				<-sem
			}
		}(i)
	}

	wg.Wait()
	json.NewEncoder(w).Encode(responses)
}
