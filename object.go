package dokdb

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

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

// check if object field==value return true
func (o *object) Equals(jspath, value string) bool {
	fieldValue := o.GetField(jspath)
	return value == fieldValue
}

// NewObject return new object with RANDOM uuid ID
func NewObject() object {
	return object{
		coords: coords{
			Lat:  0,
			Long: 0,
		},
		Id:          uuid.New().String(),
		ContentType: "",
		Description: "",
		Js:          "",
	}
}

// set value for object.Js
func (o *object) SetField(jspath string, newvalue string) error {
	println("func object.setfield")
	js := o.Js
	newjs, err := sjson.Set(js, jspath, newvalue)
	if err != nil {
		return err
	}
	o.Js = newjs
	return nil
}

// return string from object.js field
func (o *object) GetField(jspath string) string {
	gjs := gjson.Get(o.Js, jspath)
	if gjs.Exists() {
		return gjs.String()
	}
	return ""
}

// Object in rect? true or false
func (o *object) InRect(p1, p2 coords) bool {
	return o.coords.InRect(p1, p2)
}

// Object in radius R  true or false
func (o *object) InRadius(p1 coords, r int64) bool {
	return o.coords.InRadius(p1, r)
}

// Distance to another object
func (o *object) Distance(o1 object) float64 {
	return o.coords.Distance(o1.coords)
}

// fill object fields
func NewObjectFill(lat, long float64, ct, ds string, js string) object {
	u001 := uuid.New().String()
	ov := NewObject()
	ov.Id = u001

	ov.Lat = lat
	ov.Long = long

	ov.ContentType = ct
	ov.Description = ds
	str000 := strings.ReplaceAll(js, "  ", " ")
	str001 := strings.ReplaceAll(str000, "\t", " ")
	ov.Js = str001

	return ov
}

// Print
func (o *object) Print() {
	println("id=       ", o.Id)
	println("ContentType=", o.ContentType)
	fmt.Printf("latitude= %8.2f \n", o.Lat)
	fmt.Printf("longitude=%8.2f \n", o.Long)
	println("json=", o.Js)
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
