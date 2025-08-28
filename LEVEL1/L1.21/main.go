package main

import (
	"fmt"
	"strings"
)

/*
Реализовать паттерн проектирования «Адаптер» на любом примере.
Описание: паттерн Adapter позволяет сконвертировать интерфейс
одного класса в интерфейс другого, который ожидает клиент.

Продемонстрируйте на простом примере в Go: у вас есть существующий
интерфейс (или структура) и другой, несовместимый по интерфейсу потребитель —
напишите адаптер, который реализует нужный интерфейс и делегирует вызовы
к встроенному объекту.

Поясните применимость паттерна, его плюсы и минусы, а также
приведите реальные примеры использования.
*/

// НЕОПДДЕРЖИВАЕМАЯ клиентом устарелая структура
type deprecatedStruct struct {
	name string //name+surname
	age  uint
	sex  string //"male"/"female"
}

// ВНУТРЕННИЙ интерфейс, ожидаемый клиентом
type internalInterface interface {
	getName() string
	getSurname() string
	getAge() int
	getAppeal() string
}

func newInternalFromForeign(f *deprecatedStruct) internalInterface {
	return &adapterStruct{adaptor: f}
}

// ВНУТРЕННЯЯ структура
type adapterStruct struct {
	adaptor *deprecatedStruct
}

func (p *adapterStruct) getName() string {
	parts := strings.Split(p.adaptor.name, " ")
	return parts[0]
}
func (p *adapterStruct) getSurname() string {
	parts := strings.Split(p.adaptor.name, " ")
	if len(parts) >= 2 {
		return parts[1]
	}
	return ""
}
func (p *adapterStruct) getAge() int {
	return int(p.adaptor.age)
}
func (p *adapterStruct) getAppeal() string {
	appeal := "Dear"
	switch p.adaptor.sex {
	case "male":
		appeal = "Mister"
	case "female":
		appeal = "Mis"
	}
	return appeal
}

// Клиентский код, поддерживает только internalInterface
func printInfo(p internalInterface) {
	fmt.Printf("Customer's info is: %s %s %s, their age is %d.", p.getAppeal(), p.getName(), p.getSurname(), p.getAge())
}

func main() {
	//экземпляр старой структуры
	oldPerson := deprecatedStruct{name: "John Doe", age: 32, sex: "male"}
	//засовываем структуру в адаптер
	person := newInternalFromForeign(&oldPerson)

	//клиент получает доступ к методам internalInterface не зная про адаптер
	printInfo(person)
}
