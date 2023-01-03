package dokdb

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/google/uuid"
)

type object struct {
	lat, long float64 `json:string` //coords
	// contetntype is MIME Content-type  text/html image/jpg
	contentType string
	// raw json byte
	jb []byte
}

// dokdb -  struct for store json
type dokdb struct {
	// path to filename with raw json byte inside
	filename string
	// string is a random UUID
	store map[string]object
}

// ------------------------------------------------------
//
//	NEW
//
// new make new db
func New(fn string) *dokdb {
	return &dokdb{
		filename: fn,
		store:    make(map[string]object),
	}
}

// ------------------------------------------------------
//
//	ADD object
//
// Add new object and return UUID
func (d *dokdb) Add(o object) (id string) {
	newuuid := uuid.New()
	println(newuuid.String())
	return newuuid.String()
}

// ------------------------------------------------------
//
//	SAVE
//
// Save "store map[string]placestruct to filename
func (d *dokdb) Save() (er error) {
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
func (d *dokdb) Load(f string) (er error) {
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
