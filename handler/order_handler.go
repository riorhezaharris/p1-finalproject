package handler

import (
	"fmt"
	"p1finalproject/entity"
)

func (h *Handler) CreateOrder(userId int) (int, error) {
	// Fetch the current cart
	cart, subtotal, err := h.GetCart(userId)
	if err != nil {
		return 0, err
	}

	// Create an order based the current cart
	query := fmt.Sprintf(`INSERT INTO orders (user_id, total_price, status) VALUES(%d, %.2f, 'waiting_for_payment');`, userId, subtotal)

	// Run the query
	createdOrder, err := h.db.Exec(query)
	if err != nil {
		return 0, err
	}

	orderId, err := createdOrder.LastInsertId()
	if err != nil {
		return 0, err
	}

	// Input all the cart items into order details
	query = `INSERT INTO order_details (product_id, quantity, price, order_id) VALUES(`
	for _, cartItem := range cart {
		query += fmt.Sprintf(`%d, %d, %.2f, %d`, cartItem.ProductId, cartItem.Quantity, cartItem.ProductPrice, orderId)
	}
	query += `);`

	// Run the query
	_, err = h.db.Exec(query)
	if err != nil {
		return 0, err
	}

	// Reset the cart after the order is created
	err = h.ResetCart(userId)
	if err != nil {
		return 0, err
	}

	return int(orderId), err
}

func (h *Handler) GetOrders(userId int) ([]entity.Order, error) {
	query := fmt.Sprintf(`
	SELECT orders.id, users.email, orders.created_at, orders.total_price, orders.status
FROM orders
JOIN users ON orders.user_id = users.id
WHERE orders.user_id = %d
ORDER BY orders.created_at DESC;`, userId)

	// Run the query
	var result []entity.Order
	rows, err := h.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Handling the result from database
	for rows.Next() {
		var row entity.Order
		err = rows.Scan(&row.OrderId, &row.UserEmail, &row.CreatedAt, &row.TotalPrice, &row.Status)
		if err != nil {
			return nil, err
		}

		// Handle order details of the order
		row.OrderDetails, err = h.getOrderDetails(row.OrderId)
		if err != nil {
			return nil, err
		}

		// If the order is paid, include the payment details
		if row.Status != "waiting_for_payment" {
			payment, err := h.getPayment(row.OrderId)
			if err != nil {
				return nil, err
			}
			row.TotalPayment = payment.TotalPayment
			row.PaymentStatus = payment.Status
		}
		result = append(result, row)
	}

	return result, nil
}

func (h *Handler) getOrderDetails(orderId int) ([]entity.OrderItem, error) {
	query := fmt.Sprintf(`
	SELECT products.name, sizes.name, order_details.quantity, order_details.price
FROM order_details
JOIN products ON order_details.product_id = products.id
JOIN sizes ON products.size_id = sizes.id
WHERE order_details.order_id = %d;`, orderId)

	// Run the query
	var result []entity.OrderItem
	rows, err := h.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Handling the result from database
	for rows.Next() {
		var row entity.OrderItem
		err = rows.Scan(&row.ProductName, &row.Sizename, &row.Quantity, &row.Price)
		if err != nil {
			return nil, err
		}
		result = append(result, row)
	}

	return result, nil
}

func (h *Handler) getPayment(orderId int) (entity.Payment, error) {
	query := fmt.Sprintf(`
	SELECT payments.total_payment, payments.status FROM payments
WHERE payments.order_id = %d;`, orderId)

	// Run the query
	var result entity.Payment
	rows, err := h.db.Query(query)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	if err := rows.Err(); err != nil {
		return result, err
	}

	// Handling the result from database
	if rows.Next() {
		err = rows.Scan(&result.TotalPayment, &result.Status)
		if err != nil {
			return result, err
		}
	}

	return result, nil
}
