package purchase

import (
	"belimang/internal/merchant"
	"belimang/pkg/distances"
	localError "belimang/pkg/error"
	"fmt"
	"log"
	"math"
)

type orderUsecase struct {
	repo       IOrderRepository
	merchantUc merchant.IMerchantUsecase
}

type IOrderUsecase interface {
	Estimate(dto Request) (*OrderEstimationResponse, *localError.GlobalError)
	PlaceOrder(entity ActualOrder) (*ActualOrder, *localError.GlobalError)
	OrderHistory(userId string, dto GetOrderHistQueryParams) ([]GetOrderHistResponseWithOrderId, *localError.GlobalError)
}

func NewOrderUsecase(repo IOrderRepository, mUc merchant.IMerchantUsecase) IOrderUsecase {
	return &orderUsecase{
		repo:       repo,
		merchantUc: mUc,
	}
}

func (uc *orderUsecase) Estimate(dto Request) (*OrderEstimationResponse, *localError.GlobalError) {
	// Create maps of points
	// This maps is used to calculate area and distance
	var (
		points              []distances.Point
		merchantIDs         []string
		itemIDs             []string
		estimationMerchants []OrderEstimationDetail
		totalPrice          int
	)

	userPoint := distances.Point{
		Name: "user",
		Lat:  dto.UserLocation.Lat,
		Long: dto.UserLocation.Long,
	}
	points = append(points, userPoint)

	// Generate slice of merchant point
	// Find merchant should use where In
	// This block is also create quantity map
	quantities := make(map[string]int)
	for _, v := range dto.Orders {
		// Append ID Merchant to get checked later
		merchantIDs = append(merchantIDs, v.MerchantID)

		// Loop to get Item IDs
		for _, item := range v.Items {
			itemIDs = append(itemIDs, item.ItemID)

			// Add quantities by item ID
			quantities[item.ItemID] = item.Quantity

			// Append order estimation merchants
			estimationMerchants = append(estimationMerchants, OrderEstimationDetail{
				ItemID:   item.ItemID,
				Quantity: item.Quantity,
			})
		}
	}

	// Get and check the ID of merchant and item
	merchants, err := uc.merchantUc.CheckMerchantIDs(merchantIDs)
	if err != nil {
		return nil, err
	}

	items, err := uc.merchantUc.CheckItemIDs(itemIDs)
	if err != nil {
		return nil, err
	}

	// Check the count
	if (len(merchants) != len(merchantIDs)) || (len(items) != len(itemIDs)) {
		return nil, localError.ErrNotFound("ID Merchant / Item not valid", fmt.Errorf("merchant or item is invalid"))
	}

	// Loop merchant
	for _, merchant := range merchants {
		var merchantPoint distances.Point

		merchantPoint.Lat = float64(merchant.LocationLat)
		merchantPoint.Long = float64(merchant.LocationLong)
		merchantPoint.Name = merchant.Name

		points = append(points, merchantPoint)
	}

	// Loop item to get total price
	for _, item := range items {
		totalPrice += quantities[item.ID] * item.Price
	}

	// Throw error if the area more than 3km^2
	area := distances.CalculateArea(points)

	if area > 3*math.Pow10(6) {
		return nil, localError.ErrBadRequest("Area too far", fmt.Errorf("area too far"))
	}

	// Permutate shortest distance
	track, d := distances.ShortestDistance(points)
	log.Println(track, d)

	// Calculate fastest / shortest delivery time
	var time float64

	for i := 0; i <= len(points)-2; i++ {
		p1 := distances.Point{
			Lat:  track[i].Lat,
			Long: track[i].Long,
		}

		p2 := distances.Point{
			Lat:  track[i+1].Lat,
			Long: track[i+1].Long,
		}

		twoPointDistance := distances.Calculate(distances.DistanceRaw{
			Start: p1,
			End:   p2,
		})

		time += twoPointDistance / float64(DeliveryVelocity)
	}

	absTime := int(math.Round(time * 60))

	// Store user estimation
	var estimation OrderEstimation = OrderEstimation{
		UserID:        dto.UserId,
		UserLat:       userPoint.Lat,
		UserLong:      userPoint.Long,
		Price:         totalPrice,
		EstimatedTime: absTime,
	}

	estimationID, err := uc.repo.CreateEstimation(&estimation)
	if err != nil {
		return nil, err
	}

	// Store order estimation merchants
	err = uc.repo.CreateOrderMerchant(estimationID, estimationMerchants)
	if err != nil {
		return nil, err
	}

	// Generate response
	response := OrderEstimationResponse{
		TotalPrice:                     totalPrice,
		EstimatedDeliveryTimeInMinutes: absTime,
		CalculatedEstimateID:           estimationID,
	}

	return &response, nil
}

func (uc *orderUsecase) PlaceOrder(entity ActualOrder) (*ActualOrder, *localError.GlobalError) {
	result, err := uc.repo.PlaceOrder(entity.OrderEstimationId)
	if err != nil {
		return nil, err
	}

	return &ActualOrder{
		OrderId: result,
	}, nil
}

func (uc *orderUsecase) OrderHistory(userId string, dto GetOrderHistQueryParams) ([]GetOrderHistResponseWithOrderId, *localError.GlobalError) {
	orders, err := uc.repo.OrderHistory(userId, dto)
  	if err != nil {
		return nil, err
	}
  
	ordersWithOrderId := FormatGetOrderHistResponse(orders)

	resp := FormatGetOrderHistResponseWithOrderId(ordersWithOrderId)

	limit := 5
	offset := 0
	if dto.Limit != 0 {
		limit = dto.Limit
	}
	if dto.Offset != 0 {
		offset = dto.Offset
	}

	if offset >= len(resp) {
		return []GetOrderHistResponseWithOrderId{}, nil
	}
	if limit < 0 {
		return []GetOrderHistResponseWithOrderId{}, nil
	}
	if offset+limit > len(resp) {
		cut := offset+limit-len(resp)
		limit = limit-cut
	}
	
	return resp[offset:offset+limit], nil
}