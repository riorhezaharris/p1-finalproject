package handler

import (
	"fmt"
	"p1finalproject/entity"
)

func (h *Handler) CreatePayment(userId int, orderId int, totalPayment float32) error {
	query := fmt.Sprintf(`
	SELECT orders.id, users.email, orders.created_at, orders.total_price, orders.status
FROM orders
JOIN users ON orders.user_id = users.id
WHERE orders.user_id = %d && orders.id = %d
ORDER BY orders.created_at DESC;`, userId, orderId)

	// Run the query
	var order entity.Order
	rows, err := h.db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	if err := rows.Err(); err != nil {
		return err
	}

	// Handling the result from database
	if rows.Next() {
		err = rows.Scan(&order.OrderId, &order.UserEmail, &order.CreatedAt, &order.TotalPrice, &order.Status)
		if err != nil {
			return err
		}
	}

	// Handle if the payment not sufficient
	var status string
	if order.TotalPrice > totalPayment {
		status = "insufficient"
	} else {
		status = "sufficient"
	}

	query = fmt.Sprintf(`INSERT INTO payments (order_id, total_payment, status) VALUES(%d, %.2f, '%s');`, orderId, totalPayment, status)

	// Run the query
	_, err = h.db.Exec(query)
	if err != nil {
		return err
	}

	// Update the order status based on the payment
	var orderStatus string
	if status == "insufficient" {
		orderStatus = "failed"
	} else {
		orderStatus = "success"
	}
	query = fmt.Sprintf(`UPDATE orders SET status = '%s' WHERE orders.id = %d;`, orderStatus, orderId)

	// Run the query
	_, err = h.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
