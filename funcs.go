package dokdb

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

func MakePoint(x, y float64) coords {
	return coords{Lat: x, Long: y}
}

func checkPointInRect(p1, p2 coords, p coords) bool {
	minx, maxx := sortX(p1.Lat, p2.Lat)
	betweenX := checkBetweenX(minx, maxx, p.Lat)
	if betweenX == false {
		return false
	}

	miny, maxy := sortX(p1.Long, p2.Long)
	betweenY := checkBetweenX(miny, maxy, p.Long)
	if betweenY == false {
		return false
	}

	return true
}
