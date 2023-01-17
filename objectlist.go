package dokdb

import (
	"errors"
	"sort"
	"sync"
)

// query возвращает этот обьект
// и заполняет его данными
type objectList struct {
	sync.Mutex
	list []object
	// only Next() change pos
	pos int
	// Limit()
	limit int
	// Offset()
	offset int
}

// return *ol
func NewObjectList(dP *db) *objectList {
	ol := &objectList{
		// d:      dP,
		Mutex:  sync.Mutex{},
		list:   []object{},
		pos:    0,
		limit:  0,
		offset: 0,
	}
	return ol
}

// upload list to d.store
func (ol *objectList) UpdateToDB() *objectList {
	return ol
}

//
//
//
//
//

func (ol *objectList) Filter(f string) *objectList {
	// filter - use jsonQ library as filter
	// func (q *query) Filter(jq string) *query {
	// 	q.Lock()
	// 	defer q.Unlock()
	return ol
}

//
//
//
//
//

// set list=nil, pos,limit,offset=0
func (ol *objectList) Reset() *objectList {
	ol.Lock()
	defer ol.Unlock()

	ol.list = nil
	ol.pos, ol.limit, ol.offset = 0, 0, 0
	return ol
}

//
//
//
//
//

func (q *objectList) Limit(i int) *objectList {
	q.limit = i
	return q
}

//
//
//
//
//

// return []objects with (offset,limit)
func (q *objectList) GetN(offset, limit int) ([]object, error) {
	q.offset = offset
	q.limit = limit
	return q.Get()
}

func (ol *objectList) GetAll() ([]object, error) {

	return ol.list, nil
}

//
//
//
//
//

// возвращает результаты из списка с отступом offset в количестве limit
// и не меняет значение pos, offset, limit
func (q *objectList) Get() ([]object, error) {

	// 0  1  2  3   4    5   6    7  8  9  len=10
	// [] [] [] []  []  []  []   [] [] []
	// ofset=0..5   ^^ - по сути начальное значение индекса
	// limit=3      3    2   1  - количество обьектов возврата

	if q.offset < 0 || q.offset >= len(q.list) {
		return nil, errors.New("out of range offset")
	}

	if q.offset+q.limit > len(q.list) {
		return nil, errors.New("out of range offset+limitoffset")
	}

	// или можно  ol []object
	ol := make([]object, 0)

	o := q.offset
	l := q.limit

	for index := o; index < len(q.list); index++ {
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

// change pos+=offset
func (q *objectList) Offset(i int) *objectList {
	q.offset = i
	return q
}

//
//
//
//

// return list[pos] object and pos++
func (q *objectList) Next() (*object, error) {
	// [] [] [] []
	// 0  1  2  3
	// len = 4 and max pos is 3
	if q.pos < 0 || q.pos >= len(q.list) {
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

// get *object[pos]
func (q *objectList) GetObjectIndex(i int) (object, error) {
	if i < 0 || i >= len(q.list) {
		return object{}, errors.New("out of range")
	}

	return q.list[i], nil
}

//

//
//
//
//
//

// pos++
func (q *objectList) PosIncrease() *objectList {
	q.pos++
	return q
}

//
//
//
//
//

// pos--
func (q *objectList) ListPosDecrease() *objectList {
	q.pos--
	return q
}

// return *object[id] from result map
func (q *objectList) GetByID(id string) (object, error) {
	q.Lock()
	defer q.Unlock()

	for _, v := range q.list {
		if v.Id == id {
			return v, nil
		}
	}

	return object{}, errors.New("no id")

}

//
//
//
//
//

func (ol *objectList) UpdateFromQuery(q *query) *objectList {

	ol.list = nil
	ol.pos = 0
	ol.limit = 1
	ol.offset = 0

	for _, v := range q.result {
		ol.list = append(ol.list, v)
	}

	return ol
}

//
//
//
//
//

// sort list[] by distance to point from min to max distance
func (ol *objectList) SortByDistance(p coords) *objectList {
	ol.Lock()
	defer ol.Unlock()

	// q.ListUpdate()

	sort.Slice(ol.list,
		func(i int, j int) bool {
			p1 := ol.list[i].coords
			p2 := ol.list[j].coords
			d1 := DistanceBetween(p1, p)
			d2 := DistanceBetween(p2, p)
			return d1 < d2
		})

	return ol
}

//
//
//
//
//

// Len return int number of objects in list[]ibject
func (ol *objectList) Len() int {
	return len(ol.list)
}

//
//
//
//
//

// return objects that in radius  point1, radius=r  meters
// and delete all another
func (ol *objectList) GetInRadius(p1 coords, radius int64) ([]object, error) {
	ol.Lock()
	defer ol.Unlock()

	rl := make([]object, 0)

	for _, v := range ol.list {
		if v.InRadius(p1, radius) {
			rl = append(rl, v)
		}
	}

	return rl, nil
}

//
//
//
//
//

func (ol *objectList) GetInRect(p1, p2 coords) ([]object, error) {
	ol.Lock()
	defer ol.Unlock()

	rl := make([]object, 0)

	for _, v := range ol.list {
		if v.InRect(p1, p2) {
			rl = append(rl, v)
		}
	}

	return rl, nil
}

//
//
//
//
//

func (ol *objectList) GetContain(field, substr string) ([]object, error) {
	ol.Lock()
	defer ol.Unlock()

	rl := make([]object, 0)

	for _, v := range ol.list {

		if v.Contain(field, substr) {
			rl = append(rl, v)
		}
	}

	return rl, nil
}

func (ol *objectList) GetByFieldValue(field, value string) ([]object, error) {
	ol.Lock()
	defer ol.Unlock()

	rl := make([]object, 0)

	for _, v := range ol.list {

		if v.Equals(field, value) {
			rl = append(rl, v)
		}
	}

	return rl, nil
}

func (ol *objectList) Print() *objectList {

	if ol.list == nil {
		return ol
	}

	for _, v := range ol.list {
		v.Print()
	}

	return ol
}
