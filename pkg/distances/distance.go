package distances

import (
	"belimang/pkg/helper"
	"math"
)

type DistanceRaw struct {
	Start Point
	End   Point
}

type Point struct {
	Name string
	Lat  float64
	Long float64
}

type Vertex []Point

func CalculateArea(points Vertex) float64 {
	var area float64

	// Area is valid if data has more than 2 point
	if len(points) > 2 {
		for i := 0; i < len(points)-1; i++ {
			var (
				p1, p2 Point
			)
			// Check if the loop is reaching max element
			p1 = points[i]

			if i < len(points)-1 {
				p2 = points[i+1]
			} else {
				p2 = points[0]
			}

			area += helper.ToRad(p2.Long-p1.Long) * (2 + math.Sin(helper.ToRad(p1.Lat)) + math.Sin(helper.ToRad(p2.Lat)))
		}

		area = EARTH_RADIUS * area * EARTH_RADIUS / 2
	}

	return area
}

// Distance between two point
func pointDistance(p1, p2 Point) float64 {
	dx := p2.Long - p1.Long
	dy := p2.Lat - p1.Lat
	return math.Sqrt(float64(dx*dx + dy*dy))
}

func (t Vertex) distance() float64 {
	total := 0.0
	for i := 0; i < len(t)-1; i++ {
		total += pointDistance(t[i], t[i+1])
	}
	// total += pointDistance(t[len(t)-1], t[0]) // return to start
	return total
}

func ShortestDistance(points Vertex) (Vertex, float64) {
	shortest := points
	shortestDistance := points.distance()

	// Generate permutation
	var generate func([]Point, int)

	generate = func(p []Point, i int) {
		if i == 1 {
			d := Vertex(p).distance()

			if d < shortestDistance && p[0].Name == "user" {
				shortest = append(Vertex(nil), p...)
				shortestDistance = d
			}
		} else {
			for i := 0; i < i-1; i++ {
				generate(p, i-1)
				if i%2 == 0 {
					p[i], p[i-1] = p[i-1], p[i]
				} else {
					p[0], p[i-1] = p[i-1], p[0]
				}
			}
			generate(p, i-1)
		}
	}

	generate(points, len(points))

	return shortest, shortestDistance
}
