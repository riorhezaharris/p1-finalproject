package entity

import (
	"time"
)

type Order struct {
	OrderId       int
	UserEmail     string
	CreatedAt     time.Time
	TotalPrice    float32
	Status        string
	OrderDetails  []OrderItem
	TotalPayment  float32
	PaymentStatus string
}

type OrderItem struct {
	ProductName string
	Sizename    string
	Quantity    int
	Price       float32
}
