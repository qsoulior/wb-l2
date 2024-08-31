package pattern

import (
	"fmt"
)

/*
Реализовать паттерн «посетитель».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
https://en.wikipedia.org/wiki/Visitor_pattern
*/

// Применимость: когда нужно добавить новую функциональность, изменяя исходные структуры минимально.
// Плюсы: упрощает добавление новой функциональности с минимальным изменением исходных структур, упрощает выполнение операций над полиморфными объектами.
// Минусы: для корректной реализации требует метод accept в исходных структурах, сложно применим к изменяющейся иерархии структур.
// Реальный пример: добавление функциональности определения итоговой стоимости разнородных товаров и услуг.

// Интерфейс элемента, принимающего посетителя.
type VisitedElement interface {
	Accept(v Visitor)
}

// Структура конкретного элемента, реализующая интерфейс.
type ConcreteVisitedElement1 struct {
}

// Принимает посетителя.
func (e *ConcreteVisitedElement1) Accept(v Visitor) { v.VisitConcreteElement1(e) }

// Структура конкретного элемента, реализующая интерфейс.
type ConcreteVisitedElement2 struct {
}

// Принимает посетителя.
func (e *ConcreteVisitedElement2) Accept(v Visitor) { v.VisitConcreteElement2(e) }

// Интерфейс посетителя.
type Visitor interface {
	VisitConcreteElement1(e *ConcreteVisitedElement1)
	VisitConcreteElement2(e *ConcreteVisitedElement2)
}

// Структура конкретного посетителя, реализующая интерфейс.
type ConcreteVisitor1 struct {
}

// Посещает элемент №1 и выполняет какое-то действие с ним.
func (v ConcreteVisitor1) VisitConcreteElement1(e *ConcreteVisitedElement1) {
	fmt.Printf("ConcreteVisitor1.VisitConcreteElement1: %v\n", e)
}

// Посещает элемент №2 и выполняет какое-то действие с ним.
func (v ConcreteVisitor1) VisitConcreteElement2(e *ConcreteVisitedElement2) {
	fmt.Printf("ConcreteVisitor1.VisitConcreteElement2: %v\n", e)
}

// Структура конкретного посетителя, реализующая интерфейс.
type ConcreteVisitor2 struct {
}

// Посещает элемент №1 и выполняет какое-то действие с ним.
func (v ConcreteVisitor2) VisitConcreteElement1(e *ConcreteVisitedElement1) {
	fmt.Printf("ConcreteVisitor2.VisitConcreteElement1: %v\n", e)
}

// Посещает элемент №2 и выполняет какое-то действие с ним.
func (v ConcreteVisitor2) VisitConcreteElement2(e *ConcreteVisitedElement2) {
	fmt.Printf("ConcreteVisitor2.VisitConcreteElement2: %v\n", e)
}

/*
visitor1 := ConcreteVisitor1{}
visitor2 := ConcreteVisitor2{}

elements := []VisitedElement{&ConcreteVisitedElement1{}, &ConcreteVisitedElement2{}}
for _, element := range elements {
	element.Accept(visitor1)
	element.Accept(visitor2)
}
*/
