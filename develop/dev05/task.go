package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

/*
=== Утилита grep ===

Реализовать утилиту фильтрации (man grep)

Поддержать флаги:
-A - "after" печатать +N строк после совпадения
-B - "before" печатать +N строк до совпадения
-C - "context" (A+B) печатать ±N строк вокруг совпадения
-c - "count" (количество строк)
-i - "ignore-case" (игнорировать регистр)
-v - "invert" (вместо совпадения, исключать)
-F - "fixed", точное совпадение со строкой, не паттерн
-n - "line num", печатать номер строки

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

// Структура опциональных параметров выборки.
type GrepOptions struct {
	after   int  // печатать N строк после совпадения
	before  int  // печатать N строк до совпадения
	context int  // печатать N строк вокруг совпадения
	count   bool // количество строк
	invert  bool // вместо совпадения, исключать
	num     bool // печатать номер строки
}

// Структура опциональных параметров соответствия.
type MatchOptions struct {
	ignoreCase bool // игнорировать регистр
	fixed      bool // точное совпадение со строкой, не паттерн
}

// Интерфейс правила соответствия.
type Matcher interface {
	Match(str string) bool
}

// Структура правила точного соответствия.
type fixedMatcher struct {
	cmp func(str string) bool
}

// NewFixedMatcher создает и возвращает объект FixedMatcher.
func NewFixedMatcher(pattern string, opts MatchOptions) Matcher {
	if opts.ignoreCase {
		return fixedMatcher{cmp: func(in string) bool { return strings.EqualFold(in, pattern) }}
	}
	return fixedMatcher{cmp: func(in string) bool { return in == pattern }}
}

// Match возвращает результат сравнения str с паттерном.
func (m fixedMatcher) Match(str string) bool { return m.cmp(str) }

// Структура правила соответствия по паттерну.
type regexpMatcher struct {
	re *regexp.Regexp
}

// NewFixedMatcher создает и возвращает объект NewRegexpMatcher.
func NewRegexpMatcher(pattern string, opts MatchOptions) (Matcher, error) {
	if opts.ignoreCase {
		pattern = "(?i)" + pattern
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	return regexpMatcher{re}, nil
}

// Match возвращает результат сравнения str с паттерном.
func (m regexpMatcher) Match(str string) bool { return m.re.MatchString(str) }

// Grep проверяет строки из input на соответствие правилу из matcher.
// Возвращает срез, содержащий результаты сравнения строк с паттерном.
func Grep(input []string, matcher Matcher, opts GrepOptions) []bool {
	n := len(input)
	matched := make([]bool, n)
	for i := 0; i < n; i++ {
		// Пропускаем строку, если она не соответствует паттерну.
		if (!opts.invert && !matcher.Match(input[i])) || (opts.invert && matcher.Match(input[i])) {
			continue
		}
		// Отмечаем текущую строку.
		matched[i] = true
		// Отмечаем N строк после совпадения.
		for j := 1; (j <= opts.after || j <= opts.context) && j < n; j++ {
			matched[i+j] = true
		}
		// Отмечаем N строк до совпадения.
		for j := 1; (j <= opts.before || j <= opts.context) && i >= j; j++ {
			matched[i-j] = true
		}
	}
	return matched
}

// PrintCount выводит количество строк, которые соответствуют паттерну.
func PrintCount(matched []bool) {
	count := 0
	for _, m := range matched {
		if m {
			count++
		}
	}
	fmt.Println(count)
}

// PrintStrings выводит строки, которые соответствуют паттерну
func PrintStrings(matched []bool, input []string, num bool) {
	for i, m := range matched {
		if m {
			if num {
				fmt.Printf("%d:", i+1)
			}
			fmt.Println(input[i])
		}
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

// Переменные под параметры grep.
var (
	after, before, context                int
	count, ignoreCase, invert, fixed, num bool
)

func init() {
	// Настраиваем flag на парсинг параметров grep.
	flag.IntVar(&after, "A", 0, "print NUM lines of trailing context after matching lines")
	flag.IntVar(&before, "B", 0, "print NUM lines of leading context before matching lines")
	flag.IntVar(&context, "C", 0, "print NUM lines of output context")
	flag.BoolVar(&count, "c", false, "suppress normal output; instead print a count of matching lines for each input file.")
	flag.BoolVar(&ignoreCase, "i", false, "ignore case distinctions in patterns and input data.")
	flag.BoolVar(&invert, "v", false, "invert the sense of matching, to select non-matching lines")
	flag.BoolVar(&fixed, "F", false, "interpret PATTERNS as fixed strings, not regular expressions")
	flag.BoolVar(&num, "n", false, "prefix each line of output with the 1-based line number within its input file")
}

// grep [OPTION...] PATTERN [FILE...]
func main() {
	// Парсим аргументы параметров сортировки.
	flag.Parse()

	// Парсим путь до файла из аргумента.
	args := flag.Args()
	if len(args) < 2 {
		Errorf("too many command-line arguments")
		return
	}

	var (
		matcher Matcher
		err     error
	)

	// Назначем правило соотвествия.
	matchOpts := MatchOptions{ignoreCase, fixed}
	if pattern := args[0]; matchOpts.fixed {
		matcher = NewFixedMatcher(pattern, matchOpts)
	} else {
		matcher, err = NewRegexpMatcher(pattern, matchOpts)
		if err != nil {
			Errorf("failed to compile pattern: %s", err)
			return
		}
	}

	// Собираем срез строк из нескольких файлов.
	input := make([]string, 0)
	for _, arg := range args[1:] {
		// Открываем файл на чтение.
		file, err := os.Open(arg)
		if err != nil {
			Errorf("failed to open file: %s", err)
			continue
		}
		// Читаем строки из файла.
		part, err := ReadStrings(file)
		if err != nil {
			Errorf("failed to read file: %s", err)
			file.Close()
			continue
		}
		// Записываем строки в общий срез строк.
		input = append(input, part...)
		// Закрываем файл.
		if err := file.Close(); err != nil {
			Errorf("failed to close file: %s", err)
		}
	}

	// Вычисляем индексы строк, которые соответствуют паттерну.
	grepOpts := GrepOptions{after, before, context, count, invert, num}
	matched := Grep(input, matcher, grepOpts)

	// Выводим количество строк, которые соответствуют паттерну.
	if grepOpts.count {
		PrintCount(matched)
		return
	}

	// Выводим строки, которые соответствуют паттерну.
	PrintStrings(matched, input, grepOpts.num)
}
