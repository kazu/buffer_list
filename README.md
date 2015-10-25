# buffer_list
go Package  double linked list with (slice like) sequencial buffer

# example
```go
type Hoge struct {
  a int64
  b int64
}

buffer_list := buffer_list.New(&Hoge{}, 100)
buffer_list.SetCastFunc(func(p interface{}) interface{} {
  return (*Hoge)(p)
}

hoge := buffer_list.ValueWithCast()
hoge.a = 200
hoge.b = 500


cur := blist.GetElement()
new_e := blist.InsertNewElem(cur)
hoge2 := new_e.ValueWithCast
hoge2.a = 222
hoge2.b = 2222

for e:= blist.Front(); e != nil; e.Next() {
  fmt.Println("value", e.ValueWithCast() )
}

```
