package merchant

import (
	// "errors"
	localError "belimang/pkg/error"
	// "strconv"
	// "time"
	"github.com/google/uuid"
	// "log"
)

type IMerchantUsecase interface {
	CreateMerchant(req CreateMerchantDTO) (*CreateMerchantResponse, *localError.GlobalError)
	CreateItem(merchantId string, req CreateItemDTO) (*CreateItemResponse, *localError.GlobalError)
	FindAllMerchants(query GetMerchantQueryParams) ([]GetMerchantResponse, *localError.GlobalError)
	FindMerchantById(id string) (*Merchant, *localError.GlobalError)
	FindAllItem(query GetItemQueryParam, merchatId string) ([]ItemResponse, *localError.GlobalError)
}

type merchantUsecase struct {
	repo IMerchantRepository
}

func NewMerchantUsecase(repo IMerchantRepository) IMerchantUsecase {
	return &merchantUsecase{
		repo: repo,
	}
}

func (uc *merchantUsecase) CreateMerchant(req CreateMerchantDTO) (*CreateMerchantResponse, *localError.GlobalError) {
	merchant := Merchant{
		ID:               uuid.NewString(),
		Name:             req.Name,
		MerchantCategory: req.MerchantCategory,
		ImageUrl:         req.ImageUrl,
		LocationLat:      req.Location.Lat,
		LocationLong:     req.Location.Long,
	}

	err := uc.repo.CreateMerchant(merchant)
	if err != nil {
		return nil, err
	}

	response := CreateMerchantResponse{
		MerchantID: merchant.ID,
	}

	return &response, nil
}

func (uc *merchantUsecase) CreateItem(merchantId string, req CreateItemDTO) (*CreateItemResponse, *localError.GlobalError) {
	_, err := uc.repo.FindMerchantById(merchantId)
	if err != nil {
		return nil, localError.ErrNotFound("merchant not found", err.Error)
	}

	item := Item{
		ID:              uuid.NewString(),
		MerchantID:      merchantId,
		Name:            req.Name,
		ProductCategory: req.ProductCategory,
		Price:           req.Price,
		ImageUrl:        req.ImageUrl,
	}

	err = uc.repo.CreateItem(item)
	if err != nil {
		return nil, err
	}

	response := CreateItemResponse{
		ItemID: item.ID,
	}

	return &response, nil
}

func (uc *merchantUsecase) FindAllMerchants(query GetMerchantQueryParams) ([]GetMerchantResponse, *localError.GlobalError) {
	merchants, err := uc.repo.FindAllMerchants(query)
	if err != nil {
		return nil, err
	}

	resp := FormatGetMerchantResponse(merchants)

	return resp, nil
}

func (uc *merchantUsecase) FindMerchantById(id string) (*Merchant, *localError.GlobalError) {
	return uc.repo.FindMerchantById(id)
}

func (uc *merchantUsecase) FindAllItem(query GetItemQueryParam, merchantId string) ([]ItemResponse, *localError.GlobalError) {
	// Check if the merchant is exists
	merchant, err := uc.repo.FindMerchantById(merchantId)
	if merchant == nil {
		return nil, err
	}

	items, err := uc.repo.FindAllItem(query, merchantId)

	if err != nil {
		return nil, err
	}

	response := FormatItemResponse(items)

	return response, nil
}
