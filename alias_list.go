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

package buffer_list

import (
	"sync"
)

type AList struct {
	m      sync.Mutex
	parent *List
	root   *AElement
	Len    int
}

// AElement is an element of alias linked list.
type AElement struct {
	list   *AList
	prev   *AElement
	next   *AElement
	parent *Element
}

// create new alias element with element
func (al *AList) NewElem() (ae *AElement) {
	ae = &AElement{list: al, parent: al.parent.InsertLast()}
	return ae
}

// add element to last of list
func (al *AList) Push(ae *AElement) bool {
	al.m.Lock()
	defer al.m.Unlock()
	if al.root == nil {
		al.root = ae
		ae.prev = ae
		ae.next = nil
	} else {
		ae.prev = al.root.prev
		al.root.prev.next = ae
		al.root.prev = ae
		ae.next = nil
	}
	al.Len++
	return true
}

// first element of list
func (al *AList) Front() *AElement {
	return al.root
}

// last element of list
func (al *AList) Back() *AElement {
	return al.root.prev
}

// add elemnet after at_element
func (at *AElement) Insert(e *AElement) *AElement {
	l := at.list

	l.m.Lock()
	defer l.m.Unlock()

	e.list = at.list
	n := at.next
	at.next = e
	e.prev = at
	if n != nil {
		n.prev = e
		e.next = n
	} else {
		e.list.root.prev = e
	}
	l.Len++
	return e
}

// return value of reel element object
func (ae *AElement) Value() interface{} {
	return ae.parent.Value()
}

// register pointer protect GC
func (ae *AElement) Commit() {
	ae.parent.Commit()
}

// remove element from alias list
func (e *AElement) Remove() bool {
	e.list.m.Lock()
	defer e.list.m.Unlock()

	at := e.prev
	n := e.next
	if at.next == e {
		at.next = n
	}
	if n != nil {
		n.prev = at
	}
	e.list.Len -= 1
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
