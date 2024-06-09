package purchase

import (
	merchantModule "belimang/internal/merchant"
	"fmt"
	"time"
)

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
	ItemID   string `json:"itemId" binding:"required,uuid"`
	Quantity int    `json:"quantity" binding:"required"`
}

// Order estimation
type Order struct {
	MerchantID      string `json:"merchantId" binding:"required,uuid"`
	IsStartingPoint bool   `json:"isStartingPoint"`
	Items           []Item `json:"items" binding:"required,dive"`
}

// Order that has been placed / confirmed
type ActualOrder struct {
	OrderId           string `json:"orderId"`
	OrderEstimationId string `json:"calculatedEstimateId,omitempty" binding:"required"`
}

type Request struct {
	UserId       string
	UserLocation UserLocation `json:"userLocation" binding:"required"`
	Orders       []Order      `json:"orders" binding:"required,dive"`
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

type GetOrderHistQueryParams struct {
	MerchantID       string             `form:"merchantId"`
	Limit            int                `form:"limit"`
	Offset           int                `form:"offset"`
	Name             string             `form:"name"`
	MerchantCategory merchantModule.MerchantCategories `form:"merchantCategory"`
}

type GetOrderHistQueryResult struct {
	OrderId       	 string             `json:"orderId" db:"order_id"`
	MerchantId        string             `json:"merchantId" db:"merchant_id"`
	MerchantName      string             `json:"merchantName" db:"merchant_name"`
	MerchantCategory  merchantModule.MerchantCategories `json:"merchantCategory" db:"merchant_category"`
	MerchantImageUrl  string             `json:"merchantImageUrl" db:"merchant_image_url"`
	LocationLat       float64            `json:"locationLat" db:"location_lat"`
	LocationLong      float64            `json:"locationLong" db:"location_long"`
	MerchantCreatedAt time.Time          `json:"merchantCreatedAt" db:"merchant_created_at"`
	ItemId            string             `json:"itemId" db:"item_id"`
	ItemName          string             `json:"itemName" db:"item_name"`
	ProductCategory   merchantModule.ProductCategories  `json:"productCategory" db:"product_category"`
	Price             int                `json:"price" db:"price"`
	ItemImageUrl      string             `json:"itemImageUrl" db:"item_image_url"`
	ItemCreatedAt     time.Time          `json:"itemCreatedAt" db:"item_created_at"`
	Quantity          int            	 `json:"quantity" db:"quantity"`
}

type OrderHistMerchant struct {
	ID               string             `json:"merchantId"`
	Name             string             `json:"name"`
	MerchantCategory merchantModule.MerchantCategories `json:"merchantCategory"`
	ImageUrl         string             `json:"imageUrl"`
	Location         merchantModule.Location           `json:"location"`
	CreatedAt        time.Time          `json:"createdAt"`
}

type OrderHistItem struct {
	ID              string            `json:"itemId"`
	Name            string            `json:"name"`
	ProductCategory merchantModule.ProductCategories `json:"productCategory"`
	Price           int               `json:"price"`
	Quantity        int            	  `json:"quantity"`
	ImageUrl        string            `json:"imageUrl"`
	CreatedAt       time.Time         `json:"createdAt"`
}

type GetOrderHistResponse struct {
	OrderId	 string				`json:"orderId"`
	Merchant OrderHistMerchant  `json:"merchant"`
	Items    []OrderHistItem	`json:"items"`
}

type GetOrderHistResponseOnly struct {
	Merchant OrderHistMerchant  `json:"merchant"`
	Items    []OrderHistItem	`json:"items"`
}

type GetOrderHistResponseWithOrderId struct {
	OrderId	 	string						`json:"orderId"`
	Orders		[]GetOrderHistResponseOnly	`json:"orders"`
}

func FormatGetOrderHistResponseWithOrderId(data []GetOrderHistResponse) []GetOrderHistResponseWithOrderId {
	ordersWithId := []GetOrderHistResponseWithOrderId{}
	orderWithId := GetOrderHistResponseWithOrderId{}
	orders := []GetOrderHistResponseOnly{}
	order := GetOrderHistResponseOnly{}
	totalLen := len(data)

	for ix, o := range data {
		order = GetOrderHistResponseOnly{
			Merchant: o.Merchant,
			Items: o.Items,
		}
		orders = append(orders, order)

		if ix+1 == totalLen {
			orderWithId = GetOrderHistResponseWithOrderId{
				OrderId: o.OrderId,
				Orders: orders,
			}
			ordersWithId = append(ordersWithId, orderWithId)
		} else {
			if o.OrderId != data[ix+1].OrderId {
				orderWithId = GetOrderHistResponseWithOrderId{
					OrderId: o.OrderId,
					Orders: orders,
				}
				ordersWithId = append(ordersWithId, orderWithId)
				orders = []GetOrderHistResponseOnly{}
			}
		}
	}

	return ordersWithId
}

func FormatGetOrderHistResponse(orders []GetOrderHistQueryResult) []GetOrderHistResponse {
	res := []GetOrderHistResponse{}
	merchantAndItems := GetOrderHistResponse{}
	merchant := OrderHistMerchant{}
	item := OrderHistItem{}
	items := []OrderHistItem{}
	loc := merchantModule.Location{}
	totalLen := len(orders)

	for ix, m := range orders {
		loc = merchantModule.Location{
			Lat:  m.LocationLat,
			Long: m.LocationLong,
		}
		merchant = OrderHistMerchant{
			Location:         loc,
			ID:               m.MerchantId,
			Name:             m.MerchantName,
			MerchantCategory: m.MerchantCategory,
			ImageUrl:         m.MerchantImageUrl,
			CreatedAt:        m.MerchantCreatedAt,
		}
		item = OrderHistItem{
			ID:              m.ItemId,
			Name:            m.ItemName,
			ProductCategory: m.ProductCategory,
			Price:           m.Price,
			ImageUrl:        m.ItemImageUrl,
			CreatedAt:       m.ItemCreatedAt,
			Quantity:		 m.Quantity,
		}
		items = append(items, item)

		if ix+1 == totalLen {
			merchantAndItems = GetOrderHistResponse{
				OrderId: m.OrderId,
				Merchant: merchant,
				Items:    items,
			}
			res = append(res, merchantAndItems)
		} else {
			if m.OrderId+m.MerchantId != orders[ix+1].OrderId+orders[ix+1].MerchantId {
				merchantAndItems = GetOrderHistResponse{
					OrderId: m.OrderId,
					Merchant: merchant,
					Items:    items,
				}
				res = append(res, merchantAndItems)
				items = []OrderHistItem{}
			}
		}
	}

	return res
}