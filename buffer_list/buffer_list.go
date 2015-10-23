// Copyright 2015 Kazuhisa TAKEI<xtakei@me.com>. All rights reserved.
// Use of this source code is governed by MPL-2.0 license tha can be
// found in the LICENSE file

package buffer_list

import (
	"reflect"
	"unsafe"
)

const (
	DEFAULT_BUF_SIZE = 1024
)

type Element struct {
	list  *List
	next  *Element
	prev  *Element
	value unsafe.Pointer
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
}

func New(first_value interface{}) *List {
	l := new(List)
	l.Init(first_value)
	return l
	//	return new(List).Init(value_struct)
}

func (l *List) getElemData(idx int64) *Element {
	elm := (*Element)(unsafe.Pointer(&l.elms[int(l.SizeElm)*int(idx)]))
	elm.value = unsafe.Pointer(&l.datas[int(l.SizeData)*int(idx)])
	return elm
}
func (l *List) GetElement() *Element {
	return l.Used
}
func (e *Element) Next() *Element {
	if e.next != nil {
		return e.next
	} else {
		return nil
	}
}

func (e *Element) Prev() *Element {
	if e.prev != nil {
		return e.prev
	} else {
		return nil
	}
}

func (e *Element) Value() unsafe.Pointer {
	return e.value
}

func (e *Element) Free() {
	at := e.prev
	n := e.next
	at.next = n
	n.prev = at

	e.list.Len -= 1
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
}
func (l *List) InsertNewElem(at *Element) *Element {
	var e *Element
	if l.Freed == nil {
		e = l.getElemData(l.Used_idx)
		l.Used_idx += 1
	} else {
		e = l.Freed
		if l.Freed.next == nil {
			l.Freed = nil
		} else {
			l.Freed = l.Freed.next
		}
	}
	n := at.next
	at.next = e
	e.prev = n
	if n != nil {
		n.prev = e
		e.list.Used.prev = e
	}
	e.list = l
	l.Len++
	return e
}

func (l *List) Init(first_value interface{}) *List {
	if l.Used == nil {
		l.Value_inf = first_value
		l.SizeData = int64(reflect.TypeOf(first_value).Size())
		l.SizeElm = int64(reflect.TypeOf(Element{}).Size())
		l.elms = make([]byte, DEFAULT_BUF_SIZE*l.SizeElm,
			DEFAULT_BUF_SIZE*l.SizeElm)
		l.datas = make([]byte, DEFAULT_BUF_SIZE*l.SizeData,
			DEFAULT_BUF_SIZE*l.SizeData)

		elm := (*Element)(unsafe.Pointer(&l.elms[0]))
		elm.value = unsafe.Pointer(&l.datas[0])
		elm.prev = nil
		elm.next = nil
		l.Used = elm
		l.Freed = nil
		l.Used_idx = 1
		l.Len = 1
	}
	return l
}

func (l *List) Front() *Element {
	return l.Used
}

func (l *List) Back() *Element {
	return l.Used.prev
}

func (l *List) Inf() interface{} {
	return l.Value_inf
}

func (l *List) Value() unsafe.Pointer {
	return l.Used.value
}
