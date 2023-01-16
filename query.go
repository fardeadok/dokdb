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
	d      *db               //link to db
	result map[string]object //internal map of objects
	in     chan<- object
	out    <-chan object
	sync.Mutex
	list []object //слайс тоже переупорядочивается после каждого запроса!
	// pos  Next()  Offset() - change pos
	pos   int // pos=0 after updatelist()
	limit int // limit for List
}

//  РАБОТАТЬ С КАРТОЙ А НЕ []object
// ^^^ не получится тогда выдавать через Sort отсортированную карту- она не сортируется
// или можно сортировать? - ответЖ нет нельзя. карта преднамеренно рандомизируется при работе поэтому
// даже отсротированная карта будет смешана

// add field  pos int64  - current position

// только методы    .len()  .limit() .offset() работают с позицией по list[]object
// и меняют ее

// func Next() or Next(10)  {
// 	offset(pos)
// 	getfrom(pos)
// 	pos++
// }

// add  .Delete() - вызывать в конце для удаления элементов?

func (q *query) Limit(i int) *query {
	q.limit = i
	return q
}

// func (q *query) GetLimit(i int) int {
// 	return q.limit
// }

// limit - RETURN limit []object from list[]object
//
// from  [pos]    _do not modify pos_
//
// example:  query.Where().List() then query.Limit(10)
func (q *query) GetList(i int) ([]object, error) {

	// 0  1  2  3  4  5  6  7  8  9  len=10
	// [] [] [] [] [] [] [] [] [] []
	// pos=4       ^^
	// limit=3     4  5  6

	if q.pos < 0 || q.pos >= q.Len() {
		return nil, errors.New("out of range")
	}

	ol := make([]object, 0)

	for index := q.pos; index < q.Len(); index++ {
		o := q.list[index]
		ol = append(ol, o)
	}

	return ol, nil

}

// change pos+=offset
func (q *query) Offset(i int) *query {
	q.pos += i

	return q
}

// return list[pos] object and pos++
func (q *query) Next() (object, error) {
	// [] [] [] []
	// 0  1  2  3
	// len = 4 and max pos is 3
	if q.pos < 0 || q.pos >= q.Len() {
		return object{}, errors.New("out of range")
	}

	o := q.list[q.pos]
	q.pos++
	return o, nil
}

// get object[pos]
func (q *query) GetObject(i int) object {
	return q.list[q.pos]
}

// pos++
func (q *query) PosIncrease() *query {
	q.pos++
	return q
}

// pos--
func (q *query) PosDecrease() *query {
	q.pos--
	return q
}

// get position in list[]
func (q *query) GetPos() int {
	return q.pos
}

// set pos in list[]
func (q *query) SetPos(i int) *query {
	q.pos = i
	return q
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
	}

	return q
}

// filter - use jsonQ library as filter
// func (q *query) Filter(jq string) *query {
// 	q.Lock()
// 	defer q.Unlock()

// 	for k, v := range q.result {

// 		// jqfilter(........)

// 	}

// 	q.updatelist()
// 	return q
// }

//  func (q *query) ContentType(ct string) *query{

// 	return only ct types
//  }

//
//
//
//
//

// Id find one object with unique id (uuid)
// and delete all another
func (q *query) Id(id string) *query {
	q.Lock()
	defer q.Unlock()

	for k, v := range q.result {
		if v.Id != id {
			delete(q.result, k)
		}
	}

	// q.UpdateList()
	return q
}

// return object[id] from result map
func (q *query) GetId(id string) (object, error) {
	q.Lock()
	defer q.Unlock()

	if o, ok := q.result[id]; ok {
		return o, nil
	}

	return object{}, errors.New("no id")

}

// internal func for update list from map \n
// q.list[]=nil
// порядок элементов в [] меняется изза карты
func (q *query) UpdateList() *query {
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

// sort list[] by distance to point from min to max distance
func (q *query) Sort(p coords) *query {
	q.Lock()
	defer q.Unlock()

	q.UpdateList()

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

// Len return int number of objects in list[]ibject
func (q *query) Len() int {
	return len(q.list)
}

//	Return copy  []object
//
// GetList - use it after UpdateList() or Sort()
func (q *query) GetListAll() (ol []object) {
	// q.UpdateList()  // сбросит если перед ним было Sort() !!!
	return q.list
}

// return objects that in radius  point1, radius=r  meters
// and delete all another
func (q *query) InRadius(p1 coords, radius int64) *query {
	q.Lock()
	defer q.Unlock()
	for k, v := range q.result {
		if !v.InRadius(p1, radius) {
			delete(q.result, k)
		}
	}

	return q
}

// return objects that in rect point1, point2
// and delete all another
func (q *query) InRect(p1, p2 coords) *query {
	q.Lock()
	defer q.Unlock()
	for k, v := range q.result {
		if !v.InRect(p1, p2) {
			delete(q.result, k)
		}
	}

	return q
}

// fill result map with ALL data from db store
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

func (q *query) Start() *query {
	q.GetAll()
	return q
}

// get all where Field string contain substring.
// and remove all another objects from result map
func (q *query) Contain(jspath, substring string) *query {
	for k, v := range q.result {
		fieldValue := v.GetField(jspath) //get value of field
		if !strings.Contains(fieldValue, substring) {
			q.Lock()
			delete(q.result, k)
			q.Unlock()
		}
	}

	return q
}

// remove from map all != contentype
func (q *query) ContentType(value string) *query {
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

// remove from map all illegal objects
func (q *query) Where(jspath, value string) *query {
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

// print query result list
func (q *query) Print() *query {
	println("")
	println("  print list. ")

	// q.updatelist() - иначе лист обьектов перемешается при обновлении из карты

	if q.list == nil {
		return q
	}

	for k, v := range q.list {
		println(k)
		v.Print()
	}
	return q
}

// print Map

// print query result list
func (q *query) PrintMap() *query {
	println("")
	println("  print map")
	for k, v := range q.result {
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
