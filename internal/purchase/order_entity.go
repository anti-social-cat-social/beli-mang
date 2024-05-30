package purchase

import "fmt"

const (
	DeliveryVelocity int = 40 // In kph
)

type OrderEstimation struct {
	ID            string  `json:"calculatedEstimateId" db:"id"`
	UserID        string  `json:"-" db:"user_id"`
	UserLat       float64 `json:"-" db:"user_location_lat"`
	UserLong      float64 `json:"-" db:"user_location_long"`
	Price         int     `json:"totalPrice" db:"total_price"`
	EstimatedTime int     `json:"estimatedDeliveryTimeInMinutes" db:"estimated_delivery_time"`
}

type OrderEstimationDetail struct {
	OrderEstimationID string `db:"order_estimation_id"`
	ItemID            string `db:"item_id"`
	Quantity          int    `db:"quantity"`
}

type UserLocation struct {
	Lat  float64 `json:"lat" binding:"required"`
	Long float64 `json:"long" binding:"required"`
}

type Item struct {
	ItemID   string `json:"itemId" binding:"required"`
	Quantity int    `json:"quantity" binding:"required"`
}

// Order estimation
type Order struct {
	MerchantID      string `json:"merchantId" binding:"required"`
	IsStartingPoint bool   `json:"isStartingPoint"`
	Items           []Item `json:"items" binding:"required"`
}

// Order that has been placed / confirmed
type ActualOrder struct {
	OrderId           string `json:"orderId"`
	OrderEstimationId string `json:"calculatedEstimateId,omitempty" binding:"required"`
}

type Request struct {
	UserId       string
	UserLocation UserLocation `json:"userLocation" binding:"required"`
	Orders       []Order      `json:"orders" binding:"required"`
}

type OrderEstimationResponse struct {
	TotalPrice                     int    `json:"totalPrice"`
	EstimatedDeliveryTimeInMinutes int    `json:"estimatedDeliveryTimeInMinutes"`
	CalculatedEstimateID           string `json:"calculatedEstimateId"`
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
