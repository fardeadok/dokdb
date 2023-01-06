package dokdb

import "fmt"

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

func (o *object) NewObject() object {
	return object{
		coords: coords{
			Lat:  0,
			Long: 0,
		},
		Id:          "",
		ContentType: "",
		Description: "",
		Js:          "",
	}
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
