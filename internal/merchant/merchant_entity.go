package merchant

import "time"

type MerchantCategories string

const (
	SmallRestaurant			MerchantCategories = "SmallRestaurant"
	MediumRestaurant		MerchantCategories = "MediumRestaurant"
	LargeRestaurant			MerchantCategories = "LargeRestaurant"
	MerchandiseRestaurant	MerchantCategories = "MerchandiseRestaurant"
	BoothKiosk				MerchantCategories = "BoothKiosk"
	ConvenienceStore		MerchantCategories = "ConvenienceStore"
)

type Merchant struct {
	ID                  string	`json:"id" db:"id"`
	Name            	string	`json:"name" db:"name"`
	MerchantCategory 	MerchantCategories	`json:"merchantCategory" db:"merchant_category"`
	ImageUrl 			string	`json:"imageUrl" db:"image_url"`
	LocationLat 		float	`json:"locationLat" db:"location_lat"`
	LocationLong 		float	`json:"locationLong" db:"location_long"`
	CreatedAt           time.Time 	`json:"createdAt" db:"created_at"`	
}

type Location struct {
	Lat float `json:"lat" binding:"required"`
	Long float `json:"long" binding:"required"`
}

type CreateMerchantDTO struct {
	Name            	string	`json:"name" binding:"required,min=2,max=30"`
	MerchantCategory 	MerchantCategories	`json:"merchantCategory" binding:"required"`
	ImageUrl 			string	`json:"imageUrl" binding:"required,url"`
	Location			Location `json:"location" binding:"required"`
}

type CreateMerchantResponse struct {
	MerchantID string `json:"merchantId"`
}

type ProductCategories string

const (
	Beverage	ProductCategories = "Beverage"
	Food		ProductCategories = "Food"
	Snack		ProductCategories = "Snack"
	Condiments	ProductCategories = "Condiments"
	Additions	ProductCategories = "Additions"
)

type Item struct {
	ID                  string	`json:"id" db:"id"`
	MerchantID          string	`json:"merchantId" db:"merchant_id"`
	Name            	string	`json:"name" db:"name"`
	ProductCategory 	ProductCategories	`json:"productCategory" db:"product_category"`
	Price 				integer	`json:"price" db:"price"`
	ImageUrl 			string	`json:"imageUrl" db:"image_url"`
	CreatedAt           time.Time 	`json:"createdAt" db:"created_at"`	
}

type CreateItemDTO struct {
	Name            	string	`json:"name" binding:"required,min=2,max=30"`
	ProductCategory 	ProductCategories	`json:"productCategory" binding:"required"`
	Price				integer	`json:"price" binding:"required,min=1"`
	ImageUrl 			string	`json:"imageUrl" binding:"required,url"`
}

type CreateItemResponse struct {
	ItemID string `json:"itemId"`
}
