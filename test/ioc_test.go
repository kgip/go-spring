package test

import (
	"fmt"
	_ "go-spring/test/a"
	_ "go-spring/test/b"
	"testing"
)

func TestInject(t *testing.T) {
	fmt.Println("开始运行")
}
