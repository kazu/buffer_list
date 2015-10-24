package buffer_list

import (
	"testing"
)

type TestData struct {
	a int64
	b int32
	c int64
}

func createList() *List {
	list := New(TestData{}, 4096)
	data := (*TestData)(list.Front().Value())
	data.a = 1
	data.b = 11

	return list
}

func TestBufferListiCreate(t *testing.T) {

	list := createList()

	if list.Len != 1 {
		t.Error("list.Len != 1")
	}
}

func TestBufferListInsertNewElem(t *testing.T) {

	list := createList()

	e := list.InsertNewElem(list.Front())
	data := (*TestData)(e.Value())

	data.a = 2
	data.b = 22

	if list.Len != 2 {
		t.Error("list.Len != 2")
	}

	data2 := (*TestData)(list.Back().Value())

	if data2.a != 2 {
		t.Error("data2.a != 2")
	}
}

func TestBufferListCreate1000(t *testing.T) {

	list := createList()
	var data *TestData
	var e *Element
	for i := 1; i < 1000; i++ {
		e = list.InsertNewElem(list.Back())
		data = (*TestData)(e.Value())
		data.a = int64(i) * 1
		data.b = int32(i) * 11
	}

	if list.Len != 1000 {
		t.Error("list.len != 10")
	}

	data = (*TestData)(list.Back().Prev().Value())

	if data.b != 998*11 {
		t.Error("data.b != 998*11", data.b)
	}
}

func TestBufferListConcurrentCreate1000(t *testing.T) {

	list := createList()
	fin := make(chan bool, 10)

	c_elm := func(list *List, i int, fin chan bool) {
		ee := list.InsertNewElem(list.Back())
		tdata := (*TestData)(ee.Value())
		tdata.a = int64(i) * 1
		tdata.b = int32(i) * 11
		fin <- true
	}

	for i := 1; i < 1000; i++ {
		go c_elm(list, i, fin)
	}

	for i := 0; i < 999; i++ {
		select {
		case <-fin:
			continue
		}
	}

	if list.Len != 1000 {
		t.Error("list.len != 10", list.Len)
	}
}
