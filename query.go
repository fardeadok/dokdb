package dokdb

import (
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
	list []object
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

// getid return one object with unique id (uuid)
// and delete all another from map
func (q *query) Id(id string) *query {
	q.Lock()
	defer q.Unlock()

	for k, v := range q.result {
		if v.Id != id {
			delete(q.result, k)
		}
	}

	q.updatelist()
	return q
}

// internal func for update list from map
func (q *query) updatelist() {
	println("updateList")

	q.list = nil
	for _, v := range q.result {
		q.list = append(q.list, v)
	}
	println("end updateList")
}

// sort by distance to point from min to max distance
func (q *query) Sort(p coords) {
	q.Lock()
	defer q.Unlock()
	sort.Slice(q.list,
		func(i int, j int) bool {
			p1 := q.list[i].coords
			p2 := q.list[j].coords
			d1 := DistanceBetween(p1, p)
			d2 := DistanceBetween(p2, p)
			return d1 < d2
		})
}

// Len return int number of objects in map
func (q *query) Len() int {
	return len(q.result)
}

// return results as []object
func (q *query) List() (ol []object) {
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
	q.updatelist()
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
	q.updatelist()
	return q
}

// fill result map with ALL data from db store
func (q *query) GetAll() *query {
	println(" func getall")
	q.Lock()
	defer q.Unlock()

	for k, v := range q.d.store {
		q.result[k] = v
	}

	println("  func getall end for")

	q.updatelist()

	println(" end func getall")
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
	q.updatelist()
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

	q.updatelist()
	println(" end func Where")
	return q
}

// print query result list
func (q *query) Print() *query {
	println("")
	println("  print list")
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
