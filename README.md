# buffer_list
go Package  double linked list with (slice like) sequencial buffer

# example
```go
type Hoge struct {
  a int64
  b int64
  c *int64
}

buffer_list := buffer_list.New(&Hoge{}, 100)

hoge := buffer_list.Value().(*Hoge)
hoge.a = 200
hoge.b = 500


cur := blist.GetElement()
new_e := blist.InsertNewElem(cur)
hoge2 := new_e.Value().(*Hoge)
hoge2.a = 222
hoge2.b = 2222
hoge2.c = 1234
new_e.Commit() // protect c from gc Free

for e:= blist.Front(); e != nil; e.Next() {
  fmt.Println("value", e.Value() )
}

```
