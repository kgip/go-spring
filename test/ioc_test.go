package test

import (
	"fmt"
	_ "go-spring/test/a"
	_ "go-spring/test/b"
	"reflect"
	"sort"
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
