package main

import (
	"fmt"
	"os"
	"time"

	"github.com/beevik/ntp"
)

/*
=== Базовая задача ===

Создать программу печатающую точное время с использованием NTP библиотеки. Инициализировать как go module.
Использовать библиотеку https://github.com/beevik/ntp.
Написать программу печатающую текущее время / точное время с использованием этой библиотеки.

Программа должна быть оформлена с использованием как go module.
Программа должна корректно обрабатывать ошибки библиотеки: распечатывать их в STDERR и возвращать ненулевой код выхода в OS.
Программа должна проходить проверки go vet и golint.
*/

func main() {
	// Получаем время по NTP.
	networkTime, err := ntp.Time("0.ru.pool.ntp.org")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get network time: %s\n", err)
		os.Exit(1)
	}

	// Получаем локальное время.
	localTime := time.Now()

	// Выводим текущее локальное время и точное время по NTP.
	fmt.Printf("local time:\t%s\n", localTime)
	fmt.Printf("network time:\t%s\n", networkTime)
}
