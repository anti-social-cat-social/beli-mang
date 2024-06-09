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
	FindAllMerchants(query GetMerchantQueryParams) (GetMerchantResponseAndMeta, *localError.GlobalError)
	FindMerchantById(id string) (*Merchant, *localError.GlobalError)
	FindAllItem(query GetItemQueryParam, merchatId string) (ItemResponseAndMeta, *localError.GlobalError)
	CheckMerchantIDs(IDs []string) ([]Merchant, *localError.GlobalError)
	CheckItemIDs(IDs []string) ([]Item, *localError.GlobalError)
	FindNearbyMerchants(location Location, query GetMerchantQueryParams) (NearbyMerchantWithItemResponseAndMeta, *localError.GlobalError)
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

type Meta struct {
	Limit int `json:"limit"`
	Offset int `json:"offset"`
	Total int `json:"total"`
}

type NearbyMerchantWithItemResponseAndMeta struct {
	Data []NearbyMerchantWithItemResponse `json:"data"`
	Meta Meta `json:"meta"`
}

type GetMerchantResponseAndMeta struct {
	Data []GetMerchantResponse `json:"data"`
	Meta Meta `json:"meta"`
}

type ItemResponseAndMeta struct {
	Data []ItemResponse `json:"data"`
	Meta Meta `json:"meta"`
}

func (uc *merchantUsecase) FindAllMerchants(query GetMerchantQueryParams) (GetMerchantResponseAndMeta, *localError.GlobalError) {
	merchants, err := uc.repo.FindAllMerchants(query)
	if err != nil {
		return GetMerchantResponseAndMeta{}, err
	}

	resp := FormatGetMerchantResponse(merchants)

	limit := 5
	offset := 0
	if query.Limit != 0 {
		limit = query.Limit
	}
	if query.Offset != 0 {
		offset = query.Offset
	}

	meta := Meta{
		Limit: limit,
		Offset: offset,
		Total: len(resp),
	}
	
	return GetMerchantResponseAndMeta{
		Data: resp,
		Meta: meta,
	}, nil
}

func (uc *merchantUsecase) FindMerchantById(id string) (*Merchant, *localError.GlobalError) {
	return uc.repo.FindMerchantById(id)
}

func (uc *merchantUsecase) FindAllItem(query GetItemQueryParam, merchantId string) (ItemResponseAndMeta, *localError.GlobalError) {
	// Check if the merchant is exists
	merchant, err := uc.repo.FindMerchantById(merchantId)
	if merchant == nil {
		return ItemResponseAndMeta{}, err
	}

	items, err := uc.repo.FindAllItem(query, merchantId)

	if err != nil {
		return ItemResponseAndMeta{}, err
	}

	response := FormatItemResponse(items)

	limit := 5
	offset := 0
	if query.Limit != 0 {
		limit = query.Limit
	}
	if query.Offset != 0 {
		offset = query.Offset
	}

	meta := Meta{
		Limit: limit,
		Offset: offset,
		Total: len(response),
	}

	return ItemResponseAndMeta{
		Data: response,
		Meta: meta,
	}, nil
}

func (uc *merchantUsecase) FindNearbyMerchants(location Location, query GetMerchantQueryParams) (NearbyMerchantWithItemResponseAndMeta, *localError.GlobalError) {
	merchants, err := uc.repo.FindNearbyMerchants(location, query)
  	if err != nil {
		return NearbyMerchantWithItemResponseAndMeta{}, err
	}
  
  	resp := FormatNearbyMerchantWithItemResponse(merchants)

	limit := 5
	offset := 0
	if query.Limit != 0 {
		limit = query.Limit
	}
	if query.Offset != 0 {
		offset = query.Offset
	}

	if offset >= len(resp) {
		return NearbyMerchantWithItemResponseAndMeta{}, nil
	}
	if limit < 0 {
		return NearbyMerchantWithItemResponseAndMeta{}, nil
	}
	if offset+limit > len(resp) {
		cut := offset+limit-len(resp)
		limit = limit-cut
	}

	meta := Meta{
		Limit: limit,
		Offset: offset,
		Total: len(resp),
	}
	
	return NearbyMerchantWithItemResponseAndMeta{
		Data: resp[offset:offset+limit],
		Meta: meta,
	}, nil
}

func (uc *merchantUsecase) CheckMerchantIDs(IDs []string) ([]Merchant, *localError.GlobalError) {
	return uc.repo.CheckMerchantIDs(IDs)
}

func (uc *merchantUsecase) CheckItemIDs(IDs []string) ([]Item, *localError.GlobalError) {
	return uc.repo.CheckItemIDs(IDs)
}
