package pattern

import "fmt"

/*
Реализовать паттерн «фасад».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
https://en.wikipedia.org/wiki/Facade_pattern
*/

// Применимость: когда нужно предоставить упрощенный интерфейс к некоторой сложной системе.
// Плюсы: обособляет и инкапсулирует сложную систему от пользователя, предоставляя достаточное взаимодействие.
// Минусы: может стать слишком крупным, что усложнит его использование и тестирование.
// Реальный пример: сервис заказов, использующий множество других подсервисов (товаров, оплаты и доставки).

// Некоторый сервис, отвечающий за товары в заказе.
type ItemService struct {
	items []any
}

// Добавить товар.
func (s *ItemService) AddItem(item any) { fmt.Printf("AddItem: %v\n", item) }

// Удалить товар.
func (s *ItemService) RemoveItem(itemID string) { fmt.Printf("RemoveItem: %s\n", itemID) }

// Получить список добавленных товаров.
func (s *ItemService) GetItems() []any { return s.items }

// Некоторый сервис, отвечающий за оплату заказа.
type PaymentService struct {
	status string
	amount float64
}

// Выполнить оплату.
func (s *PaymentService) ProcessPayment(amount float64) { fmt.Printf("ProcessPayment: %f\n", amount) }

// Вернуть оплату.
func (s *PaymentService) RefundPayment() { fmt.Println("RefundPayment") }

// Получить статус оплаты.
func (s *PaymentService) GetPaymentStatus() string { return s.status }

// Некоторый сервис, отвечающий за доставку заказа.
type DeliveryService struct {
	state    string
	statuses []string
}

// Организовать доставку по адресу.
func (s *DeliveryService) ArrangeDelivery(addr string) { fmt.Printf("ArrangeDelivery: %s\n", addr) }

func (s *DeliveryService) TrackDelivery() []string { return s.statuses }

func (s *DeliveryService) CancelDelivery() { fmt.Println("CancelDelivery") }

// Сервис, реализующий упрощенный интерфейс к сложной системе и отвечающий за заказ.
type OrderService struct {
	itemService     *ItemService
	paymentService  *PaymentService
	deliveryService *DeliveryService
}

// Создание сервиса для заказов
func NewOrderService() *OrderService {
	return &OrderService{
		itemService:     &ItemService{},
		paymentService:  &PaymentService{},
		deliveryService: &DeliveryService{},
	}
}

// Сделать заказ.
func (s *OrderService) MakeOrder(items []any, amount float64, address string) {
	for _, item := range items {
		s.itemService.AddItem(item)
	}
	s.paymentService.ProcessPayment(amount)
	s.deliveryService.ArrangeDelivery(address)
}

// Отменить заказ.
func (s *OrderService) CancelOrder() {
	s.paymentService.RefundPayment()
	s.deliveryService.CancelDelivery()
}

// Получить статус заказа.
func (s *OrderService) GetOrderStatus() {
	fmt.Printf("Items: %v\n", s.itemService.GetItems())
	fmt.Printf("Payment Status: %s\n", s.paymentService.GetPaymentStatus())
	fmt.Printf("Delivery Status: %v\n", s.deliveryService.TrackDelivery())
}
