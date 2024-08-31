package pattern

import "fmt"

/*
Реализовать паттерн «команда».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
https://en.wikipedia.org/wiki/Command_pattern
*/

// Применимость: когда нужно выполнять действия в любой момент времени без прямой зависимости между исполнителями и инициаторами.
// Плюсы: снижает зависимость между исполнителями и инициаторами действий, позволяет реализовать отложенное выполнение действий, их повторение и отмену.
// Минусы: требует введения вспомогательных структур и интерфейсов, усложняет код программы.
// Реальный пример: любое приложение, выполняющее различные действия и имеющее несколько способов запустить выполнение этих действий (через GUI, CLI или как-то иначе).

// Интерфейс команды.
type Command interface {
	Execute()
}

// Структура конкретной команды, реализующая интерфейс.
type ConcreteCommand1 struct {
	receiver CommandReceiver // получатель, который выполнит действие
}

// Выполняет команду.
func (c ConcreteCommand1) Execute() { c.receiver.Action1() }

// Структура конкретной команды, реализующей интерфейс.
type ConcreteCommand2 struct {
	receiver CommandReceiver // получатель, который выполнит действие
}

// Выполняет команду.
func (c ConcreteCommand2) Execute() { c.receiver.Action2() }

// Структура инициатора команды.
type CommandInvoker struct {
	command Command // команда, которая будет запущена
}

// Устанавливает команду для запуска.
func (i *CommandInvoker) SetCommand(cmd Command) { i.command = cmd }

// Запускает команду.
func (i CommandInvoker) Invoke() { i.command.Execute() }

// Интерфейс получателя и исполнителя команды.
type CommandReceiver interface {
	Action1()
	Action2()
}

// Структура конкретного получателя команды, реализующая интерфейс.
type ConcreteCommandReceiver struct {
}

// Выполняет некоторое действие.
func (r ConcreteCommandReceiver) Action1() { fmt.Println("Action1") }

// Выполняет некоторое действие.
func (r ConcreteCommandReceiver) Action2() { fmt.Println("Action2") }

/*
receiver := ConcreteReceiver{}

cmd1 := ConcreteCommand1{receiver: receiver}
cmd2 := ConcreteCommand2{receiver: receiver}

invoker := Invoker{}
invoker.SetCommand(cmd1)
invoker.Invoke() // Action1
invoker.SetCommand(cmd2)
invoker.Invoke() // Action2
*/
