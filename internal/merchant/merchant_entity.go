package merchant

import (
	"time"
	// "log"	
)

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
	LocationLat 		float32	`json:"locationLat" db:"location_lat"`
	LocationLong 		float32	`json:"locationLong" db:"location_long"`
	CreatedAt           time.Time 	`json:"createdAt" db:"created_at"`	
}

type Location struct {
	Lat float32 `json:"lat" binding:"required"`
	Long float32 `json:"long" binding:"required"`
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
	Price 				int	`json:"price" db:"price"`
	ImageUrl 			string	`json:"imageUrl" db:"image_url"`
	CreatedAt           time.Time 	`json:"createdAt" db:"created_at"`	
}

type CreateItemDTO struct {
	Name            	string	`json:"name" binding:"required,min=2,max=30"`
	ProductCategory 	ProductCategories	`json:"productCategory" binding:"required"`
	Price				int	`json:"price" binding:"required,min=1"`
	ImageUrl 			string	`json:"imageUrl" binding:"required,url"`
}

type CreateItemResponse struct {
	ItemID string `json:"itemId"`
}

type Sort string

const (
	Asc  Sort = "asc"
	Desc Sort = "desc"
)

type GetMerchantQueryParams struct {
	MerchantID			string 	`form:"merchantId"`
	Limit       		int    	`form:"limit"`
	Offset      		int    	`form:"offset"`
	Name        		string 	`form:"name"`
	MerchantCategory	MerchantCategories	`form:"merchantCategory"`
	CreatedAt 			Sort	`form:"createdAt"`
}

type GetMerchantResponse struct {
	MerchantID			string 	`json:"merchantId"`
	Name            	string	`json:"name"`
	MerchantCategory 	MerchantCategories	`json:"merchantCategory"`
	ImageUrl 			string	`json:"imageUrl"`
	Location			Location `json:"location"`
	CreatedAt           string 	`json:"createdAt"`
}

func FormatGetMerchantResponse(merchants []Merchant) []GetMerchantResponse {
	getMerchantResponse := []GetMerchantResponse{}

	for _, merchant := range merchants {
		row := GetMerchantResponse{
			MerchantID: merchant.ID,
			Name: merchant.Name,
			MerchantCategory: merchant.MerchantCategory,
			ImageUrl: merchant.ImageUrl,
			Location: Location{
				Lat: merchant.LocationLat,
				Long: merchant.LocationLong,
			},
			CreatedAt: merchant.CreatedAt.Format(time.RFC3339),
		}
		getMerchantResponse = append(getMerchantResponse, row)
	}

	return getMerchantResponse
}

type ItemForNearbyMerchant struct {
	ID                  string	`json:"itemId"`
	Name            	string	`json:"name"`
	ProductCategory 	ProductCategories	`json:"productCategory"`
	Price 				int	`json:"price"`
	ImageUrl 			string	`json:"imageUrl"`
	CreatedAt           time.Time 	`json:"createdAt"`	
}

type NearbyMerchantWithItem struct {
	ID                  string	`json:"merchantId"`
	Name            	string	`json:"name"`
	MerchantCategory 	MerchantCategories	`json:"merchantCategory"`
	ImageUrl 			string	`json:"imageUrl"`
	Location			Location `json:"location"`
	CreatedAt           time.Time 	`json:"createdAt"`	
	// Items				[]ItemForNearbyMerchant `json:"items"`
}

type MerchantWithItemQueryResult struct {
	MerchantId	string `json:"merchantId" db:"merchant_id"`
	MerchantName	string `json:"merchantName" db:"merchant_name"`
	MerchantCategory	MerchantCategories `json:"merchantCategory" db:"merchant_category"`
	MerchantImageUrl	string `json:"merchantImageUrl" db:"merchant_image_url"`
	LocationLat	float32 `json:"locationLat" db:"location_lat"`
	LocationLong	float32 `json:"locationLong" db:"location_long"`
	MerchantCreatedAt	time.Time `json:"merchantCreatedAt" db:"merchant_created_at"`
	ItemId	string `json:"itemId" db:"item_id"`
	ItemName	string `json:"itemName" db:"item_name"`
	ProductCategory	ProductCategories `json:"productCategory" db:"product_category"`
	Price	int `json:"price" db:"price"`
	ItemImageUrl	string `json:"itemImageUrl" db:"item_image_url"`
	ItemCreatedAt	time.Time `json:"itemCreatedAt" db:"item_created_at"`
	Distance			float32	`json:"distance" db:"distance"`
}

type NearbyMerchantWithItemResponse struct {
	Merchant	NearbyMerchantWithItem `json:"merchant"`
	Items		[]ItemForNearbyMerchant `json:"items"`
}

func FormatNearbyMerchantWithItemResponse(merchants []MerchantWithItemQueryResult) []NearbyMerchantWithItemResponse {
	res := []NearbyMerchantWithItemResponse{}
	merchantAndItems := NearbyMerchantWithItemResponse{}
	merchant := NearbyMerchantWithItem{}
	item := ItemForNearbyMerchant{}
	items := []ItemForNearbyMerchant{}
	loc := Location{}
	totalLen := len(merchants)

	for ix, m := range merchants {
		loc = Location{
			Lat: m.LocationLat,
			Long: m.LocationLong,
		}
		merchant = NearbyMerchantWithItem{
			Location: loc,
			ID: m.MerchantId,
			Name: m.MerchantName,
			MerchantCategory: m.MerchantCategory,
			ImageUrl: m.MerchantImageUrl,
			CreatedAt: m.MerchantCreatedAt,       
		}
		item = ItemForNearbyMerchant{
			ID: m.ItemId,           
			Name: m.ItemName,         
			ProductCategory: m.ProductCategory,
			Price: m.Price,
			ImageUrl: m.ItemImageUrl,
			CreatedAt: m.ItemCreatedAt,
		}
		items = append(items, item)

		if ix+1 == totalLen {
			merchantAndItems = NearbyMerchantWithItemResponse{
				Merchant: merchant,
				Items: items,
			} 
			res = append(res, merchantAndItems)
		} else {
			if m.MerchantId != merchants[ix+1].MerchantId {
				merchantAndItems = NearbyMerchantWithItemResponse{
					Merchant: merchant,
					Items: items,
				} 
				res = append(res, merchantAndItems)
				items = []ItemForNearbyMerchant{}
			}
		}
	}

	return res
}