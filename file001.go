package dokdb

import (
	"encoding/json"
	"os"
)

type placeStruct struct {
	lat, long float64 `json:string` //coords
	content   string  // type MIME Content-type  text/html image/jpg
	jb        []byte  //raw json
	// uj        map[string]interface{} //unmarshalled jsonb
}

// dokdb -  struct for store json
type dokdb struct {
	filename string
	fd       *os.File
	store    map[string]placeStruct
	//  добавить структуру с координатами
}

// Save "store map[string]placestruct to filename
func (d *dokdb) Save() (er error) {
	println("func dokdb save")
	fd, err := os.Create(d.filename)

	if err != nil {
		return err
	}
	defer fd.Close()

	tmpJ, err := json.Marshal(d.store)
	if err != nil {
		return err
	}
	println("store is converted to tmpJ []byte. len=", len(tmpJ))

	writedLen, err := fd.Write(tmpJ)
	if err != nil {
		return err
	}

	println("writed to filename. len=", writedLen)

	return nil
}

// Load - load from "filename" to "store"
func (d *dokdb) Load(f string) (er error) {
	return nil
}
