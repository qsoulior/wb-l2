package pattern

import "fmt"

/*
Реализовать паттерн «цепочка вызовов».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
https://en.wikipedia.org/wiki/Chain-of-responsibility_pattern
*/

// Применимость: когда нужно обработать какой-то вызов или запрос в последовательной цепочке операций.
// Плюсы: снижает зависимость между операциями обработки, обособляет ответственность каждой из операций.
// Минусы: при неправильной реализации цепочки или ошибке в обработчике вызов может быть не обработан.
// Реальный пример: обработка HTTP-запроса в промежуточных обработчиках (middleware) перед выполнением операций целевым обработчиком.

// Структура вызова.
type Call struct {
	x, y, z bool
}

// Интерфейс обработчика вызова.
type CallHandler interface {
	Handle(call *Call)
}

// Структура конкретного обработчика вызова, реализующая интерфейс.
type ConcreteCallHandler1 struct {
	Next CallHandler
}

// Обрабатывает вызов.
func (h *ConcreteCallHandler1) Handle(call *Call) {
	call.x = true
	fmt.Printf("ConcreteCallHandler1.Handle: %v\n", call)
	h.Next.Handle(call)
}

// Структура конкретного обработчика вызова, реализующая интерфейс.
type ConcreteCallHandler2 struct {
	Next CallHandler
}

// Обрабатывает вызов, прерывает при необходимости.
func (h *ConcreteCallHandler2) Handle(call *Call) {
	if !call.x {
		return
	}

	call.y = true
	fmt.Printf("ConcreteCallHandler2.Handle: %v\n", call)
	h.Next.Handle(call)
}

// Структура конкретного обработчика вызова, реализующая интерфейс.
type ConcreteCallHandler3 struct {
}

// Обрабатывает вызов, прерывает при необходимости.
func (h *ConcreteCallHandler3) Handle(call *Call) {
	if !call.y {
		return
	}

	call.z = true
	fmt.Printf("ConcreteCallHandler3.Handle: %v\n", call)
}

/*
handler3 := &ConcreteCallHandler3{}
handler2 := &ConcreteCallHandler2{Next: handler3}
handler1 := &ConcreteCallHandler1{Next: handler2}

call := Call{}
handler1.Handle(&call)
*/
