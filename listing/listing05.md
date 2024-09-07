Что выведет программа? Объяснить вывод программы.

```go
package main

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func test() *customError {
	{
		// do something
	}
	return nil
}

func main() {
	var err error
	err = test()
	if err != nil {
		println("error")
		return
	}
	println("ok")
}
```

Ответ:
```
error
```
Причина кроется в устройстве интерфейса, подробно описанном в листинге 3. При присваивании `err = test()` в поле `data` интерфейса записывается `nil`-значение `*customError`, тогда как в поле `tab` создается таблица интерфейса, связывающая статический тип интерфейса с конкретным типом `*customError`. Структура интерфейса не является `nil`: следовательно, условие в `if` выполняется и выводится "error".