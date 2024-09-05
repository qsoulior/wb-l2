package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/mitchellh/go-ps"
)

/*
=== Взаимодействие с ОС ===

Необходимо реализовать свой собственный UNIX-шелл-утилиту с поддержкой ряда простейших команд:

- cd <args> - смена директории (в качестве аргумента могут быть то-то и то)
- pwd - показать путь до текущего каталога
- echo <args> - вывод аргумента в STDOUT
- kill <args> - "убить" процесс, переданный в качесте аргумента (пример: такой-то пример)
- ps - выводит общую информацию по запущенным процессам в формате *такой-то формат*

Так же требуется поддерживать функционал fork/exec-команд

Дополнительно необходимо поддерживать конвейер на пайпах (linux pipes, пример cmd1 | cmd2 | .... | cmdN).

*Шелл — это обычная консольная программа, которая будучи запущенной, в интерактивном сеансе выводит некое приглашение
в STDOUT и ожидает ввода пользователя через STDIN. Дождавшись ввода, обрабатывает команду согласно своей логике
и при необходимости выводит результат на экран. Интерактивный сеанс поддерживается до тех пор, пока не будет введена команда выхода (например \quit).
*/

// Интерфейс команды, которую может выполнить шелл.
type Command interface {
	Execute(args ...string) (string, error)
}

// Структура команды cd, реализующая интерфейс.
type CD struct {
}

// Execute выполняет команду с аргументами args и возвращает строку результата или ошибку.
func (c CD) Execute(args ...string) (string, error) {
	if len(args) > 0 {
		dir := args[0]
		err := os.Chdir(dir)
		if err != nil {
			return "", fmt.Errorf("failed to cd %s: %w", dir, err)
		}
	}
	return "", nil
}

// Структура команды pwd, реализующая интерфейс.
type PWD struct {
}

// Execute выполняет команду с аргументами args и возвращает строку результата или ошибку.
func (c PWD) Execute(args ...string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to pwd: %w", err)
	}

	return wd, nil
}

// Структура команды echo, реализующая интерфейс.
type Echo struct {
}

// Execute выполняет команду с аргументами args и возвращает строку результата или ошибку.
func (c Echo) Execute(args ...string) (string, error) {
	var builder strings.Builder
	for i, arg := range args {
		builder.WriteString(arg)
		if i < len(args)-1 {
			builder.WriteRune(' ')
		}
	}
	return builder.String(), nil
}

// Структура команды kill, реализующая интерфейс.
type Kill struct {
}

// Execute выполняет команду с аргументами args и возвращает строку результата или ошибку.
func (c Kill) Execute(args ...string) (string, error) {
	if len(args) == 0 {
		return "", nil
	}
	pid, err := strconv.Atoi(args[0])
	if err != nil {
		return "", fmt.Errorf("failed to get pid %d: %w", pid, err)
	}
	process, err := os.FindProcess(pid)
	if err != nil {
		return "", fmt.Errorf("failed to find process %d: %w", pid, err)
	}
	err = process.Kill()
	if err != nil {
		return "", fmt.Errorf("failed to kill processs %d: %w", pid, err)
	}
	return fmt.Sprintf("process %d was killed", pid), nil
}

// Структура команды ps, реализующая интерфейс.
type PS struct {
}

// Execute выполняет команду с аргументами args и возвращает строку результата или ошибку.
func (c PS) Execute(args ...string) (string, error) {
	var limit int
	if len(args) > 0 {
		limit, _ = strconv.Atoi(args[0])
	}

	processes, err := ps.Processes()
	if err != nil {
		return "", fmt.Errorf("failed to list processes: %w", err)
	}

	if limit <= 0 || limit > len(processes) {
		limit = len(processes)
	}

	var builder strings.Builder
	builder.WriteString("PID\tPPID\tExecutable\n")
	for i := 0; i < limit; i++ {
		builder.WriteString(fmt.Sprintf("%d\t%d\t%s", processes[i].Pid(), processes[i].PPid(), processes[i].Executable()))
		if i < limit-1 {
			builder.WriteRune('\n')
		}
	}
	return builder.String(), nil
}

// Структура команды exec, реализующая интерфейс.
type Exec struct {
}

// Execute выполняет команду с аргументами args и возвращает строку результата или ошибку.
func (c Exec) Execute(args ...string) (string, error) {
	if len(args) == 0 {
		return "", nil
	}
	name := args[0]
	cmd := exec.Command(name, args[1:]...)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to exec: %w", err)
	}
	return string(output), nil
}

// Мапа доступных команд.
var commands = map[string]Command{
	"cd":   new(CD),
	"pwd":  new(PWD),
	"echo": new(Echo),
	"kill": new(Kill),
	"ps":   new(PS),
	"exec": new(Exec),
}

// ParseLine принимает строку для вызова команды и парсит ее.
// Возвращает объект с интерфейсом Command и срез аргументов или ошибку.
func ParseLine(line string) (Command, []string, error) {
	args := strings.Split(strings.TrimSpace(line), " ")
	cmd, ok := commands[args[0]]
	if !ok {
		return nil, nil, fmt.Errorf("command not found: %s", args[0])
	}
	return cmd, args[1:], nil
}

// ExecuteLines принимает строки для вызова команд, парсит каждую из них
// и вызывает по очереди, возвращая результат последней вызванной команды.
// Аргументы используются только в первой команде из среза.
func ExecuteLines(lines []string) (string, error) {
	if len(lines) == 0 {
		return "", nil
	}

	cmd, args, err := ParseLine(lines[0])
	if err != nil {
		return "", err
	}

	result, err := cmd.Execute(args...)
	if err != nil {
		return "", err
	}

	for i := 1; i < len(lines); i++ {
		cmd, _, err := ParseLine(lines[i])
		if err != nil {
			return "", err
		}

		result, err = cmd.Execute(result)
		if err != nil {
			return "", err
		}
	}

	return result, nil
}

func main() {
	fmt.Fprint(os.Stdout, "Welcome to shell! Enter \\q or \\quit to exit.\n> ")

	// Сканируем stdin построчно.
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "\\quit" || text == "\\q" {
			break
		}

		if text == "" {
			fmt.Fprint(os.Stdout, "> ")
			continue
		}

		// Если команда одна, то lines будет срезом из одного элемента.
		lines := strings.Split(text, "|")
		result, err := ExecuteLines(lines)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
		} else {
			fmt.Fprint(os.Stdout, result)
		}
		fmt.Fprint(os.Stdout, "\n> ")
	}
}
