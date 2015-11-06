// Copyright 2015 Kazuhisa TAKEI<xtakei@me.com>. All rights reserved.
// Use of this source code is governed by MPL-2.0 license tha can be
// found in the LICENSE file

// Package buffer_list implements a double linked list with sequencial buffer data.
//
// To Get New First Data from buffer_list(l is a *List)
//		type Hoge Struct {
//			a int
//			b int
//		}
//		l := buffer_list.New(Hoge{})
//		hoge := l.GetElement(),Value().(*Hoge)
//		hoge.a = 1
//		hoge.b = 2
// To iterate over a list
//		for e := l.Front(); e != nil ; e = e.Next() {
//			a := (*Hoge)(e.Value())  // Hoge is Value type
//			// do something
//		}

package buffer_list

import (
	//	"fmt" // FIXME remove
	"reflect"
	"sync"
	"unsafe"
)

const (
	DEFAULT_BUF_SIZE = 1024
)

type Element struct {
	list      *List
	next      *Element
	prev      *Element
	old_value unsafe.Pointer
	value     interface{}
}

type List struct {
	Used      *Element
	Freed     *Element
	SizeElm   int64
	SizeData  int64
	Used_idx  int64
	Value_inf interface{}
	elms      []byte
	datas     []byte
	Len       int
	m         sync.Mutex
	cast_f    func(interface{}) interface{}
	pointers  map[uintptr]map[int]unsafe.Pointer
}

func New(first_value interface{}, buf_cnt int) (l *List) {
	l = new(List).Init(first_value, buf_cnt)
	l.pointers = make(map[uintptr]map[int]unsafe.Pointer)
	return l
}

func (e *Element) Commit() {
	e.List().Pick_ptr(e)
}

func (e *Element) free_pick_ptr() {
	l := e.list
	f_num := reflect.ValueOf(e.Value()).Elem().NumField()
	v_ptr := reflect.ValueOf(e.Value()).Pointer()

	if l.pointers[uintptr(v_ptr)] == nil {
		return
	}

	for i := 0; i < f_num; i++ {
		if l.pointers[uintptr(v_ptr)][i] != nil {
			//			fmt.Println("free value.member", i, l.pointers[uintptr(v_ptr)][i])

			delete(l.pointers[uintptr(v_ptr)], i)
		}
	}
}

func (l *List) Pick_ptr(e *Element) {
	f_num := reflect.ValueOf(e.Value()).Elem().NumField()
	v_ptr := reflect.ValueOf(e.Value()).Pointer()
	if l.pointers[uintptr(v_ptr)] == nil {
		l.pointers[uintptr(v_ptr)] = make(map[int]unsafe.Pointer)
	}

	for i := 0; i < f_num; i++ {
		m := reflect.ValueOf(e.Value()).Elem().Field(i)
		switch m.Kind() {
		case reflect.UnsafePointer:
			fallthrough
		case reflect.String:
			fallthrough
		case reflect.Slice:
			fallthrough
		case reflect.Map:
			fallthrough
		case reflect.Chan:
			fallthrough
		case reflect.Array:
			fallthrough
		case reflect.Func:
			fallthrough
		case reflect.Ptr:
			//			fmt.Println("detect ptr member", i, m.Kind(), m.Pointer())
			l.pointers[uintptr(v_ptr)][i] = unsafe.Pointer(m.Pointer())
		default:
		}

	}
}

func (l *List) GetDataPtr() uintptr {
	return uintptr(unsafe.Pointer(&l.datas[0]))
}
func (l *List) getElemData(idx int64) *Element {
	elm := (*Element)(unsafe.Pointer(&l.elms[int(l.SizeElm)*int(idx)]))
	elm.value = reflect.NewAt(l.TypeOfValue_inf(), unsafe.Pointer(&l.datas[int(l.SizeData)*int(idx)])).Interface()
	return elm
}
func (l *List) GetElement() *Element {
	return l.Used
}
func (e *Element) Next() *Element {
	e.list.m.Lock()
	defer e.list.m.Unlock()

	if e.next != nil {
		return e.next
	} else {
		return nil
	}
}

func (e *Element) Prev() *Element {
	e.list.m.Lock()
	defer e.list.m.Unlock()

	if e.prev != nil {
		return e.prev
	} else {
		return nil
	}
}

func (e *Element) Value() interface{} {
	return e.value
}

