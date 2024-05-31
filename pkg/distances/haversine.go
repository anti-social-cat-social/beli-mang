package distances

import (
	"belimang/pkg/helper"
	"math"
)

const EARTH_RADIUS float64 = 6371.009

func Calculate(raw DistanceRaw) float64 {
	// Calculate degree difference of lat and long
	// Then convert the value to radians
	dLat := helper.ToRad(raw.End.Lat - raw.Start.Lat)
	dLong := helper.ToRad(raw.End.Long - raw.Start.Long)

	lat1 := helper.ToRad(raw.Start.Lat)
	lat2 := helper.ToRad(raw.End.Lat)

	// 1st step to calculate using haversine
	s := math.Pow(math.Sin(dLat/2), 2) + (math.Pow(math.Sin(dLong/2), 2) * math.Cos(lat1) * math.Cos(lat2))

	// 2nd step
	a := 2 * math.Asin(math.Sqrt(s))

	return a * EARTH_RADIUS
}
