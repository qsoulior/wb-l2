package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

/*
=== Утилита sort ===

Отсортировать строки в файле по аналогии с консольной утилитой sort (man sort — смотрим описание и основные параметры): на входе подается файл из несортированными строками, на выходе — файл с отсортированными.

Реализовать поддержку утилитой следующих ключей:
	-k — указание колонки для сортировки (слова в строке могут выступать в качестве колонок, по умолчанию разделитель — пробел)
	-n — сортировать по числовому значению
	-r — сортировать в обратном порядке
	-u — не выводить повторяющиеся строки

Дополнительно:
Реализовать поддержку утилитой следующих ключей:
	-M — сортировать по названию месяца
	-b — игнорировать хвостовые пробелы
	-c — проверять отсортированы ли данные
	-h — сортировать по числовому значению с учетом суффиксов
*/

// Структура опциональных параметров сортировки.
type SortOptions struct {
	key     int  // номер колонки
	numeric bool // сортировка по целочисленному значению
	reverse bool // сортировка в обратном порядке
	unique  bool // не выводить повторяющиеся строки

	month        bool // сортировка по названию месяца
	ignoreBlanks bool // игнорировать ведущие и хвостовые пробельные символы
	check        bool // проверка, отсортированы ли строки
	humanNumeric bool // сортировка по целочисленному значению с суффиксом
}

// Тип функции сравнения двух строк.
type CompareFunc func(a, b string) int

// Шаблон для парсинга месяца в строке.
const MonthLayout = "Jan"

// ParseNumeric извлекает целочисленное значение из строки s.
func ParseNumeric(s string) int {
	re := regexp.MustCompile(`\d+`)
	value := re.FindString(s)
	d, _ := strconv.Atoi(value)
	return d
}

// CompareNumeric возвращает результат сравнения двух строк,
// содержащих целочисленные значения.
func CompareNumeric(a, b string) int {
	return ParseNumeric(a) - ParseNumeric(b)
}

// ParseNumeric извлекает целочисленное значение с определенным суффиксом из строки s.
func ParseNumericWithSuffix(s string) int {
	multipliers := map[string]int{"K": 1e3, "M": 1e6, "G": 1e9, "T": 1e12}

	re := regexp.MustCompile(`(\d+)([KMGT])`)
	values := re.FindStringSubmatch(s)
	if len(values) < 3 {
		return 0
	}

	multiplier, _ := multipliers[values[2]]
	d, _ := strconv.Atoi(values[1])
	return d * multiplier
}

// CompareNumericWithSuffix возвращает результат сравнения двух строк,
// содержащих целочисленные значения с определенными суффиксами.
func CompareNumericWithSuffix(a, b string) int {
	return ParseNumericWithSuffix(a) - ParseNumericWithSuffix(b)
}

// ParseMonth извлекает порядковый номер месяца из строки s.
func ParseMonth(s string) int {
	months := map[string]int{"jan": 0, "feb": 1, "mar": 2, "apr": 3, "may": 4, "jun": 5, "jul": 6, "aug": 7, "sep": 8, "oct": 9, "nov": 10, "dec": 11}
	re := regexp.MustCompile(`jan|feb|mar|apr|may|jun|jul|aug|sep|oct|nov|dec`)

	key := re.FindString(strings.ToLower(s))
	month, _ := months[key]
	return month
}

// CompareMonth возвращает результат сравнения двух строк,
// содержащих строковые представления месяцов.
func CompareMonth(a, b string) int {
	return ParseMonth(a) - ParseMonth(b)
}

// CompareKey возвращает функцию сравнения cmp для колонки с номером key.
// Колонка - часть строки, отделенная пробельными символами.
func CompareKey(cmp CompareFunc, key int) CompareFunc {
	return func(a, b string) int {
		colsA := strings.Fields(a)
		colsB := strings.Fields(b)
		if i := key - 1; i < len(colsA) && i < len(colsB) {
			return cmp(colsA[i], colsB[i])
		}
		return 0
	}
}

// Compare возвращает результат сравнения двух строк
// в зависимости от опциональных параметров.
func Compare(opts SortOptions) CompareFunc {
	var cmp CompareFunc

	switch {
	case opts.numeric:
		cmp = CompareNumeric
	case opts.humanNumeric:
		cmp = CompareNumericWithSuffix
	case opts.month:
		cmp = CompareMonth
	default:
		cmp = strings.Compare
	}

	if opts.key > 0 {
		return CompareKey(cmp, opts.key)
	}

	return cmp
}

