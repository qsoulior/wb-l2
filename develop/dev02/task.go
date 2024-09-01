package main

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

/*
=== Задача на распаковку ===

Создать Go функцию, осуществляющую примитивную распаковку строки, содержащую повторяющиеся символы / руны, например:
	- "a4bc2d5e" => "aaaabccddddde"
	- "abcd" => "abcd"
	- "45" => "" (некорректная строка)
	- "" => ""
Дополнительное задание: поддержка escape - последовательностей
	- qwe\4\5 => qwe45 (*)
	- qwe\45 => qwe44444 (*)
	- qwe\\5 => qwe\\\\\ (*)

В случае если была передана некорректная строка функция должна возвращать ошибку. Написать unit-тесты.

Функция должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

var ErrStringInvalid = errors.New("invalid string")

// Unpack осуществляет распаковку строки и возвращает распакованную строку.
func Unpack(inp string) (string, error) {
	var (
		out         strings.Builder
		prev        rune
		prevEscaped bool
	)

	for _, r := range inp {
		if prev == '\\' && !prevEscaped { // текущий символ был экранирован, а предыдущий - нет
			out.WriteRune(r)
			prevEscaped = true
		} else {
			if prev != 0 && unicode.IsDigit(r) { // цифра, расположенная не в начале
				if unicode.IsDigit(prev) && !prevEscaped { // предыдущий символ - неэкранированная цифра
					return "", ErrStringInvalid
				}

				// Записываем предыдущий символ n-1 раз, где n - текущая цифра.
				for i := 0; i < int(r-'0')-1; i++ {
					out.WriteRune(prev)
				}
			} else if r != '\\' { // любой символ, кроме обратного слэша
				out.WriteRune(r)
			}
			prevEscaped = false
		}
		prev = r
	}

	fmt.Println()
	return out.String(), nil
}
