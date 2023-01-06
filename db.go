package dokdb

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/tidwall/sjson"
)

// db -  struct for store json
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
// New make new *db
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
//	Print all store
func (d *db) Print() {
	println("FUNC PRINT")

	for k, v := range d.store {
		println("uuid=", k)
		printObject(v)
	}
}

// ----------------------
//
//	AddNewObjectFields
//
// AddNewObjectFields - add new json string  and return UUID
func (d *db) AddNewObjectFields(lat, long float64, ct, ds string, js string) (id string) {
	println("")
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

	str000 := strings.ReplaceAll(js, "  ", " ")
	str001 := strings.ReplaceAll(str000, "\t", " ")

	ov.Js = str001

	d.Lock()
	defer d.Unlock()
	d.store[myuuidString] = ov

	return myuuidString
}

// ------------------
//
//	ADD OBJECT without uuid inside
//	id will be added by db
//
//	add new object and return his new generated uuid
func (d *db) AddNewObject(o object) (id string, err error) {
	println("")
	println("func AddNewObjectFields")
	id001 := uuid.New()
	idString := id001.String()

	d.Lock()
	defer d.Unlock()
	d.store[idString] = o

	return idString, nil
}

// ------------------
//
//	UPDATE= search object by his id and update
//
//
//	update object by his id field
func (d *db) UpdateObject(o object) (err error) {
	println("")
	println("UPDATE existing object")
	// id001 := uuid.New()
	id := o.Id

	d.Lock()
	defer d.Unlock()
	d.store[id] = o

	return nil
}

// ------------------
//
//	FIND OBJECT BY UUID
//
//	return object
func (d *db) FindUUID(id string) (o object, err error) {
	println("")
	println("FIND by uuid")

	object001, ok := d.store[id]
	if !ok || id == "" {
		return object{
			coords: coords{
				Lat:  0,
				Long: 0,
			},
			Id:          id,
			ContentType: "",
			Description: "",
			Js:          "",
		}, errors.New("no id in db")
	}

	return object001, nil
}

// ------------------
//
//	UPDATE JSON
//
// update existing record by UUID
func (d *db) UpdateField(id string, field string, newfalue string) (err error) {
	println("UPDATE JSOB by uuid")

	object001, ok := d.store[id]
	if !ok {
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
	println("")
	println("FUNC SAVE")

	fd, err := os.Create(d.filename)
	if err != nil {
		return err
	}
	defer func(fd *os.File) {
		err := fd.Close()
		if err != nil {
			println("error close file")
		}
	}(fd)

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
	println("")
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
//	FIND IN RECT
//
// find all objects in border and return slice
func (d *db) FindInRect(point1, point2 coords) (objectList []object) {
	println("")
	println("FUNC FindInRect")

	for k, v := range d.store {
		if checkPointInRect(point1, point2, v.coords) {
			objectList = append(objectList, v)
			println()
			println("uuid=", k)
			printObject(v)
		}
	}

	return objectList
}

//	--------------------
//
//	FIND IN RADIUS
//
// return objects in radius (meters)
func (d *db) FindInRadius(point coords, radiusMeters int64) (objectList []object) {
	println("")
	println("FUNC FindInRadius")
	fmt.Printf("center lat= %8.2f    long=  %8.2f   \n", point.Lat, point.Long)
	println("radius=", radiusMeters)

	for k, v := range d.store {
		if checkPointInradius(point, radiusMeters, v.coords) {
			objectList = append(objectList, v)
			println("uuid=       ", k)
			printObject(v)
		}
	}

	return objectList
}

// --------------------
//
// get all objects and return channel
func (d *db) GetAll_chan(in chan<- object) <-chan object {
	println("")
	out := make(chan object)
	return out
}

//	 --------------------
//
//	find in rect
//
// find all objects in border and return chan
func (d *db) FindInRect_chan(point1, point2 coords, in chan<- object) <-chan object {
	println("")
	out := make(chan object)
	return out
}