// IgnoreBlanks преобразует исходный срез так,
// чтобы строки не содержали пробельные символы в начале и конце.
func IgnoreBlanks(in []string) {
	for i, str := range in {
		in[i] = strings.TrimSpace(str)
	}
}

// IsSorted возвращает результат проверки на то, отсортирован ли срез
// в порядке, определенном опциональными параметрами.
func IsSorted(in []string, opts SortOptions) bool {
	// Переворачиваем срез, если надо проверить сортировку в обратном порядке.
	if opts.reverse {
		slices.Reverse(in)
	}

	// Проверяем, отсортирован ли срез строк.
	isSorted := slices.IsSortedFunc(in, Compare(opts))

	// Проверяем, соблюдена ли строгая сортировка.
	if opts.unique {
		// Сортировка считается не строгой, если есть равные строки.
		return isSorted && len(in) == len(slices.Compact(in))
	}

	return isSorted
}

// Sort преобразует исходный срез так, чтобы строки были расположены
// в порядке, определенном опциональными параметрами.
func Sort(in *[]string, opts SortOptions) {
	strs := *in

	// Сортируем срез строк.
	slices.SortFunc(strs, Compare(opts))

	// Переворачиваем срез, если надо отсортировать в обратном порядке.
	if opts.reverse {
		slices.Reverse(strs)
	}

	// Удаляем дубликаты строк из среза.
	if opts.unique {
		*in = slices.Compact(strs)
	}
}

// Errorf пишет данные в stderr.
func Errorf(format string, a ...any) { fmt.Fprintf(os.Stderr, format, a...) }

// ReadStrings читает строки из r и записывает их в срез строк.
func ReadStrings(r io.Reader) ([]string, error) {
	strs := make([]string, 0)
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		strs = append(strs, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return strs, nil
}

// Переменные под параметры сортировки.
var (
	key                                                                int
	numeric, reverse, unique, month, ignoreBlanks, check, humanNumeric bool
)

func init() {
	// Настраиваем flag на парсинг параметров сортировки.
	flag.IntVar(&key, "k", 0, "sort via a key")
	flag.BoolVar(&numeric, "n", false, "compare according to string numerical value")
	flag.BoolVar(&reverse, "r", false, "reverse the result of comparisons")
	flag.BoolVar(&unique, "u", false, "with -c, check for strict ordering; without -c, output only the first of an equal run")
	flag.BoolVar(&month, "M", false, "compare (unknown) < 'JAN' < ... < 'DEC'")
	flag.BoolVar(&ignoreBlanks, "b", false, "ignore leading blanks")
	flag.BoolVar(&check, "c", false, "check for sorted input; do not sort")
	flag.BoolVar(&humanNumeric, "h", false, "compare human readable numbers (e.g., 2K 1G)")
}

func main() {
	// Парсим аргументы параметров сортировки.
	flag.Parse()
	opts := SortOptions{key, numeric, reverse, unique, month, ignoreBlanks, check, humanNumeric}

	// Парсим путь до файла из аргумента.
	path := flag.Arg(0)
	if path == "" {
		Errorf("empty file path")
		return
	}

	// Открываем файл на чтение.
	file, err := os.Open(path)
	if err != nil {
		Errorf("failed to open file: %s", err)
		return
	}

	// Читаем строки из файла и записываем в срез строк.
	strs, err := ReadStrings(file)
	if err != nil {
		Errorf("failed to read file: %s", err)
		file.Close()
		return
	}

	// Закрываем файл.
	if err := file.Close(); err != nil {
		Errorf("failed to close file: %s", err)
	}

	// Удаляем пробельные символы из строк.
	if opts.ignoreBlanks {
		IgnoreBlanks(strs)
	}

	// Проверяем, отсортированы ли строки.
	if opts.check {
		sorted := IsSorted(strs, opts)
		fmt.Println(sorted)
		return
	}

	// Сортируем строки и выводим отсортированный результат.
	Sort(&strs, opts)
	for _, str := range strs {
		fmt.Println(str)
	}
}
