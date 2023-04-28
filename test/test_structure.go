package test

import (
	"fmt"
	"testing"
)

type class_school struct {
	class  int
	school string
}

type student struct {
	name string
	age  int
	//c_s class_school
	//匿名嵌套
	class_school
}

type teacher struct {
	name   string
	age    int
	course string //老师讲授课程
	c_s    class_school
}

func TestStructure(t *testing.T) {
	s1 := student{
		name: "心安",
		age:  5,
		class_school: class_school{
			class:  8,
			school: "叮当幼儿园",
		},
	}
	t1 := teacher{
		name:   "花花老师",
		age:    27,
		course: "美术课",
		c_s: class_school{
			class:  4,
			school: "叮当幼儿园",
		},
	}
	//Go语言语法糖
	fmt.Printf("t1: %#v", t1)
	fmt.Println(t1.c_s.school)
	fmt.Printf("s1: %#v", s1)
	fmt.Println(s1.school)
}
