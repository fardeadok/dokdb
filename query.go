// query.Start().UpdateList().GetListAll()
//
// query.GetId(id)  - return [id]object
//
// query.Id(id) - return [id] and delete all another
//
// query.Offset(2).Limit(3).UpdateList().GetList()
//
// query.SetPos(0).Limit(5).GetList()
//
// query.Sort().GetListAll()
//
// query.Inradius(5
// )
// o := query.Next() - return object and pos++
package dokdb

import (
	"errors"
	"sort"
	"strings"
	"sync"
)

// 01:08

// QUERY FOR ASK FROM DB
// query DO NOT MODUFY db with Sort,List,Contain etc...
// use UpdateDB for upliad results from query to db
type query struct {
	d *db //link to db
	// УБРАТЬ КАРТУ И РАБОТАТЬ НАПРЯМУЮ С d.store
	result map[string]object //internal map of objects
	in     chan<- object
	out    <-chan object
	sync.Mutex
	list []object
	// только Next() меняет позицию.
	// pos=0 after updatelist() Sort()
	pos int
	// переменные для Offset().Limit().GetList()
	limit  int // limit for List
	offset int // offset for GetList
}

// return new  *query
func NewQuery(dP *db) *query {
	q := &query{
		d:      dP,
		result: make(map[string]object),
		in:     make(chan<- object),
		out:    make(<-chan object),
		Mutex:  sync.Mutex{},
		list:   []object{},
		pos:    0,
		limit:  0,
		offset: 0,
	}

	return q
}

//
//
//
//
//

// upload list to d.store
func (q *query) UpdateToDB() *query {

	return q
}

//
//
//
//
//

func (q *query) Filter(f string) *query {
	// filter - use jsonQ library as filter
	// func (q *query) Filter(jq string) *query {
	// 	q.Lock()
	// 	defer q.Unlock()
	return q
}

//
//
//
//
//

// set list=nil, pos,limit,offset=0
func (q *query) Reset() *query {
	q.Lock()
	defer q.Unlock()

	q.list = nil
	q.pos, q.limit, q.offset = 0, 0, 0
	return q
}

//
//
//
//
//

func (q *query) Start() *query {
	q.GetAll()
	return q
}

//
//
//
//
//

// fill result map with ALL data from db store
// УБРАТЬ ФУНКЦИЮ И РАБОТАТЬ НАПРЯМУЮ С d.store через ссылку
func (q *query) GetAll() *query {
	println(" func getall")
	q.Lock()
	defer q.Unlock()

	//clear map
	for k := range q.result {
		delete(q.result, k)
	}

	for k, v := range q.d.store {
		q.result[k] = v
	}

	// clear []object
	q.list = nil

	println(" end func getall")
	return q
}

//
//
//
//
//

func (q *query) ListLimit(i int) *query {
	q.limit = i
	return q
}

//
//
//
//
//

// return []objects with (offset,limit)
func (q *query) ListGetN(offset, limit int) ([]object, error) {
	q.offset = offset
	q.limit = limit
	return q.ListGet()
}

//
//
//
//
//

// возвращает результаты из списка с отступом offset в количестве limit
// и не меняет значение pos, offset, limit
func (q *query) ListGet() ([]object, error) {

	// 0  1  2  3   4    5   6    7  8  9  len=10
	// [] [] [] []  []  []  []   [] [] []
	// ofset=0..5   ^^ - по сути начальное значение индекса
	// limit=3      3    2   1  - количество обьектов возврата

	if q.offset < 0 || q.offset >= q.Len() {
		return nil, errors.New("out of range offset")
	}

	if q.offset+q.limit > q.Len() {
		return nil, errors.New("out of range offset+limitoffset")
	}

	// или можно  ol []object
	ol := make([]object, 0)

	o := q.offset
	l := q.limit

	for index := o; index < q.Len(); index++ {
		ol = append(ol, q.list[index])
		l--
		if l == 0 {
			break
		}
	}

	return ol, nil

}

//
//
//
//
//

// change pos+=offset
func (q *query) ListOffset(i int) *query {
	q.offset = i
	return q
}

//
//
//
//
//

// return list[pos] object and pos++
func (q *query) ListNextObject() (*object, error) {
	// [] [] [] []
	// 0  1  2  3
	// len = 4 and max pos is 3
	if q.pos < 0 || q.pos >= q.Len() {
		return nil, errors.New("out of range")
	}

	o := &q.list[q.pos]
	q.pos++

	return o, nil
}

//
//
//
//
//

// get *object[pos]
func (q *query) ListGetObject(i int) *object {
	if i < 0 || i >= len(q.list) {
		return nil
	}

	return &q.list[i]
}

//
//
//
//
//

// pos++
func (q *query) ListPosIncrease() *query {
	q.pos++
	return q
}