func (e *Element) Free() {

	e.list.m.Lock()
	defer e.list.m.Unlock()

	for ee := e.list.Used; ee != nil; ee = ee.next {
		if e == ee {
			goto DO_FREE
		}
	}

	//	fmt.Println("dont Free() e is not used ")
	return

DO_FREE:
	//	fmt.Println("do Free()")

	at := e.prev
	n := e.next
	if at.next == e {
		at.next = n
	}
	if n != nil {
		n.prev = at
	}

	e.list.Len -= 1

	if e.list.Used == e {
		e.list.Used = n
	}
	// move to free buffer
	if e.list.Freed == nil {
		e.prev = nil
		e.next = nil
		e.list.Freed = e
	} else {
		f_at := e.list.Freed
		e.next = f_at
		e.prev = nil
		f_at.prev = e
		e.list.Freed = e
	}
	e.free_pick_ptr()
}

func (e *Element) InitValue() {

	e.free_pick_ptr()

	diff := int(uint64(reflect.ValueOf(e.value).Pointer()) - uint64(uintptr(unsafe.Pointer(&e.list.datas[0]))))

	for i := range e.list.datas[diff : diff+int(e.list.SizeData)-1] {
		e.list.datas[diff+i] = 0
	}

	return
	//	fmt.Println(ref_byte, databyte)
}
func (l *List) newFirstElem() *Element {
	var e *Element

	//	l.m.Lock()
	//	defer l.m.Unlock()

	if l.Freed == nil {
		e = l.getElemData(l.Used_idx)
		l.Used_idx += 1
	} else {
		e = l.Freed
		if l.Freed.next == nil {
			l.Freed = nil
		} else {
			l.Freed = l.Freed.next
			l.Freed.prev = nil
		}
	}
	e.prev = e
	e.next = nil
	e.list = l
	if l.Used == nil {
		l.Used = e
	}
	l.Len++
	return e
}

func (l *List) InsertNewElem(at *Element) *Element {
	var e *Element

	l.m.Lock()
	defer l.m.Unlock()

	if l.Len == 0 && at == nil {
		return l.newFirstElem()
	}

	if l != at.list {
		return nil
	}

	if l.Freed == nil {
		e = l.getElemData(l.Used_idx)
		l.Used_idx += 1
	} else {
		e = l.Freed
		e.prev = nil
		e.next = nil
		if l.Freed.next == nil {
			l.Freed = nil
		} else {
			l.Freed = l.Freed.next
			l.Freed.prev = nil
		}
	}
	e.list = l
	n := at.next
	at.next = e
	e.prev = at
	if n != nil {
		n.prev = e
		e.next = n
	} else {
		e.list.Used.prev = e
	}

	l.Len++
	return e
}

func (l *List) TypeOfValue_inf() reflect.Type {
	if reflect.TypeOf(l.Value_inf).Kind() == reflect.Ptr {
		return reflect.ValueOf(l.Value_inf).Elem().Type()
	} else {
		return reflect.TypeOf(l.Value_inf)
	}
}

func (l *List) Init(first_value interface{}, value_len int) *List {
	l.m.Lock()
	defer l.m.Unlock()
	if l.Used == nil {
		var buf_len int64
		if value_len < 1024 {
			buf_len = int64(DEFAULT_BUF_SIZE)
		} else {
			buf_len = int64(value_len)
		}
		l.Value_inf = first_value
		l.SizeData = int64(l.TypeOfValue_inf().Size())
		l.SizeElm = int64(reflect.TypeOf(Element{}).Size())
		l.elms = make([]byte, buf_len*l.SizeElm,
			buf_len*l.SizeElm)
		l.datas = make([]byte, buf_len*l.SizeData,
			buf_len*l.SizeData)
		elm := (*Element)(unsafe.Pointer(&l.elms[0]))
		elm.value = reflect.NewAt(l.TypeOfValue_inf(), unsafe.Pointer(&l.datas[0])).Interface()
		elm.prev = elm
		elm.next = nil
		elm.list = l
		l.Used = elm
		l.Freed = nil
		l.Used_idx = 1
		l.Len = 1
	}
	return l
}

func (l *List) Front() *Element {
	l.m.Lock()
	defer l.m.Unlock()

	return l.Used
}

func (l *List) Back() *Element {
	l.m.Lock()
	defer l.m.Unlock()

	if l.Used == nil {
		return nil
	} else {
		return l.Used.prev
	}
}

func (l *List) Inf() interface{} {
	return l.Value_inf
}

func (l *List) Value() interface{} {
	return l.Used.value
}
func (l *List) SetCastFunc(f func(val interface{}) interface{}) {
	l.cast_f = f
}

func (e *Element) List() *List {
	return e.list
}

func (e *Element) ValueWithCast() interface{} {
	return e.list.cast_f(e.Value())
}
