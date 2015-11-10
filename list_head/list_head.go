// Copyright 2015 Kazuhisa TAKEI<xtakei@me.com>. All rights reserved.
// Use of this source code is governed by MPL-2.0 license tha can be
// found in the LICENSE file
//
// list_head is kernel list_head like double linked list
//

package list_head

import (
	"fmt"
	"reflect"
	"unsafe"
)

var diff map[string]uintptr = make(map[string]uintptr)
var diff_type map[string]bool = make(map[string]bool)

type ListHead struct {
	Prev *ListHead
	Next *ListHead
}

func NewListHead(i interface{}) interface{} {
	list_head := &ListHead{}
	list_head.Prev = list_head
	list_head.Next = list_head
}

func Register(l *ListHead, i interface{}) {
	val := reflect.ValueOf(i)
	var val_p, f_p reflect.Value
	if val.Kind() == reflect.Ptr {
		val_p = val
	} else if val.CanAddr() {
		val_p = val.Addr()
	}
	f_p = val_p.Elem().FieldByName("ListHead")
	diff_type[val_p.Type().Name()] = false
	if f_p.Kind() != reflect.Ptr && f_p.CanAddr() {
		f_p = f_p.Addr()
		diff_type[val_p.Type().Name()] = true
	}
	diff[val_p] = uintptr(f_p.Pointer() - uintptr(val_p.Pointer()))
}

// ContainOf(l *ListHead, i interface{}) interface{} {

func ContainOf(l *ListHead, i interface{}) interface{} {

	val := reflect.ValueOf(i)
	var val_p, f_p reflect.Value
	if val.Kind() == reflect.Ptr {
		val_p = val
	} else if val.CanAddr() {
		val_p = val.Addr()
	}
	f_p = val_p.Elem().FieldByName("ListHead")
	if f_p.Kind() != reflect.Ptr && f_p.CanAddr() {
		f_p = f_p.Addr()
	}

	fmt.Printf("%x %x\n", f_p.Pointer(), val_p.Pointer())
	fmt.Println(val_p.Elem().Type())
	ret_p := uintptr(reflect.ValueOf(l).Pointer()) - (uintptr(f_p.Pointer() - uintptr(val_p.Pointer())))
	fmt.Printf("arg=%x ret=%x ", uintptr(reflect.ValueOf(l).Pointer()), ret_p)
	return reflect.NewAt(val_p.Elem().Type(), unsafe.Pointer(ret_p)).Interface()
}
