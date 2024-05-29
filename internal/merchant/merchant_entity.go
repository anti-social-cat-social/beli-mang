package merchant

import "time"

type MerchantCategories string

const (
	SmallRestaurant       MerchantCategories = "SmallRestaurant"
	MediumRestaurant      MerchantCategories = "MediumRestaurant"
	LargeRestaurant       MerchantCategories = "LargeRestaurant"
	MerchandiseRestaurant MerchantCategories = "MerchandiseRestaurant"
	BoothKiosk            MerchantCategories = "BoothKiosk"
	ConvenienceStore      MerchantCategories = "ConvenienceStore"
)

type Merchant struct {
	ID               string             `json:"id" db:"id"`
	Name             string             `json:"name" db:"name"`
	MerchantCategory MerchantCategories `json:"merchantCategory" db:"merchant_category"`
	ImageUrl         string             `json:"imageUrl" db:"image_url"`
	LocationLat      float32            `json:"locationLat" db:"location_lat"`
	LocationLong     float32            `json:"locationLong" db:"location_long"`
	CreatedAt        time.Time          `json:"createdAt" db:"created_at"`
}

type Location struct {
	Lat  float32 `json:"lat" binding:"required"`
	Long float32 `json:"long" binding:"required"`
}

type CreateMerchantDTO struct {
	Name             string             `json:"name" binding:"required,min=2,max=30"`
	MerchantCategory MerchantCategories `json:"merchantCategory" binding:"required"`
	ImageUrl         string             `json:"imageUrl" binding:"required,url"`
	Location         Location           `json:"location" binding:"required"`
}

type CreateMerchantResponse struct {
	MerchantID string `json:"merchantId"`
}

type ProductCategories string

const (
	Beverage   ProductCategories = "Beverage"
	Food       ProductCategories = "Food"
	Snack      ProductCategories = "Snack"
	Condiments ProductCategories = "Condiments"
	Additions  ProductCategories = "Additions"
)

type Item struct {
	ID              string            `json:"id" db:"id"`
	MerchantID      string            `json:"merchantId" db:"merchant_id"`
	Name            string            `json:"name" db:"name"`
	ProductCategory ProductCategories `json:"productCategory" db:"product_category"`
	Price           int               `json:"price" db:"price"`
	ImageUrl        string            `json:"imageUrl" db:"image_url"`
	CreatedAt       time.Time         `json:"createdAt" db:"created_at"`
}

type CreateItemDTO struct {
	Name            string            `json:"name" binding:"required,min=2,max=30"`
	ProductCategory ProductCategories `json:"productCategory" binding:"required"`
	Price           int               `json:"price" binding:"required,min=1"`
	ImageUrl        string            `json:"imageUrl" binding:"required,url"`
}

type GetItemQueryParam struct {
	ItemID          string            `form:"itemId"`
	Limit           int               `form:"limit"`
	Offset          int               `form:"offset"`
	Name            string            `form:"name"`
	ProductCategory ProductCategories `form:"productCategory"`
	CreatedAt       Sort              `form:"createdAt"`
}

type ItemResponse struct {
	ID              string            `json:"id"`
	Name            string            `json:"name"`
	ProductCategory ProductCategories `json:"productCategory"`
	Price           int               `json:"price"`
	ImageUrl        string            `json:"imageUrl"`
	CreatedAt       string            `json:"createdAt"`
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
	MerchantID       string             `form:"merchantId"`
	Limit            int                `form:"limit"`
	Offset           int                `form:"offset"`
	Name             string             `form:"name"`
	MerchantCategory MerchantCategories `form:"merchantCategory"`
	CreatedAt        Sort               `form:"createdAt"`
}

type GetMerchantResponse struct {
	MerchantID       string             `json:"merchantId"`
	Name             string             `json:"name"`
	MerchantCategory MerchantCategories `json:"merchantCategory"`
	ImageUrl         string             `json:"imageUrl"`
	Location         Location           `json:"location"`
	CreatedAt        string             `json:"createdAt"`
}

func FormatGetMerchantResponse(merchants []Merchant) []GetMerchantResponse {
	getMerchantResponse := []GetMerchantResponse{}

	for _, merchant := range merchants {
		row := GetMerchantResponse{
			MerchantID:       merchant.ID,
			Name:             merchant.Name,
			MerchantCategory: merchant.MerchantCategory,
			ImageUrl:         merchant.ImageUrl,
			Location: Location{
				Lat:  merchant.LocationLat,
				Long: merchant.LocationLong,
			},
			CreatedAt: merchant.CreatedAt.Format(time.RFC3339),
		}
		getMerchantResponse = append(getMerchantResponse, row)
	}

	return getMerchantResponse
}

func FormatItemResponse(items []Item) []ItemResponse {
	itemsResponse := []ItemResponse{}

	for _, item := range items {
		row := ItemResponse{
			ID:              item.ID,
			Name:            item.Name,
			ProductCategory: item.ProductCategory,
			Price:           item.Price,
			ImageUrl:        item.ImageUrl,
			CreatedAt:       item.CreatedAt.Format(time.RFC3339),
		}

		itemsResponse = append(itemsResponse, row)
	}

	return itemsResponse
}
