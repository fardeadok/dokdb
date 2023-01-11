package dokdb

// test sync 18:52

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

// db -  struct for store json
type db struct {
	// path to filename with raw json byte inside
	filename string
	// string is a random UUID
	sync.Mutex
	store map[string]object
}

// find all object with field=value
//
//	FindByField
func (d *db) FindByField(jspath, value string) []object {
	var objectList []object
	println("")
	println("FUNC Find by field")

	for k, v := range d.store {
		if v.Equals(jspath, value) {
			objectList = append(objectList, v)
			// println()
			println("uuid=", k)
			// printObject(v)
		}
	}

	return objectList

}

// ----------------------
//
//	NEW
//
// fn is filename.
// New make new *db.
// Then us Load for load file in db
func NewDB(fn string) *db {
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
	// println("")
	// println("FUNC ADDJSON")

	myUuid := uuid.New()
	myuuidString := myUuid.String()
	// println(myuuidString)

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

//	-----------------
//
//	AddObject
//
// if object id == "" then generate id and save to db store
func (d *db) AddObject(o object) (id string, err error) {
	println("")
	println("ADD object to db")

	if o.Id == "" {
		o.Id = uuid.New().String()
	}

	d.Lock()
	defer d.Unlock()
	d.store[o.Id] = o

	return o.Id, nil
}

// ------------------
//
//	UPDATE
//
//	search object by his id and update
func (d *db) UpdateObject(o object) (err error) {
	println("")
	println("UPDATE existing object")

	d.Lock()
	defer d.Unlock()
	d.store[o.Id] = o

	return nil
}

// ------------------
//
//	FIND OBJECT BY UUID
//
//	return object or make New object if empty
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

//	-----------------
//
// get value from object fild
func (d *db) GetField(id string, jspath string) (string, error) {
	println("func getfield")
	o001, ok := d.store[id]
	if !ok {
		return "", errors.New("no id in db")
	}
	return o001.GetField(jspath), nil
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

	err = object001.SetField(field, newfalue)
	if err != nil {
		return err
	}

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

	d.Lock()
	defer d.Unlock()

	fd, err := os.Create(d.filename)
	if err != nil {
		return err
	}
	// defer func(fd *os.File) {
	// 	err := fd.Close()
	// 	if err != nil {
	// 		println("error close file")
	// 	}
	// }(fd)

	// записываем красиво с отступами
	tmpJ, err := json.MarshalIndent(d.store, "", "\t")
	if err != nil {
		return err
	}

	// println("print string(tmpJ)")
	// println(string(tmpJ))
	// println("")

	// println("store is converted to tmpJ []byte. len=", len(tmpJ))

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

//	--------------------
//
// # FIND IN RECT
//
// find all objects in border and return slice
func (d *db) FindInRect(point1, point2 coords) (objectList []object) {
	println("")
	println("FUNC FindInRect")

	for k, v := range d.store {
		if v.InRect(point1, point2) {
			objectList = append(objectList, v)
			println()
			println("uuid=", k)
			printObject(v)
		}
	}

	return objectList
}

//	-------------------
//
// sort objects in rect by distance to p1
func (d *db) SortNearest(p coords, ol []object) {

	sort.Slice(ol,
		func(i int, j int) bool {
			p1 := ol[i].coords
			p2 := ol[j].coords
			d1 := DistanceBetween(p1, p)
			d2 := DistanceBetween(p2, p)
			return d1 < d2
		})
}

//	--------------------
//
// # FIND IN RADIUS
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

//	------------------
//
// # GetAll
//
//	get all []object
func (d *db) GetAll() (ol []object) {
	println("")
	println("func GetAll")
	for k, v := range d.store {
		ol = append(ol, v)
		println(k)
	}
	return ol
}

// --------------------
//
// # getall_chan
//
// get all objects and return channel
func (d *db) GetAll_chan() chan object {

	time.Sleep(3 * time.Second)
	println("")
	println("func GetAll_chan")

	out := make(chan object)

	go func() {
		for k, v := range d.store {
			println("	for k,v := range d.store   key=", k)
			v.Print()
			out <- v
		}
		close(out)
	}()

	println("")
	println("func GetAll_chan exit")
	return out

}

// for k, v := range chan001 {
// 	if o, ok := <-v; ok == true {
// 		o.Print()
// 	}
// }

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
