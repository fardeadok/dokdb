package dokdb

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/google/uuid"
)

type object struct {
	// coords latitude, longitude
	lat, long float64
	// contetntype is MIME Content-type  text/html image/jpg
	contentType string
	description string
	// raw json byte
	jb []byte
}

// dokdb -  struct for store json
type db struct {
	// path to filename with raw json byte inside
	filename string
	// string is a random UUID
	store map[string]object
}

// ------------------------------------------------------
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

// ------------------------------------------------------
//
//	PRINT store
func (d *db) Print() {
	println("func Print store")

	for k, v := range d.store {
		println("uuid=", k)
		println(v.lat, v.long, v.contentType)
		println(v.jb)
	}
}

// ------------------------------------------------------
//
//	ADD json
//
// Add new json []byte  and return UUID
func (d *db) AddJson(lat, long float64, ct, ds string, jr []byte) (id string) {
	println("func AddJson")
	myUuid := uuid.New()

	myuuidString := myUuid.String()

	ov := object{}
	ov.lat = lat
	ov.long = long
	ov.contentType = ct
	ov.description = ds
	ov.jb = jr

	d.store[myuuidString] = ov

	return myuuidString
}

// ------------------------------------------------------
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
	// tmpJ, err := json.Marshal(dstore)
	if err != nil {
		return err
	}
	println("store is converted to tmpJ []byte. len=", len(tmpJ))

	writedLen, err := fd.Write(tmpJ)
	if err != nil {
		return err
	}

	println("writed to filename. len=", writedLen)
	println("func Save ok")
	return nil
}

// ------------------------------------------------------
//
//	LOAD
//
// Load - load from "filename" to "store"
func (d *db) Load(f string) (er error) {
	println("func Load json byte from filename")
	rawbytes, err := ioutil.ReadFile(d.filename)
	if err != nil {
		return err
	}

	err002 := json.Unmarshal(rawbytes, &d.store)
	if err002 != nil {
		return err002
	}

	println("func Load ok")
	return nil
}
