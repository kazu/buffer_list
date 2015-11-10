# buffer_list [![wercker status](https://app.wercker.com/status/af71821c3a51e35a170766fdab30e1b8/s "wercker status")](https://app.wercker.com/project/bykey/af71821c3a51e35a170766fdab30e1b8)
Package is double linked list with (slice like) sequencial buffer.

container/list's value is not allocated sequncially. so element value occur fragmentation.
buffer_list has []byte buffer in List.


## buffer_list.List

buffer_list.List is double linked list with member data buffer.
if you want to have buffer of list's value, you cant get effect its.


```go
type Hoge struct {
  a int64
  b int64
  c *int64
}

buffer_list := buffer_list.New(&Hoge{}, 100)

hoge := buffer_list.Front().Value().(*Hoge)
hoge.a = 200
hoge.b = 500


new_e := blist.InsertLast() // allocate new element/value
hoge2 := new_e.Value().(*Hoge)
hoge2.a = 222
hoge2.b = 2222
hoge2.c = 1234
new_e.Commit() // protect from gc Free

for e:= blist.Front(); e != nil; e.Next() {
  fmt.Println("value", e.Value() )
}

```

## buffer_list.AList

AList is 'alias list' of buffer_list.List. AList permit you to have multi list with shared buffer data. 
Alist has Genetator pattern. 

```go
bl :=  buffer_list.New(Hoge{}, 10000)
alist := Alist{parent: bl}
ae := alist.NewElem()
alist.Front().Insert(ae)
v := ae.Value().(*Hoge)
v.Commit()

// iteration
for e := range alist.Generator() {
     v := e.Value().(*Hoge)
}

```

## WARNING
this packages buffer is []byte so pointer member of nested struct is commonly freed by GC.
buffer_list implement to protect struct's  pointer member. but this problem is not fully fixed.


if struct member is chan/slice/embedded pointer interface/Array/map/function/struct , protected.


## Contributing

Bug reports and pull requests are welcome on GitHub at https://github.com/kazu/buffer_list
