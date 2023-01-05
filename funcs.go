package dokdb

import "math"

//	--------------------
//
// return min, then max
func sortX(x1, x2 float64) (minX, maxX float64) {
	if x1 < x2 {
		return x1, x2
	}

	// else
	return x2, x1
}

// return true if  (minx < X < maxx)
func checkBetweenX(minX, maxX, x float64) bool {
	if minX < x && x < maxX {
		return true
	}

	return false
}

//	--------------------
//
// make new point from lat, long
func MakePoint(x, y float64) coords {
	return coords{Lat: x, Long: y}
}

//	--------------------
//
// check point in rect
func checkPointInRect(p1, p2 coords, p coords) bool {
	minx, maxx := sortX(p1.Lat, p2.Lat)
	betweenX := checkBetweenX(minx, maxx, p.Lat)
	if !betweenX {
		return false
	}

	miny, maxy := sortX(p1.Long, p2.Long)
	betweenY := checkBetweenX(miny, maxy, p.Long)

	return betweenY
}

// ------------------
//
// distance between 2 points.
// distance is float64 represents the raw
// number of meters
func distanceBetween(p1, p2 coords) float64 {

	value := 0.5 - math.Cos((p2.Lat-p1.Lat)*PiOver180)/2
	value += math.Cos(p1.Lat*PiOver180) * math.Cos(p2.Lat*PiOver180) * (1 - math.Cos((p2.Long-p1.Long)*PiOver180)) / 2

	return float64(DoubleEarthRadius * Distance(math.Asin(math.Sqrt(value))))

}

// ------------------
//
// check if point in radius
// radius is METERS!
func checkPointInradius(p1 coords, radius int64, p coords) bool {

	// meters
	dist001 := distanceBetween(p1, p)

	return dist001 <= float64(radius)

}
