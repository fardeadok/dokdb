package dokdb

import (
	"fmt"
	"math"
)

// COORDS IS A POINT(X,Y)
type coords struct {
	// latitude - X
	Lat float64 `json:"lt"`
	// longitude - Y
	Long float64 `json:"lg"`
}

type Point coords

func (c *coords) String() string {
	str001 := fmt.Sprintf(" %8.2f , ", c.Lat) + fmt.Sprintf(" %8.2f", c.Long)
	return str001
}

//	--------------------
//
// make new coords  from lat, long
func MakePoint(x, y float64) coords {
	return coords{Lat: x, Long: y}
}

//	--------------------
//
// make coords from lat, long
func Coords(x, y float64) coords {
	return coords{Lat: x, Long: y}
}

//	--------------------
//
// make coords from lat, long
func NewCoords(x, y float64) coords {
	return coords{Lat: x, Long: y}
}

//	--------------------
//
// check point in rect
func (c *coords) InRect(p1, p2 coords) bool {
	if c.BetweenLat(p1, p2) && c.BetweenLong(p1, p2) {
		return true
	}
	return false
}

//	--------------------
//
// return true if  (a.lat < c.Lat < b.lat)
func (c *coords) BetweenLat(a, b coords) bool {
	min, max := sortX(a.Lat, b.Lat)
	if min < c.Lat && c.Lat < max {
		return true
	}
	return false
}

//	--------------------
//
// return true if  (a.long < c.long < b.long)
func (c *coords) BetweenLong(a, b coords) bool {
	min, max := sortX(a.Long, b.Long)
	if min < c.Long && c.Long < max {
		return true
	}
	return false
}

// ------------------
//
// distance to another Point
// distance is float64 represents the raw
// number of meters
func (c *coords) Distance(p2 coords) float64 {
	// println("	func distance")
	value := 0.5 - math.Cos((p2.Lat-c.Lat)*PiOver180)/2 + math.Cos(c.Lat*PiOver180)*math.Cos(p2.Lat*PiOver180)*(1-math.Cos((p2.Long-c.Long)*PiOver180))/2
	// fmt.Printf("value=  %10.6f  \n", value)
	d := float64(DoubleEarthRadius * distance(math.Asin(math.Sqrt(value))))
	// fmt.Printf("distance meters=  %8.2f \n", d)
	// fmt.Printf("distance km=  %8.2f \n", d/1000)
	// return meters
	return d
}

// ------------------
//
// check if point in radius
// radius is METERS!
func (c *coords) InRadius(p1 coords, radius int64) bool {
	// println("")
	// println("func inradius")
	// fmt.Printf("center lat= %8.2f    long=  %8.2f   \n", p1.Lat, p1.Long)
	// fmt.Printf("point  lat= %8.2f    long=  %8.2f   \n", c.Lat, c.Long)
	// meters
	return c.Distance(p1) <= float64(radius)
}
