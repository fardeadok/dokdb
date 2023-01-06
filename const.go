package dokdb

import "math"

type distance float64

const (
	MilesPerKilometer = 0.6213712

	Millimeter = distance(0.001)
	Centimeter = distance(0.01)
	Meter      = distance(1)
	Kilometer  = distance(1000)
	Mile       = distance(1 / MilesPerKilometer * 1000)

	//	meters
	EarthRadius = 6371 * Kilometer
	//	meters
	DoubleEarthRadius = 2 * EarthRadius
	PiOver180         = math.Pi / 180
)

const (
	YachtSailing     = "yacht/sail"
	YachtMotor       = "yacht/motor"
	BoatMotor        = "boat/motor"
	CatamaranSailing = "catamaran/sail"
	CatamaranMotor   = "catamaran/motor"
	BoatHouse        = "boat/house"
	BoatHotel        = "boat/hotel"
)
