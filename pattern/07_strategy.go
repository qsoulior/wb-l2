package pattern

import "fmt"

/*
Реализовать паттерн «стратегия».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
https://en.wikipedia.org/wiki/Strategy_pattern
*/

// Применимость: когда нужно выполнять схожие действия из одного семества с возможностью взаимной замены в процессе исполнения программы.
// Плюсы: инкапсулирует логику каждого действия от других действий, позволяет менять действия во время исполнения программы.
// Минусы: требует введения вспомогательных структур и интерфейсов, усложняет код программы.
// Реальный пример: замена стратегии хранения данных в реальном времени в зависимости от поступающих событий (допустим, замена кэша на постоянное хранилище при необходимости резервирования данных).

// Интерфейс стратегии, вызывающей некоторое действие.
type Strategy interface {
	Execute()
}

// Структура конкретной стратегии, реализующая интерфейс.
type ConcreteStrategy1 struct {
}

// Выполнить некоторое действие.
func (s ConcreteStrategy1) Execute() { fmt.Println("ConcreteStrategy1.Execute") }

// Структура конкретной стратегии, реализующая интерфейс.
type ConcreteStrategy2 struct {
}

// Выполнить некоторое действие.
func (s ConcreteStrategy2) Execute() { fmt.Println("ConcreteStrategy2.Execute") }

// Структура контекста стратегии.
type StrategyContext struct {
	strategy Strategy
}

// Установить активную стратегию.
func (c *StrategyContext) SetStrategy(strategy Strategy) { c.strategy = strategy }

// Вызвать действие активной стратегии.
func (c StrategyContext) Action() { c.strategy.Execute() }

/*
ctx := StrategyContext{}

ctx.SetStrategy(ConcreteStrategy1{})
ctx.Action() // ConcreteStrategy1.Execute

ctx.SetStrategy(ConcreteStrategy2{})
ctx.Action() // ConcreteStrategy2.Execute
*/