//
//
//
//
//

// pos--
func (q *query) ListPosDecrease() *query {
	q.pos--
	return q
}

//
//
//
//
//

//
//
//
//
//

// return *object[id] from result map
func (q *query) FidnID(id string) (object, error) {
	q.Lock()
	defer q.Unlock()

	if o, ok := q.result[id]; ok {
		return o, nil
	}

	return object{
		coords: coords{
			Lat:  0,
			Long: 0,
		},
		Id:          "",
		ContentType: "",
		Description: "",
		Js:          "",
	}, errors.New("no id")

}

//
//
//
//
//

// internal func for update list from map \n
// q.list[]=nil
// порядок элементов в [] меняется изза карты
func (q *query) ListUpdate() *query {
	println("updateList")
	q.pos = 0
	// clear list[]
	q.list = nil

	for _, v := range q.result {
		q.list = append(q.list, v)
	}
	println("end updateList")

	return q
}

//
//
//  !SORT ПУСТЬ ВОЗВРАЩАЕТ НОВЫЙ []object - не из себя
//
//

// sort list[] by distance to point from min to max distance
func (q *query) ListSort(p coords) *query {
	q.Lock()
	defer q.Unlock()

	q.ListUpdate()

	sort.Slice(q.list,
		func(i int, j int) bool {
			p1 := q.list[i].coords
			p2 := q.list[j].coords
			d1 := DistanceBetween(p1, p)
			d2 := DistanceBetween(p2, p)
			return d1 < d2
		})

	return q
}

//
//
//
//
//

// Len return int number of objects in list[]ibject
func (q *query) Len() int {
	return len(q.list)
}

//
//
//
//
//

//	Return copy  []object
//
// GetList - use it after UpdateList() or Sort()
func (q *query) GetListAll() (ol []object) {
	// q.UpdateList()  // сбросит если перед ним было Sort() !!!
	return q.list
}

//
//
//
//
//

// return objects that in radius  point1, radius=r  meters
// and delete all another
func (q *query) WhereInRadius(p1 coords, radius int64) *query {
	q.Lock()
	defer q.Unlock()
	for k, v := range q.result {
		if !v.InRadius(p1, radius) {
			delete(q.result, k)
		}
	}

	return q
}

//
//
//
//
//

// return objects that in rect point1, point2
// and delete all another
func (q *query) WhereInRect(p1, p2 coords) *query {
	q.Lock()
	defer q.Unlock()
	for k, v := range q.result {
		if !v.InRect(p1, p2) {
			delete(q.result, k)
		}
	}

	return q
}

//
//
//
//
//

// get all where Field string contain substring.
// and remove all another objects from result map
func (q *query) WhereContain(jspath, substring string) *query {
	println("")
	println(" func contain", jspath, substring)

	for k, v := range q.result {

		v.Print()

		fieldValue := v.GetField(jspath) //get value of field
		println("country=", fieldValue)

		if !strings.Contains(fieldValue, substring) {
			q.Lock()
			delete(q.result, k)
			q.Unlock()
		}

	}

	return q
}

//
//
//
//
//

// remove from map all != contentype
func (q *query) WhereContentType(value string) *query {
	println("  func  ContentType= ", value)

	for k, v := range q.result {
		if v.ContentType != value {
			q.Lock()
			delete(q.result, k)
			q.Unlock()
		}
	}

	println(" end func ContentType")
	return q
}

//
//
//
//
//

// WhereId find one object with unique id (uuid)
// and delete all another
func (q *query) WhereID(id string) *query {
	println("  func  WhereID ", id)

	for k := range q.result {
		if k != id {
			q.Lock()
			delete(q.result, k)
			q.Unlock()
		}
	}

	println(" end func WhereID")
	return q
}

// remove from map all illegal objects
func (q *query) WhereField(jspath, value string) *query {
	println("  func  Where ", jspath, " = ", value)
	for k, v := range q.result {
		if !v.Equals(jspath, value) {
			q.Lock()
			delete(q.result, k)
			q.Unlock()
		}
	}

	println(" end func Where")
	return q
}

//
//
//
//
//

// print query result list
func (q *query) Print() *query {
	println("")
	println("  print list. ")

	// если q.updatelist() то  обьекты  перемешается при обновлении из карты

	if q.list == nil {
		return q
	}

	for k, v := range q.list {
		println(k)
		v.Print()
	}
	return q
}

//
//
//
//
//

// print Map

// print query result list
func (q *query) PrintMap() *query {
	println("")
	println("  print map")
	for k, v := range q.result {
		println("")
		println(k)
		v.Print()
	}
	return q
}

// go func() {
// 	for _, n := range nums {
// 		out <- n
// 	}
// 	close(out)
// }()
