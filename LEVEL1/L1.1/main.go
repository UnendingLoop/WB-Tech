package main

import "fmt"

/*Дана структура Human (с произвольным набором полей и методов).
Реализовать встраивание методов в структуре Action от родительской структуры Human (аналог наследования).
Подсказка: используйте композицию (embedded struct), чтобы Action имел все методы Human.*/

// Human - родительская структура
type Human struct {
	Name string
	Age  int
}

// Action - дочерняя структура с методами Human
type Action struct {
	Human
}

func (h *Human) assignName() {
	h.Name = "Jesus Godson"
}
func (h *Human) assignAge() {
	h.Age = 33
}

func (a Action) sayName() {
	fmt.Println("My name is", a.Name)
}
func (a Action) sayAge() {
	fmt.Println("My age is", a.Age)
}

func main() {
	var person Human
	person.assignName()
	person.assignAge()

	var action Action
	action.assignName()
	action.assignAge()
	action.sayName()
	action.sayAge()
}
