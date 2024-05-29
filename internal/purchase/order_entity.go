package purchase

import "fmt"

const (
	DeliveryVelocity int = 40 // In kph
)

type UserLocation struct {
	Lat  float64 `json:"lat" binding:"required"`
	Long float64 `json:"long" binding:"required"`
}

type Item struct {
	ItemID   string `json:"itemId" binding:"required"`
	Quantity int    `json:"quantity" binding:"required"`
}

type Order struct {
	MerchantID      string `json:"merchantId" binding:"required"`
	IsStartingPoint bool   `json:"isStartingPoint"`
	Items           []Item `json:"items" binding:"required"`
}

type Request struct {
	UserLocation UserLocation `json:"userLocation" binding:"required"`
	Orders       []Order      `json:"orders" binding:"required"`
}

func (r Request) ValidateRequest() error {
	// Validate that there is exactly one order with isStartingPoint == true
	startingPointCount := 0
	for _, order := range r.Orders {
		if order.IsStartingPoint {
			startingPointCount++
		}
	}
	if startingPointCount != 1 {
		return fmt.Errorf("exactly one order must have isStartingPoint == true")
	}

	// Additional validation can be added here if needed

	return nil
}
