package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

/*
=== Утилита cut ===

Реализовать утилиту аналог консольной команды cut (man cut). Утилита должна принимать строки через STDIN, разбивать по разделителю (TAB) на колонки и выводить запрошенные.

Реализовать поддержку утилитой следующих ключей:
-f - "fields" - выбрать поля (колонки)
-d - "delimiter" - использовать другой разделитель
-s - "separated" - только строки с разделителем
*/

// Cut разбивает строку по delim и возвращает новую строку с полями, индексы которых указаны в fields.
// Если separated равен true, то Cut возвращает false вторым значением для строки, которая не содержит delim.
func Cut(s string, fields []int, delim string, separated bool) (string, bool) {
	parts := strings.Split(s, delim) // части строки, полученные по delim

	// Если строка не пустая и не содержит delim.
	if len(parts) == 1 {
		if separated {
			return "", false
		}
		return parts[0], true
	}

	// Если строка содержит delim, и fields не пустой.
	if len(fields) >= 1 {
		newParts := make([]string, 0)
		// Выбираем только те колонки, которые попали в fields.
		for _, field := range fields {
			if field < len(parts) {
				newParts = append(newParts, parts[field])
			}
		}
		// Обновляем части строки.
		parts = newParts
	}

	// Соединяем части строки
	return strings.Join(parts, delim), true
}

var (
	fields    Fields
	delimiter string
	separated bool
)

// Тип полей для парсинга пакетом flag.
type Fields []int

// String возвращает строкове представление переменной типа Fields.
func (f *Fields) String() string { return fmt.Sprintf("%v", *f) }

// Set устанавливает значение переменной типа Fields из строки s.
func (f *Fields) Set(s string) error {
	// Делим строку s по разделителю "," и добавляем индексы полей в срез.
	for _, value := range strings.Split(s, ",") {
		field, err := strconv.Atoi(value)
		if err != nil || field-1 < 0 {
			return errors.New("fields are numbered from 1")
		}
		*f = append(*f, field-1)
	}

	// Сортируем индексы полей по возрастанию.
	slices.Sort(*f)
	return nil
}

// Errorf пишет данные в stderr.
func Errorf(format string, a ...any) { fmt.Fprintf(os.Stderr, format, a...) }

func init() {
	// Настраиваем flag на парсинг параметров cut.
	flag.Var(&fields, "f", "select only these fields; also print any line that contains no delimiter character, unless the -s option is specified")
	flag.StringVar(&delimiter, "d", "\t", "use DELIM instead of TAB for field delimiter")
	flag.BoolVar(&separated, "s", false, "do not print lines not containing delimiters")
}

func main() {
	flag.Parse()

	// Проверяем наличие флага -f.
	if len(fields) == 0 {
		Errorf("you must specify fields\n")
		flag.PrintDefaults()
		os.Exit(2)
	}

	// Сканируем строки и для каждой из них возвращаем результат операции.
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		result, ok := Cut(scanner.Text(), fields, delimiter, separated)
		if ok {
			fmt.Println(result)
		}
	}

	// Возвращаем ошибку, если не удалось отсканировать строки.
	if err := scanner.Err(); err != nil {
		Errorf("failed to read file: %s\n", err)
		os.Exit(1)
	}
}
