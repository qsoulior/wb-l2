package pattern

/*
Реализовать паттерн «строитель».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
https://en.wikipedia.org/wiki/Builder_pattern
*/

// Применимость: когда нужно пошагово создавать разные представления сложного объекта.
// Плюсы: позволяет избавиться от сложных конструкторов с большим количеством аргументов, инкапсулирует логику сборки объекта, делает сборку более реюзабельной.
// Минусы: требует внедрения дополнительных классов и методов, что усложняет код программы.
// Реальный пример: сборка заказа, состоящего из нескольких этапов и частей (товаров, оплаты и доставки).

// Структура, объекты которой собирает строитель.
type BuilderProduct struct {
	x, y, z int
}

// Интерфейс строителя.
type Builder interface {
	Reset() Builder
	SetX(x int) Builder
	SetY(y int) Builder
	SetZ(z int) Builder
	Build() BuilderProduct
}

// Структура конкретного строителя, реализующая интерфейс.
type ConcreteBuilder struct {
	product BuilderProduct
}

// Сбрасывает состояние собираемого объекта.
func (b *ConcreteBuilder) Reset() Builder {
	b.product = BuilderProduct{}
	return b
}

// Устанавливает значение X собираемого объекта.
func (b *ConcreteBuilder) SetX(x int) Builder {
	b.product.x = x
	return b
}

// Устанавливает значение Y собираемого объекта.
func (b *ConcreteBuilder) SetY(y int) Builder {
	b.product.y = y
	return b
}

// Устанавливает значение Z собираемого объекта.
func (b *ConcreteBuilder) SetZ(z int) Builder {
	b.product.z = z
	return b
}

// Возвращает собранный объект.
func (b ConcreteBuilder) Build() BuilderProduct { return b.product }

// Структура директора, определяющего порядок вызова шагов строителя.
type Director struct {
	builder Builder
}

// Устанавливает активного строителя.
func (d *Director) SetBuilder(builder Builder) { d.builder = builder }

// Собирает и возвращает объект.
func (d Director) GetProduct() BuilderProduct { return d.builder.Reset().SetX(10).SetY(20).Build() }

/*
// Собираем объект без директора.
builder := ConcreteBuilder{}
product1 := builder.SetY(10).SetZ(5).Build()
fmt.Println(product1) // {0 10 5}

// Собираем объект с директором.
director := Director{}
director.SetBuilder(&builder)
product2 := director.GetProduct()
fmt.Println(product2) // {10 20 0}
*/
