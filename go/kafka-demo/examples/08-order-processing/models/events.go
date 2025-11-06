package models

import "time"

// EventType 事件类型
type EventType string

const (
	EventOrderCreated      EventType = "OrderCreated"
	EventInventoryReserved EventType = "InventoryReserved"
	EventPaymentCompleted  EventType = "PaymentCompleted"
	EventOrderCompleted    EventType = "OrderCompleted"
	EventOrderFailed       EventType = "OrderFailed"
)

// OrderItem 订单项
type OrderItem struct {
	ProductID string  `json:"product_id"`
	Name      string  `json:"name"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

// OrderCreated 订单创建事件
type OrderCreated struct {
	EventType   EventType   `json:"event_type"`
	OrderID     string      `json:"order_id"`
	UserID      string      `json:"user_id"`
	Items       []OrderItem `json:"items"`
	TotalAmount float64     `json:"total_amount"`
	Timestamp   time.Time   `json:"timestamp"`
	TraceID     string      `json:"trace_id"`
}

// InventoryReserved 库存预留事件
type InventoryReserved struct {
	EventType     EventType `json:"event_type"`
	OrderID       string    `json:"order_id"`
	ReservationID string    `json:"reservation_id"`
	Status        string    `json:"status"` // success, failed
	Reason        string    `json:"reason,omitempty"`
	Timestamp     time.Time `json:"timestamp"`
	TraceID       string    `json:"trace_id"`
}

// PaymentCompleted 支付完成事件
type PaymentCompleted struct {
	EventType EventType `json:"event_type"`
	OrderID   string    `json:"order_id"`
	PaymentID string    `json:"payment_id"`
	Status    string    `json:"status"` // success, failed
	Amount    float64   `json:"amount"`
	Reason    string    `json:"reason,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	TraceID   string    `json:"trace_id"`
}

// OrderCompleted 订单完成事件
type OrderCompleted struct {
	EventType EventType `json:"event_type"`
	OrderID   string    `json:"order_id"`
	UserID    string    `json:"user_id"`
	Status    string    `json:"status"` // completed, failed
	Timestamp time.Time `json:"timestamp"`
	TraceID   string    `json:"trace_id"`
}
