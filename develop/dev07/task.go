package main

import (
	"fmt"
	"time"
)

/*
=== Or channel ===

Реализовать функцию, которая будет объединять один или более done каналов в single канал если один из его составляющих каналов закроется.
Одним из вариантов было бы очевидно написать выражение при помощи select, которое бы реализовывало эту связь,
однако иногда неизестно общее число done каналов, с которыми вы работаете в рантайме.
В этом случае удобнее использовать вызов единственной функции, которая, приняв на вход один или более or каналов, реализовывала весь функционал.

Определение функции:
var or func(channels ...<- chan interface{}) <- chan interface{}

Пример использования функции:
sig := func(after time.Duration) <- chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(after)
}()
return c
}

start := time.Now()
<-or (
	sig(2*time.Hour),
	sig(5*time.Minute),
	sig(1*time.Second),
	sig(1*time.Hour),
	sig(1*time.Minute),
)

fmt.Printf(“fone after %v”, time.Since(start))
*/

// or объединяет один или более done-каналов в single-канал, который закроется,
// если хотя бы один из его составляющих каналов закроется.
func or(channels ...<-chan any) <-chan any {
	// Условия для выхода из рекурсии.
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}

	// Создаем single-канал done, который будет закрыт при закрытии составляющих каналов.
	done := make(chan any)

	// Если каналов >= 2, то создаем горутину с селектом.
	// Она будет читать каналы с индексами 0 и 1, а также or для каналов с индексами >= 2.
	go func() {
		defer close(done)
		select {
		case <-channels[0]:
		case <-channels[1]:
		case <-or(append(channels[2:], done)...):
			// Рекурсивно вызываем or для каналов с индексами >= 2.
			// Добавляем в аргументы done-канал, чтобы вызовы завершались при его закрытии.
		}
	}()

	// Возвращаем single-канал done.
	return done
}

func main() {
	sig := func(after time.Duration) <-chan any {
		c := make(chan any)
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)

	fmt.Printf("done after %v\n", time.Since(start))
}
