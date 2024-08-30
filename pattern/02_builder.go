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

// Структура заказа.
type Order struct {
	Items    []Item
	Payment  Payment
	Delivery Delivery
}

// Структура товара.
type Item struct {
	Name  string
	Price float64
}

// Структура оплаты.
type Payment struct {
	Amount float64
	Bank   string
}

// Структура доставки.
type Delivery struct {
	Addr  string
	Phone string
}

// Строитель, создающий заказ.
type OrderBuilder struct {
	order Order
}

func NewOrderBuilder() *OrderBuilder { return new(OrderBuilder) }

// Добавить товар в заказ.
func (b *OrderBuilder) AddItem(item Item) *OrderBuilder {
	b.order.Items = append(b.order.Items, item)
	return b
}

// Установить оплату заказа.
func (b *OrderBuilder) SetPayment(amount float64, bank string) *OrderBuilder {
	b.order.Payment = Payment{amount, bank}
	return b
}

// Установить доставку заказа.
func (b *OrderBuilder) SetDelivery(addr string, phone string) *OrderBuilder {
	b.order.Delivery = Delivery{addr, phone}
	return b
}

// Вернуть построенный заказ.
func (b *OrderBuilder) Build() Order { return b.order }

/*
order := NewOrderBuilder().
	AddItem(Item{"laptop", 1000}).AddItem(Item{"fridge", 5000}).
	SetPayment(6000, "bank").
	SetDelivery("6 some st.", "900").
	Build()
*/
