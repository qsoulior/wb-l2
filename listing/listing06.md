Что выведет программа? Объяснить вывод программы. Рассказать про внутреннее устройство слайсов и что происходит при передачи их в качестве аргументов функции.

```go
package main

import (
	"fmt"
)

func main() {
	var s = []string{"1", "2", "3"}
	modifySlice(s)
	fmt.Println(s)
}

func modifySlice(i []string) {
	i[0] = "3"
	i = append(i, "4")
	i[1] = "5"
	i = append(i, "6")
}
```

Ответ:
```
[3 2 3]
```
Из четырех строчек тела `modifySlice` исходный массив, на который указывает срез `s` изменяется только первой строчкой. Причина этому кроется в устройстве среза.

Срез представляет из себя структуру вида:
```go
type slice struct {
	array unsafe.Pointer // указатель на массив с данными
	len   int            // длина слайса
	cap   int            // вместимость слайса
}
```
Каждый срез указывает на некоторый массив в памяти, содержит длину и емкость.
Срез можно взять от массива с помощью операции `arr[low:high:max]`, где `low` - индекс первого элемента в участке массива, на который указывает срез, а `high` -индекс элемента, идущего за последним элементом в участке массива.
Длина - это текущее количество элементов в массиве, с которыми может работать срез: `len=high-low`.
Емкость - это максимальное количество элементов в массиве начиная с элемента `low`, с которыми может работать срез без аллокации нового массива: `cap=max-low`.
Также срез можно создать сразу с аллокацией массива с помощью `make(тип_среза, len, cap)` или перечислением `[]int{1, 2, 3}`.

Функция `append` возвращает новый срез добавленными элементами. Если емкости массива не хватает, то `append` аллоцирует новый массив и новый срез связывает уже с ним. Тогда у нового среза помимо `len` изменится еще и `cap`.

Так как срез структура, то в при передаче в качестве аргумента по значению она копируется вместе с полями, однако указатель на массив продолжает указывать на тот же массив. Следовательно изменение элементов среза приводит к изменению элементов массива. Если на массив указывают и другие срезы, то из них так же будут видны эти изменения. В то же время функция `append` не изменяет исходный массив, а ее результат присваивается скопированной переменной среза, переданной в качестве аргумента. Изменение копии не приводит к изменению среза в вызывающей функции.