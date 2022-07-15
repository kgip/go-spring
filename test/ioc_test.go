package test

import (
	"fmt"
	_ "github.com/kgip/go-spring/test/a"
	_ "github.com/kgip/go-spring/test/b"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestInject(t *testing.T) {
	fmt.Println("开始运行")
}

type Reader interface {
	Read()
}

type Writer interface {
	Write()
}

type Book struct {
	age    int
	title  string
	author string
}

func (Book) Read() {
	fmt.Println("读")
}

func (Book) Write() {
	fmt.Println("写")
}

func TestStructName(t *testing.T) {
	book := Book{}
	rv := reflect.ValueOf(&book)
	t.Log(rv.Type().Elem().Name())
}

func TestSort(t *testing.T) {
	books := []Book{{
		age: 1,
	}, {age: 2}, {age: 0}}
	sort.Slice(books, func(i, j int) bool {
		return books[i].age > books[j].age
	})
	t.Log(books)
}

func TestCreateStruct(t *testing.T) {
	rt := reflect.TypeOf(Book{})
	value := reflect.New(rt).Interface()
	t.Log(value)
}

func TestCreateArray(t *testing.T) {
	arr := [3]int{1, 2, 4}
	rv := reflect.ValueOf(arr)
	t.Log(rv.Kind())
	instance := reflect.New(rv.Type()).Elem().Interface()

	t.Log(instance)
}

func R(book Reader) {
	rv := reflect.ValueOf(book)
	fmt.Println(rv.Elem().Kind())
}

func TestInterface(t *testing.T) {
	var reader Reader = &Book{}
	reader.Read()
	writer, ok := reader.(Writer)
	if ok {
		writer.Write()
	}
}

type User struct {
	Id   int
	Name string
}

//type Mapper[T any] struct{}
//
//func (Mapper[T]) List() []*T {
//	a := new(T)
//	rv := reflect.TypeOf(a).Elem()
//	for i := 0; i < rv.NumField(); i++ {
//		fmt.Println(rv.Field(i).Name)
//	}
//	return []*T{a, a}
//}
//
//type UserService struct {
//	UserMapper Mapper[User]
//	Booker     *Book
//}

func TestName(t *testing.T) {
	//userService := UserService{}
	//rv := reflect.ValueOf(&userService)
	//t.Log(rv.Elem().Field(0).CanSet())
	//frv := reflect.New(rv.Elem().Field(0).Type())

	//var mapper = &Mapper[User]{}
	//prv := reflect.ValueOf(mapper)
	//rt := prv.Type()
	//value := reflect.New(rt.Elem())
	//t.Log(prv.Elem().CanAddr())
	//t.Log(prv.CanAddr())
	//
	//rv.Elem().Field(0).Set(frv.Elem())
	//t.Log(userService.UserMapper.List())
	//if userMapper, ok := value.(*Mapper[User]); ok {
	//	list := userMapper.List()
	//	for _, user := range list {
	//		t.Log((*user).Id)
	//	}
	//}
}

func TestReflect(t *testing.T) {
	a := 1
	rv := reflect.ValueOf(&a)
	t.Log(rv.Elem().CanAddr())
}

func TestPtr(t *testing.T) {
	var a = 1
	var ptr interface{} = &a
	rv := reflect.ValueOf(ptr)
	t.Log(rv.Elem().Interface())
	//var ptr1 = &ptr
	//var ptr2 = &ptr1
	//var rv = reflect.ValueOf(ptr2)
	//t.Log(rv.Elem().Kind())
	//t.Log(rv.Elem().Elem().Kind())
	//t.Log(rv.Elem().Elem().Elem().Kind())
}

type A struct {
	b    *B
	name string
}

type B struct {
	a    *A
	name string
}

func TestCircularReference(t *testing.T) {
	var a *A
	var b *B
	a = &A{name: "AAAAA"}
	b = &B{name: "BBBBB"}
	b.a = a
	a.b = b
	t.Log(a, b)
}

func TestB(t *testing.T) {
	var splits []string
	splits = append(splits, "aa")
	t.Log(splits)
	str := "value=path"
	index := strings.Index(str, "=")
	t.Log(str[:index])
	t.Log(str[index+1:])
}

type C struct {
	int
}

func TestAnon(t *testing.T) {
	c := C{}
	rt := reflect.TypeOf(c)
	if rt.Field(0).Anonymous {
		t.Log(rt.Field(0).Name)
		t.Log(rt.Field(0).Type.Name())
	}
}
