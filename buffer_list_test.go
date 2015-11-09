package buffer_list

import (
	"fmt"
	"reflect"
	"runtime"
	"testing"
)

var g_t *testing.T = nil
var enable_gc_check bool = true

type TestData struct {
	a int64
	b int32
	c int64
}

type TestDataPtr struct {
	a int
	b *TestData
}

type TestNestData struct {
	a int
	b TestDataPtr
}

func createList() *List {
	list := New(TestData{}, 4096)
	data := list.Front().Value().(*TestData)
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
	data := e.Value().(*TestData)

	data.a = 2
	data.b = 22

	if list.Len != 2 {
		t.Error("list.Len != 2")
	}

	data2 := list.Back().Value().(*TestData)

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
		data = e.Value().(*TestData)
		data.a = int64(i) * 1
		data.b = int32(i) * 11
	}

	if list.Len != 1000 {
		t.Error("list.len != 10")
	}

	data = list.Back().Prev().Value().(*TestData)

	if data.b != 998*11 {
		t.Error("data.b != 998*11", data.b)
	}
}

func TestBufferListConcurrentCreate1000(t *testing.T) {

	list := createList()
	fin := make(chan bool, 10)

	c_elm := func(list *List, i int, fin chan bool) {
		ee := list.InsertNewElem(list.Back())
		tdata := ee.Value().(*TestData)
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

func createData(l *List, f func(*Element, int)) {
	l.Back().Free()
	for i := 0; i < l.Cap()/2; i++ {
		e := l.InsertNewElem(l.Back())
		e.InitValue()
		f(e, i)
	}
}

func on_gc(d *TestData) {
	if enable_gc_check {
		g_t.Error("Error direct ")

	}
}

func on_gc2(d *TestData) {
	if enable_gc_check {
		g_t.Error(fmt.Sprintf("Error nested %p ", g_t))
	}
}

func TestProtectFreePtr(t *testing.T) {
	g_t = t
	tlist := New(TestDataPtr{}, 50000)

	createData(tlist, func(e *Element, i int) {
		v := e.Value().(*TestDataPtr)
		v.a = i
		v.b = &TestData{a: int64(i)}
		e.Commit()
		runtime.SetFinalizer(v.b, on_gc)
	})
	runtime.GC()
	//cnt := 0

	for e := tlist.Front(); e != nil; e = e.Next() {
		v := e.Value().(*TestDataPtr)
		runtime.SetFinalizer(v.b, nil)
		if !reflect.ValueOf(v.b).IsValid() {
			t.Error("FreePtr: data is freed")
		}

	}
}

func TestProtectFreeNestPtr(t *testing.T) {
	g_t = t

	tlist := New(TestNestData{}, 50000)

	createData(tlist, func(e *Element, i int) {
		v := e.Value().(*TestNestData)
		v.a = i
		v.b = TestDataPtr{a: i, b: &TestData{a: int64(i + 1)}}
		runtime.SetFinalizer(v.b.b, on_gc)
		e.Commit()
	})
	//	fmt.Println(tlist.Front().DumpPicks())

	runtime.GC()

	for e := tlist.Front(); e != nil; e = e.Next() {
		v := e.Value().(*TestNestData)

		runtime.SetFinalizer(v.b.b, nil)

		if !reflect.ValueOf(v.b.b).IsValid() {
			t.Error("FreeNestPtr: data is freed")
		}
	}
}

func TestAList(t *testing.T) {
	g_t = t

	tlist := New(TestNestData{}, 50000)
	alist := &AList{parent: tlist}
	createData(tlist, func(e *Element, i int) {
		if i%10 == 1 {
			ae := &AElement{list: alist, parent: e}
			alist.Push(ae)
		}
		v := e.Value().(*TestNestData)
		v.a = i
		v.b = TestDataPtr{a: i, b: &TestData{a: int64(i + 1)}}
		runtime.SetFinalizer(v.b.b, on_gc)
		e.Commit()
	})

	if alist.Len > 10000 {
		t.Error("invalid alist.Len")
	}
	cnt := 0
	for ae := alist.Front(); ae != nil; ae = ae.Next() {
		v := ae.Value().(*TestNestData)
		cnt++
		if !reflect.ValueOf(v.b.b).IsValid() {
			t.Error("cannot get value via AElement")
		}
	}
	if cnt == 0 {
		t.Error("not created AElement")
	}
}
