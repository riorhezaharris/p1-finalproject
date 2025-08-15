package handler

import (
	"fmt"
	"strings"
	"p1finalproject/entity"
)

func (h *Handler) CreateOrder(userId int) (int, error) {
    // Ambil cart + subtotal
    cart, subtotal, err := h.GetCart(userId)
    if err != nil { return 0, err }
    if len(cart) == 0 { return 0, fmt.Errorf("cart empty") }

    // Transaksi biar atomic
    tx, err := h.db.Begin()
    if err != nil { return 0, err }
    defer func() {
        if err != nil { _ = tx.Rollback() }
    }()

    // 1) Buat record orders → dapat orderId
    res, err := tx.Exec(`INSERT INTO orders (user_id, total_price, status) VALUES (?, ?, 'waiting_for_payment')`,
        userId, subtotal)
    if err != nil { return 0, err }
    oid, err := res.LastInsertId()
    if err != nil { return 0, err }
    orderId := int(oid)

    // 2) Bulk insert order_details
    //    INSERT INTO order_details (product_id, quantity, price, order_id) VALUES (?, ?, ?, ?), (?, ?, ?, ?), ...
    placeholders := make([]string, 0, len(cart))
    args := make([]any, 0, len(cart)*4)
    for _, it := range cart {
        placeholders = append(placeholders, "(?, ?, ?, ?)")
        args = append(args, it.ProductId, it.Quantity, it.ProductPrice, orderId)
    }
    q := fmt.Sprintf(`INSERT INTO order_details (product_id, quantity, price, order_id) VALUES %s`,
        strings.Join(placeholders, ","))
    if _, err = tx.Exec(q, args...); err != nil { return 0, err }

    // 3) Kosongkan cart user
    if _, err = tx.Exec(`DELETE FROM carts WHERE user_id = ?`, userId); err != nil { return 0, err }

    // Commit
    if err = tx.Commit(); err != nil { return 0, err }
    return orderId, nil
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

func (h *Handler) GetOrderById(userId int, orderId int) (entity.Order, error) {
	query := fmt.Sprintf(`
	SELECT orders.id, users.email, orders.created_at, orders.total_price, orders.status
FROM orders
JOIN users ON orders.user_id = users.id
WHERE orders.user_id = %d && orders.id = %d
ORDER BY orders.created_at DESC;`, userId, orderId)

	// Run the query
	var result entity.Order
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
		err = rows.Scan(&result.OrderId, &result.UserEmail, &result.CreatedAt, &result.TotalPrice, &result.Status)
		if err != nil {
			return result, err
		}

		// Handle order details of the order
		result.OrderDetails, err = h.getOrderDetails(result.OrderId)
		if err != nil {
			return result, err
		}

		// If the order is paid, include the payment details
		if result.Status != "waiting_for_payment" {
			payment, err := h.getPayment(result.OrderId)
			if err != nil {
				return result, err
			}
			result.TotalPayment = payment.TotalPayment
			result.PaymentStatus = payment.Status
		}
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
