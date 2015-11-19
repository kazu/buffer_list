// Copyright 2015 Kazuhisa TAKEI<xtakei@me.com>. All rights reserved.
// Use of this source code is governed by MPL-2.0 license tha can be
// found in the LICENSE file
//
// Package buffer_list implements a double linked list with sequencial buffer data.
//
// implement AList(alias list)
// example
//   bl :=  buffer_list.New(Hoge{}, 10000)
//   alist := Alist{parent: bl}
//   ae := alist.NewElem()
//   alist.Front().Insert(ae)
//   v := ae.Value().(*Hoge)
//   v.Commit()
//
// iteration
//      for e := range  alist.Generator() {
//			v := e.Value().(*Hoge)
//
//	  }
//

package buffer_list

import (
	"fmt"
	"sync"
)

type AList struct {
	m      sync.Mutex
	parent *List
	root   *AElement
	Len    int
	e2ae   map[*Element]*AElement
}

// AElement is an element of alias linked list.
type AElement struct {
	list   *AList
	prev   *AElement
	next   *AElement
	parent *Element
}

// create new alias element with element

func NewAList(l *List) AList {
	return AList{parent: l}
}

func (al *AList) NewElem() (ae *AElement) {
	ae = &AElement{list: al, parent: al.parent.InsertLast()}
	return ae
}

func (al *AList) SizeOfParentCache() int {
	return len(al.e2ae)
}

// add element to last of list
func (al *AList) Push(e *AElement) bool {

	al.m.Lock()
	defer al.m.Unlock()

	al.detect_empty_record("before_Push()")

	e.prev = e
	e.next = e
	if al.root == nil {
		al.root = e
	}
	at := al.root.prev
	n := at.next
	at.next = e
	e.prev = at
	e.next = n
	n.prev = e
	e.list = al
	al.Len++
	return true
}

// first element of list
func (al *AList) Front() *AElement {
	return al.root
}

// last element of list
func (al *AList) Back() *AElement {
	if al.root == nil {
		return nil
	}
	return al.root.prev
}

// add elemnet after at_element
func (at *AElement) Insert(e *AElement) *AElement {
	l := at.list

	if l == nil {
		fmt.Printf("e.Insert() fail a.list is nil e=%#v\n", at)
		return nil
	}

	l.m.Lock()
	defer l.m.Unlock()

	l.detect_empty_record("before_Insert()")

	n := at.next
	at.next = e
	e.prev = at
	e.next = n
	n.prev = e
	e.list = at.list

	l.Len++
	return e
}

// return value of reel element object
func (ae *AElement) Value() interface{} {
	if ae.list.e2ae == nil {
		ae.list.e2ae = make(map[*Element]*AElement)
	}
	ae.list.e2ae[ae.parent] = ae

	return ae.parent.Value()
}
func (ae *AElement) InitValue() {
	ae.parent.InitValue()
}

// register pointer protect GC
func (ae *AElement) Commit() {
	ae.parent.Commit()
}
func (l *AList) detect_empty_record(s string) {
	if l.root == nil {
		return
	}

	if l.root.list == nil || l.root.parent == nil || l.root.prev == nil || l.root.next == nil {
		fmt.Printf("detect invalid on %s al=%+v AList.root=%+v\n", s, l, l.root)
	}

}

// remove element from alias list
func (e *AElement) Remove() bool {

	if e.list == nil {
		fmt.Printf("e.Remove() fail e.list empty e=%#v\n", e)
		return false
	}
	e.list.m.Lock()
	defer e.list.m.Unlock()

	e.list.detect_empty_record("before_Remove()")

	delete(e.list.e2ae, e.parent)

	e.prev.next = e.next
	e.next.prev = e.prev
	e.list.Len -= 1
	if e.list.root == e {
		if e.list.Len == 0 {
			e.list.root = nil
		} else {
			e.list.root = e.next
		}
	}
	e.next = nil
	e.prev = nil

	e.list.detect_empty_record("After_Remove()")
	e.list = nil

	return true
}

// Next return next list element or nil
func (e *AElement) Next() *AElement {
	return e.next
}

// Prev return previous list element or nil
func (e *AElement) Prev() *AElement {
	return e.prev
}

// remove  from alias list and free real list
func (e *AElement) Free() {
	e.Remove()
	e.parent.Free()
	e.parent = nil
}

func (e *AElement) Base() *Element {
	return e.parent
}
func (e *AElement) List() *AList {
	return e.list
}

func (l *AList) ElemByValue(v interface{}) *AElement {

	e := l.parent.ElemByValue(v)
	if e == nil {
		fmt.Printf("WARN: ElemByValue() fail get parent elem v=%#v\n", v)
		return nil
	}
	if l.e2ae == nil {
		fmt.Printf("WARN: ElemByValue() fail a2ae not exits v=%#v e=%#e\n", v, e)
		return nil
	}

	if l.e2ae[e] == nil {
		fmt.Printf("WARN: ElemByValue() no entry in l.e2ae[e] v=%#v e=%#v\n", v, e)
	}

	return l.e2ae[e]
}

func (l *AList) Generator() chan *AElement {
	ch_size := l.Len

	if ch_size == 0 {
		ch_size = 1
	}

	ch := make(chan *AElement, ch_size)

	go func() {
		if l == nil {
			close(ch)
			return
		}
		cnt := 0
		e := l.Back()
		for {
			if e == nil || cnt > ch_size-1 {
				break
			}
			ch <- e
			e = e.Prev()
			cnt++
		}
		close(ch)
	}()

	return ch
}
