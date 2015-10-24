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
