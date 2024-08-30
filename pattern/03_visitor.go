package pattern

import (
	"math"
)

/*
Реализовать паттерн «посетитель».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
https://en.wikipedia.org/wiki/Visitor_pattern
*/

// Применимость: когда нужно добавить новую функциональность, изменяя исходные структуры минимально.
// Плюсы: упрощает добавление новой функциональности с минимальным изменением исходных структур, упрощает выполнение операций над полиморфными объектами.
// Минусы: для корректной реализации требует метод accept в исходных структурах, сложно применим к изменяющейся иерархии структур.
// Реальный пример: добавление функциональности определения итоговой стоимости разнородных товаров и услуг.

// Интерфейс некоторой фигуры.
type Figure interface {
	Name() string
	Accept(v FigureVisitor) float64
}

// Структура круга.
type Circle struct {
	r float64
}

// Получить название фигуры.
func (c *Circle) Name() string { return "circle" }

// Принимает посетителя дял реализации дополнительной функциональности.
func (c *Circle) Accept(v FigureVisitor) float64 { return v.VisitCircle(c) }

// Структура прямоугольника.
type Rectangle struct {
	a, b float64
}

// Получить название фигуры.
func (r *Rectangle) Name() string { return "rectangle" }

// Принимает посетителя дял реализации дополнительной функциональности.
func (r *Rectangle) Accept(v FigureVisitor) float64 { return v.VisitRectangle(r) }

// Интерфейс некоторого "посетителя фигур" для реализации дополнительной функциональности.
type FigureVisitor interface {
	VisitCircle(c *Circle) float64
	VisitRectangle(r *Rectangle) float64
}

// Конкретный "посетитель фигур" для вычисления периметра.
type FigurePerimeter struct {
}

// Реализовать вычисление периметра круга.
func (s *FigurePerimeter) VisitCircle(c *Circle) float64 { return 2 * math.Pi * c.r }

// Реализовать вычисление периметра прямоугольника.
func (s *FigurePerimeter) VisitRectangle(r *Rectangle) float64 { return 2 * (r.a + r.b) }

// Конкретный "посетитель фигур" для вычисления площади.
type FigureArea struct {
}

// Реализовать вычисление площади круга.
func (s *FigureArea) VisitCircle(c *Circle) float64 { return math.Pi * math.Pow(c.r, 2) }

// Реализовать вычисление площади прямоугольника.
func (s *FigureArea) VisitRectangle(r *Rectangle) float64 { return r.a * r.b }

/*
perimeter := new(FigurePerimeter)
area := new(FigureArea)
figures := []Figure{&Circle{r: 3}, &Rectangle{a: 2, b: 3}}
for _, figure := range figures {
	fmt.Printf("%s: perimeter = %f, area = %f\n", figure.Name(), figure.Accept(perimeter), figure.Accept(area))
}
*/
