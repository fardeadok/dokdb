package dokdb

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/tidwall/sjson"
)

// OBJECT IS RECORD IN DB
type object struct {
	// coords latitude, longitude
	// поля с большой буквы иначе не экспортируются и не записываются в файл
	// Lat  float64 `json:"lt"`
	// Long float64 `json:"lg"`
	coords
	// UUID  unique long id
	Id string `json:"id"`
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
	sync.Mutex
	store map[string]object
}

// ----------------------
//
//	NEW
//
// new make new *db
func New(fn string) *db {
	println("FUNC NEW 18:02")
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
	println("FUNC PRINT")

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
	println("FUNC ADDJSON")
	myUuid := uuid.New()
	myuuidString := myUuid.String()
	println(myuuidString)

	ov := object{}
	ov.Lat = lat
	ov.Long = long
	ov.Id = myuuidString
	ov.ContentType = ct
	ov.Description = ds

	str001 := strings.ReplaceAll(js, "\\", "")
	ov.Js = str001

	d.Lock()
	defer d.Unlock()
	d.store[myuuidString] = ov

	return myuuidString
}

// ------------------
//
//	UPDATE JSON
//
// update existing record by UUID
func (d *db) UpdateJson(id string, field string, newfalue string) (err error) {
	println("UPDATE JSOB by uuid")

	object001, ok := d.store[id]
	if ok != true {
		return errors.New("no id in db")
	}

	objectJs := object001.Js
	newJson, err := sjson.Set(objectJs, field, newfalue)
	if err != nil {
		return err
	}

	object001.Js = newJson

	d.Lock()
	defer d.Unlock()
	d.store[id] = object001

	return nil
}

//	 --------------------
//
//	SAVE
//
// Save "store map[string]placestruct to filename
func (d *db) Save() (er error) {
	println("FUNC SAVE")

	fd, err := os.Create(d.filename)
	if err != nil {
		return err
	}
	defer fd.Close()

	d.Lock()
	defer d.Unlock()

	// записываем красиво с отступами
	tmpJ, err := json.MarshalIndent(d.store, "", "\t")
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
	println("FUNC LOAD json byte from filename=", d.filename)

	rawbytes, err := os.ReadFile(d.filename)
	if err != nil {
		return err
	}

	d.Lock()
	defer d.Unlock()
	if err := json.Unmarshal(rawbytes, &d.store); err != nil {
		return err
	}

	println("func Load ok")
	return nil
}

//	 --------------------
//
//	find in rect
//
// find all objects in border and return slice
func (d *db) FindInRect(point1, point2 coords) (objectList []object) {
	println("FUNC FindInRect")

	for k, v := range d.store {
		if checkPointInRect(point1, point2, v.coords) {
			objectList = append(objectList, v)
			println()
			println("uuid=", k)
			println("lat long contenttype=", v.Lat, v.Long, v.ContentType)
			println("json=", v.Js)
		}
	}

	return objectList
}

// --------------------
//
// get all objects and return channel
func (d *db) GetAll_chan(in chan<- object) <-chan object {
	out := make(chan object)
	return out
}

//	 --------------------
//
//	find in rect
//
// find all objects in border and return chan
func (d *db) FindInRect_chan(point1, point2 coords, in chan<- object) <-chan object {
	out := make(chan object)
	return out
}