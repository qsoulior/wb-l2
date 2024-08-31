package pattern

import "fmt"

/*
Реализовать паттерн «фасад».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
https://en.wikipedia.org/wiki/Facade_pattern
*/

// Применимость: когда нужно предоставить упрощенный интерфейс к некоторой сложной системе.
// Плюсы: обособляет и инкапсулирует сложную систему от пользователя, предоставляя достаточное взаимодействие.
// Минусы: может стать слишком крупным, что усложнит его использование и тестирование.
// Реальный пример: сервис заказов, использующий множество других подсервисов (товаров, оплаты и доставки).

// Структура сервиса сложной системы.
type FacadeService1 struct {
}

// Выполняет некоторое действие.
func (s FacadeService1) Action() { fmt.Println("Service1.Action") }

// Структура сервиса сложной системы.
type FacadeService2 struct {
}

// Выполняет некоторое действие.
func (s FacadeService2) Action() { fmt.Println("Service2.Action") }

// Структура сервиса сложной системы.
type FacadeService3 struct {
}

// Выполняет некоторое действие.
func (s FacadeService3) Action() { fmt.Println("Service3.Action") }

// Структура фасада сложной системы.
type Facade struct {
	service1 FacadeService1
	service2 FacadeService2
	service3 FacadeService3
}

// Предоставляет некоторую функциональность системы.
func (s Facade) ComplexAction() {
	s.service1.Action()
	s.service2.Action()
	s.service3.Action()
}

/*
service1 := FacadeService1{}
service2 := FacadeService2{}
service3 := FacadeService3{}

facade := Facade{service1: service1, service2: service2, service3: service3}
facade.ComplexAction()
*/
