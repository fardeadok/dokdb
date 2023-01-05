package dokdb

// RECT IS RECTANGULAR
type rect struct {
	// upper left point
	point1 coords
	// bottom right point
	point2 coords
}

// COORDS IS A POINT(X,Y)
type coords struct {
	// latitude - X
	Lat float64 `json:"lt"`
	// longitude - Y
	Long float64 `json:"lg"`
}
