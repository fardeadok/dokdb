package dokdb

import (
	"fmt"
	"math"
)

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
func DistanceBetween(p1, p2 coords) float64 {
	println("	func distancebetween")

	value := 0.5 - math.Cos((p2.Lat-p1.Lat)*PiOver180)/2 + math.Cos(p1.Lat*PiOver180)*math.Cos(p2.Lat*PiOver180)*(1-math.Cos((p2.Long-p1.Long)*PiOver180))/2

	fmt.Printf("value=  %10.6f  \n", value)

	d := float64(DoubleEarthRadius * distance(math.Asin(math.Sqrt(value))))

	fmt.Printf("distance meters=  %8.2f \n", d)

	fmt.Printf("distance km=  %8.2f \n", d/1000)

	return d / 1000
}

// ------------------
//
// check if point in radius
// radius is METERS!
func checkPointInradius(p1 coords, radius int64, p coords) bool {
	println("")
	println("func checkpointinradius")
	fmt.Printf("center lat= %8.2f    long=  %8.2f   \n", p1.Lat, p1.Long)
	fmt.Printf("point  lat= %8.2f    long=  %8.2f   \n", p.Lat, p.Long)

	// meters
	dist001 := DistanceBetween(p1, p)
	return dist001 <= float64(radius)
}

// PRINT OBJECT
func printObject(o object) {
	println("")
	println("func printobject")
	println("uuid=       ", o.Id)
	println("ContentType=", o.ContentType)
	fmt.Printf("latitude= %8.2f \n", o.Lat)
	fmt.Printf("longitude=%8.2f \n", o.Long)
	println("json=", o.Js)
	println()
}
