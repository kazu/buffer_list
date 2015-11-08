# buffer_list [![wercker status](https://app.wercker.com/status/af71821c3a51e35a170766fdab30e1b8/s "wercker status")](https://app.wercker.com/project/bykey/af71821c3a51e35a170766fdab30e1b8)
go Package  double linked list with (slice like) sequencial buffer


# WARNING
Element.Value() is protected by GC Free in common case.
but following case is not protected


```gp
type IData struct {
        a int
        b int
}

type IValue struct {
        a *IData
}
type TValue struct {
        a  int
        iv IValue
        b  *IData
}

l := buffer_list.New(&TValue{}, 100)
e := l.InsertLast()
v := e.Value().(*TValue)
e.InitValue()

//v.iv is not protected. it may be freed 
v.iv = IValue{a: &IData{a: 10, b: 1}}

//v.b is OK
v.b = &IData{a: i * 10, b: i + 1}

e.Commit()

```


# example
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


new_e := blist.InsertLast()
hoge2 := new_e.Value().(*Hoge)
hoge2.a = 222
hoge2.b = 2222
hoge2.c = 1234
new_e.Commit() // protect c from gc Free

for e:= blist.Front(); e != nil; e.Next() {
  fmt.Println("value", e.Value() )
}

```
