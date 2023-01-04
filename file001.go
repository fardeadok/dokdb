package dokdb

import (
	"encoding/json"
	"os"

	"github.com/google/uuid"
)

type object struct {
	// coords latitude, longitude
	// поля с большой буквы иначе не экспортируются и не записываются в файл
	Lat  float64 `json:"lt"`
	Long float64 `json:"lg"`
	// contetntype is MIME Content-type  text/html image/jpg
	ContentType string `json:"ct"`
	Description string `json:"ds"`
	// json string
	Js string `json:"js"`
}

// dokdb -  struct for store json
type db struct {
	// path to filename with raw json byte inside
	filename string
	// string is a random UUID
	store map[string]object
}

// ----------------------
//
//	NEW
//
// new make new *db
func New(fn string) *db {
	return &db{
		filename: fn,
		store:    make(map[string]object),
	}
}

// ----------------------
//
//	PRINT
//
//	print all store
func (d *db) Print() {
	println("func Print store v1549")

	for k, v := range d.store {
		println("uuid=", k)
		println("lat long contenttype=", v.Lat, v.Long, v.ContentType)
		println("json=", v.Js)
		println()
	}
}

// ----------------------
//
//	ADD json
//
// Add new json string  and return UUID
func (d *db) AddJson(lat, long float64, ct, ds string, js string) (id string) {
	println("func AddJson")
	myUuid := uuid.New()

	myuuidString := myUuid.String()

	ov := object{}
	ov.Lat = lat
	ov.Long = long
	ov.ContentType = ct
	ov.Description = ds
	ov.Js = js

	d.store[myuuidString] = ov

	return myuuidString
}

//	 --------------------
//
//	SAVE
//
// Save "store map[string]placestruct to filename
func (d *db) Save() (er error) {
	println("func dokdb save")

	fd, err := os.Create(d.filename)
	if err != nil {
		return err
	}
	defer fd.Close()

	// записываем красиво с отступами
	tmpJ, err := json.MarshalIndent(d.store, "", " ")
	if err != nil {
		return err
	}

	println("print string(tmpJ)")
	println(string(tmpJ))
	println("")

	println("store is converted to tmpJ []byte. len=", len(tmpJ))

	writedLen, err := fd.Write(tmpJ)
	if err != nil {
		return err
	}

	println("writed to filename. len=", writedLen)
	println("func Save ok")
	return nil
}

//	--------------------
//
//	LOAD
//
// Load - load from "filename" to "store"
func (d *db) Load() (er error) {
	println("func Load json byte from filename")

	rawbytes, err := os.ReadFile(d.filename)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(rawbytes, &d.store); err != nil {
		return err
	}

	println("func Load ok")
	return nil
}
