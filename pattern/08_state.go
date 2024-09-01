package pattern

import "fmt"

/*
Реализовать паттерн «состояние».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
https://en.wikipedia.org/wiki/State_pattern
*/

// Применимость: когда нужно менять поведение в зависимости от некоторых сменяющих друг друга состояний.
// Плюсы: решает проблему множества условных операторов при большом количестве состояний, инкапсулирует логику действия для каждого состояния.
// Минусы: может привести к оверинжинирингу, если состояний мало или они редко сменяются.
// Реальный пример: изменения состояний и поведения сущности заказа (к примеру, заказ может быть отменен только, если оформлен).

// Интерфейс состояния.
type State interface {
	SetContext(context *StateContext)
	Action()
}

// Структура конкретного состояния, реализующая интерфейс.
type ConcreteState1 struct {
	context *StateContext
}

// Устанавливает контекст, который хранит состояние.
func (s *ConcreteState1) SetContext(context *StateContext) { s.context = context }

// Выполняет действие и меняет состояние у контекста.
func (s ConcreteState1) Action() {
	fmt.Println("ConcreteState1.Action")
	s.context.SetState(&ConcreteState2{})
}

// Структура конкретного состояния, реализующая интерфейс.
type ConcreteState2 struct {
	context *StateContext
}

// Устанавливает контекст, который хранит состояние.
func (s *ConcreteState2) SetContext(context *StateContext) { s.context = context }

// Выполняет действие и меняет состояние у контекста.
func (s ConcreteState2) Action() {
	fmt.Println("ConcreteState2.Action")
	s.context.SetState(&ConcreteState3{})
}

// Структура конкретного состояния, реализующая интерфейс.
type ConcreteState3 struct {
	context *StateContext
}

// Устанавливает контекст, который хранит состояние.
func (s *ConcreteState3) SetContext(context *StateContext) { s.context = context }

// Выполняет действие и меняет состояние у контекста.
func (s ConcreteState3) Action() {
	fmt.Println("ConcreteState3.Action")
	s.context.SetState(&ConcreteState1{})
}

// Структура контекста состояния.
type StateContext struct {
	state State
}

// Конструктор контекста состояния.
func NewStateContext(initialState State) *StateContext {
	ctx := &StateContext{initialState}
	ctx.state.SetContext(ctx)
	return ctx
}

// Устанавливает/меняет состояние в контексте.
func (c *StateContext) SetState(state State) {
	c.state = state
	c.state.SetContext(c)
}

// Выполняет некоторое действие в зависимости от установленного состояния.
func (c StateContext) Action() { c.state.Action() }

/*
ctx := NewStateContext(&ConcreteState1{})
ctx.Action() // ConcreteState1.Action
ctx.Action() // ConcreteState2.Action
ctx.Action() // ConcreteState3.Action
ctx.Action() // ConcreteState1.Action
*/
