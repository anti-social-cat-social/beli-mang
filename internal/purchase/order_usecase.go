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
	merchantUc merchant.IMerchantUsecase
}

type IOrderUsecase interface {
	Estimate(dto Request) (float64, *localError.GlobalError)
}

func NewOrderUsecase(mUc merchant.IMerchantUsecase) IOrderUsecase {
	return &orderUsecase{
		merchantUc: mUc,
	}
}

func (uc *orderUsecase) Estimate(dto Request) (float64, *localError.GlobalError) {
	// Create maps of points
	// This maps is used to calculate area and distance
	var points []distances.Point

	userPoint := distances.Point{
		Name: "user",
		Lat:  dto.UserLocation.Lat,
		Long: dto.UserLocation.Long,
	}
	points = append(points, userPoint)

	// Generate slice of merchant point
	// Find merchant should use where In
	for _, v := range dto.Orders {
		var merchantPoint distances.Point

		merchant, mercError := uc.merchantUc.FindMerchantById(v.MerchantID)
		if mercError != nil {
			return 0.0, mercError
		}

		merchantPoint.Lat = float64(merchant.LocationLat)
		merchantPoint.Long = float64(merchant.LocationLong)
		merchantPoint.Name = merchant.Name

		points = append(points, merchantPoint)
	}

	// Throw error if the area more than 3km^2
	area := distances.CalculateArea(points)

	if area > 3*math.Pow10(6) {
		return 0.0, localError.ErrBadRequest("Area too far", fmt.Errorf("area too far"))
	}

	// Permutate shortest distance
	track, d := distances.ShortestDistance(points)
	log.Println(track, d)

	// Calculate fastest / shortest delivery time
	var time float64
	velocity := 40.0 // in kph

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

		time += twoPointDistance / velocity
	}

	return time * 60, nil
}
