package service

type Item struct {
	ID    string  `json:"id"`
	Price float64 `json:"price"`
	Qty   int     `json:"qty"`
}

type Cart struct {
	DiscountType  string  `json:"discount_type"`
	DiscountValue float64 `json:"discount_value"`
	Items         []Item  `json:"items"`
}

func (c Cart) Subtotal() float64 {
	var sum float64
	for _, it := range c.Items {
		sum += it.Price * float64(it.Qty)
	}
	return sum
}


type DiscountStrategy interface {
	Calculate(cart Cart) float64
}

func NewStrategy(t string, val float64) DiscountStrategy {
	switch t {
	case "percent":
		return &PercentDiscount{Value: val}
	case "bulk":
		return &BulkDiscount{MinQty: 5, Value: val}
	default:
		return &NoDiscount{}
	}
}

type NoDiscount struct{}

func (d *NoDiscount) Calculate(cart Cart) float64 {
	return 0
}

type PercentDiscount struct {
	Value float64
}

func (d *PercentDiscount) Calculate(cart Cart) float64 {
	return cart.Subtotal() * d.Value / 100
}

type BulkDiscount struct {
	MinQty int
	Value  float64
}

func (d *BulkDiscount) Calculate(cart Cart) float64 {
	totalQty := 0
	for _, it := range cart.Items {
		totalQty += it.Qty
	}
	if totalQty > d.MinQty {
		return cart.Subtotal() * d.Value / 100
	}
	return 0
}
