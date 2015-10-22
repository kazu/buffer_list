# buffer_list
go Package  double linked list with seuqnecial buffer

# example

type Hoge struct {
  a int64
  b int64
}

buffer_list := buffer_list.New(&Hoge{})
hoge := (*Hoge)(buffer_list.Value())

hoge.a = 200
hoge.b = 500


