package handler

import (
	"errors"
	"fmt"
	"p1finalproject/entity"
)

func (h *Handler) GetCart(userId int) ([]entity.CartItem, float32, error) {
	query := fmt.Sprintf(`SELECT carts.id, carts.product_id, products.name, products.price, sizes.name, quantity
FROM carts
JOIN users ON carts.user_id = users.id
JOIN products ON carts.product_id = products.id
JOIN sizes ON products.size_id = sizes.id
WHERE carts.user_id = %d`, userId)

	// Run the query
	var result []entity.CartItem
	rows, err := h.db.Query(query)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	// Handling the result from database
	for rows.Next() {
		var row entity.CartItem
		err = rows.Scan(&row.CartItemId, &row.ProductId, &row.ProductName, &row.ProductPrice, &row.ProductSize, &row.Quantity)
		if err != nil {
			return nil, 0, err
		}
		result = append(result, row)
	}

	// Calculate the subtotal from the cart items
	var calculation float32
	for _, item := range result {
		calculation += item.ProductPrice
	}

	return result, calculation, nil
}

func (h *Handler) UpdateQuantity(userId int, productId int, quantity int) error {
	// Check if the product exist in the cart
	cartItem, err := h.findCartItem(userId, productId)
	if err != nil {
		return err
	}

	// If the product exist in the cart, just add the quantity
	var query string
	if len(cartItem) == 0 {
		query = fmt.Sprintf(`INSERT INTO carts (user_id, product_id, quantity) VALUES(%d, %d, %d);`, userId, productId, quantity)
	} else {
		query = fmt.Sprintf(`UPDATE carts SET quantity = %d WHERE carts.id = %d;`, cartItem[0].Quantity+quantity, cartItem[0].CartItemId)
	}

	// Run the query
	_, err = h.db.Exec(query)
	return err
}

func (h *Handler) RemoveItem(userId int, productId int) error {
	// Check if the product exist in the cart
	cartItem, err := h.findCartItem(userId, productId)
	if err != nil {
		return err
	}

	// Error handling if the product not found
	var query string
	if len(cartItem) == 0 {
		return errors.New("item not found in the cart")
	}

	// Remove the product from the cart
	query = fmt.Sprintf(`DELETE FROM carts WHERE carts.id = %d;`, cartItem[0].CartItemId)

	// Run the query
	_, err = h.db.Exec(query)
	return err
}

func (h *Handler) ResetCart(userId int) error {
	query := fmt.Sprintf(`DELETE FROM carts WHERE carts.user_id = %d;`, userId)

	// Run the query
	_, err := h.db.Exec(query)
	return err
}

func (h *Handler) findCartItem(userId int, productId int) ([]entity.CartItem, error) {
	query := fmt.Sprintf(`SELECT carts.id, carts.product_id, products.name, products.price, quantity
FROM carts
JOIN users ON carts.user_id = users.id
JOIN products ON carts.product_id = products.id
WHERE carts.user_id = %d && carts.product_id = %d`, userId, productId)

	// Run the query
	var result []entity.CartItem
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
		var row entity.CartItem
		err = rows.Scan(&row.CartItemId, &row.ProductId, &row.ProductName, &row.ProductPrice, &row.Quantity)
		if err != nil {
			return nil, err
		}
		result = append(result, row)
	}

	return result, nil
}
