package dokdb

// OBJECT IS RECORD IN DB
type object struct {
	// coords latitude, longitude
	// поля с большой буквы иначе не экспортируются и не записываются в файл
	// Lat  float64 `json:"lt"`
	// Long float64 `json:"lg"`
	coords `json:"coords"`
	// UUID  unique long id
	Id string `json:"id"`
	// contetntype is MIME Content-type  text/html image/jpg
	ContentType string `json:"ct"`
	Description string `json:"ds"`
	// json string
	Js string `json:"js"`
}
