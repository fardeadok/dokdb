package dokdb

import "math"

type Distance float64

const MilesPerKilometer = 0.6213712

const (
	Millimeter = Distance(0.001)
	Centimeter = Distance(0.01)
	Meter      = Distance(1)
	Kilometer  = Distance(1000)
	Mile       = Distance(1 / MilesPerKilometer * 1000)
)

const (
	EarthRadius       = 6371 * Kilometer
	DoubleEarthRadius = 2 * EarthRadius
	PiOver180         = math.Pi / 180
)
