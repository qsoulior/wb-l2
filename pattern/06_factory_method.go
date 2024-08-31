package pattern

import "fmt"

/*
Реализовать паттерн «фабричный метод».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
https://en.wikipedia.org/wiki/Factory_method_pattern
*/

// Применимость: когда нужно реализовать настраиваемое создание различных объектов на основе общего интерфейса.
// Плюсы: упрощает создание объектов в сложной системе, инкапсулирует логику порождения объектов и упрощает поддерживаемость этой логики.
// Минусы: для каждой структуры требует создания своего фабричного метода (и структуры с методом, соотвественно), что увеличивает количество структур и их связанность.
// Реальный пример: порождение различных шаблонов писем (регистрация, новости, безопасность), рассылаемых по электронной почте.

// Интерфейс объектов, которые создаются фабричным методом.
type FactoryProduct interface {
	Action()
}

// Структура конкретного объекта, реализующая интерфейс.
type ConcreteFactoryProduct1 struct {
}

// Выполняет некоторое действие объекта.
func (p ConcreteFactoryProduct1) Action() { fmt.Println("ConcreteFactoryProduct1.Action") }

// Структура конкретного объекта, реализующая интерфейс.
type ConcreteFactoryProduct2 struct {
}

// Выполняет некоторое действие объекта.
func (p ConcreteFactoryProduct2) Action() { fmt.Println("ConcreteFactoryProduct2.Action") }

// Интерфейс с абстрактным фабричным методом.
type Factory interface {
	CreateProduct() FactoryProduct
}

// Структура с конкретным фабричным методом, реализующая интерфейс.
type ConcreteFactory1 struct {
}

// Создает и возвращает конкретный объект.
func (f ConcreteFactory1) CreateProduct() FactoryProduct { return ConcreteFactoryProduct1{} }

// Структура с конкретным фабричным методом, реализующая интерфейс.
type ConcreteFactory2 struct {
}

// Создает и возвращает конкретный объект.
func (f ConcreteFactory2) CreateProduct() FactoryProduct { return ConcreteFactoryProduct2{} }

/*
var factory Factory

factory = ConcreteFactory1{}
product1 := factory.CreateProduct()
product1.Action() // ConcreteFactoryProduct1.Action

factory = ConcreteFactory2{}
product2 := factory.CreateProduct()
product2.Action() // ConcreteFactoryProduct2.Action
*/
