package handler

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"p1finalproject/entity"
)

type Handler struct {
	DB *sql.DB
}

func NewHandler(db *sql.DB) *Handler { return &Handler{DB: db} }

// ---------- auth / bootstrap ----------
func (h *Handler) GetOrCreateUserByEmail(email string) (entity.User, error) {
	var u entity.User
	err := h.DB.QueryRow(`SELECT id, email FROM users WHERE email=? LIMIT 1`, email).Scan(&u.ID, &u.Email)
	if err == sql.ErrNoRows {
		res, err := h.DB.Exec(`INSERT INTO users (email, password) VALUES (?, ?)`, email, "")
		if err != nil {
			return u, err
		}
		id, _ := res.LastInsertId()
		u.ID = int(id)
		u.Email = email
		return u, nil
	}
	return u, err
}

// ---------- catalog ----------
func (h *Handler) ListProducts() ([]entity.Product, error) {
	rows, err := h.DB.Query(`
SELECT p.id, p.name, p.price, p.size_id, s.name AS size_name, p.supplier_id
FROM products p
LEFT JOIN sizes s ON s.id = p.size_id
ORDER BY p.id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []entity.Product
	for rows.Next() {
		var p entity.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.SizeID, &p.SizeName, &p.SupplierID); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (h *Handler) GetProductByCode(code int) (*entity.Product, error) {
	row := h.DB.QueryRow(`
SELECT p.id, p.name, p.price, p.size_id, s.name AS size_name, p.supplier_id
FROM products p
LEFT JOIN sizes s ON s.id = p.size_id
WHERE p.id=?`, code)
	var p entity.Product
	if err := row.Scan(&p.ID, &p.Name, &p.Price, &p.SizeID, &p.SizeName, &p.SupplierID); err != nil {
		return nil, err
	}
	return &p, nil
}

// ---------- cart ----------
func (h *Handler) AddToCart(userID, productID, qty int) error {
	// upsert sederhana: kalau sudah ada → tambah qty
	_, err := h.DB.Exec(`
INSERT INTO carts (user_id, product_id, quantity) VALUES (?,?,?)
ON DUPLICATE KEY UPDATE quantity = quantity + VALUES(quantity)`,
		userID, productID, qty)
	return err
}

func (h *Handler) GetCart(userID int) ([]entity.CartItem, float64, error) {
	rows, err := h.DB.Query(`
SELECT c.product_id, p.name, IFNULL(s.name,''), p.price, c.quantity, (p.price * c.quantity) AS subtotal
FROM carts c
JOIN products p ON p.id = c.product_id
LEFT JOIN sizes s ON s.id = p.size_id
WHERE c.user_id=?
ORDER BY c.product_id`, userID)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []entity.CartItem
	var total float64
	for rows.Next() {
		var it entity.CartItem
		if err := rows.Scan(&it.ProductID, &it.Name, &it.SizeName, &it.Price, &it.Quantity, &it.SubTotal); err != nil {
			return nil, 0, err
		}
		total += it.SubTotal
		items = append(items, it)
	}
	return items, total, rows.Err()
}

func (h *Handler) UpdateCartItem(userID, productID, qty int) error {
	if qty <= 0 {
		_, err := h.DB.Exec(`DELETE FROM carts WHERE user_id=? AND product_id=?`, userID, productID)
		return err
	}
	_, err := h.DB.Exec(`UPDATE carts SET quantity=? WHERE user_id=? AND product_id=?`, qty, userID, productID)
	return err
}

func (h *Handler) ClearCart(userID int) error {
	_, err := h.DB.Exec(`DELETE FROM carts WHERE user_id=?`, userID)
	return err
}

// ---------- checkout ----------
type CheckoutResult struct {
	OrderID     int
	Total       float64
	Paid        float64
	Change      float64
	OrderStatus string
	PayStatus   string
}

func (h *Handler) Checkout(userID int, paid float64) (*CheckoutResult, error) {
	items, total, err := h.GetCart(userID)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, errors.New("keranjang kosong")
	}

	status := "FAILED"
	payStatus := "FAILED"
	if paid >= total {
		status = "SUCCESS"
		payStatus = "PAID"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := h.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		// safety rollback jika belum commit
		_ = tx.Rollback()
	}()

	// orders
	res, err := tx.ExecContext(ctx, `
INSERT INTO orders (user_id, created_at, total_price, status)
VALUES (?, NOW(), ?, ?)`, userID, total, status)
	if err != nil {
		return nil, err
	}
	oid64, _ := res.LastInsertId()
	oid := int(oid64)

	// order_details
	stmt, err := tx.PrepareContext(ctx, `
INSERT INTO order_details (product_id, quantity, price, order_id) VALUES (?,?,?,?)`)
	if err != nil {
		return nil, err
	}
	for _, it := range items {
		if _, err := stmt.ExecContext(ctx, it.ProductID, it.Quantity, it.Price, oid); err != nil {
			return nil, err
		}
	}
	_ = stmt.Close()

	// payments
	_, err = tx.ExecContext(ctx, `
INSERT INTO payments (order_id, total_payment, status)
VALUES (?,?,?)`, oid, int(paid), payStatus)
	if err != nil {
		return nil, err
	}

	// bersihkan cart
	if _, err := tx.ExecContext(ctx, `DELETE FROM carts WHERE user_id=?`, userID); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &CheckoutResult{
		OrderID:     oid,
		Total:       total,
		Paid:        paid,
		Change:      paid - total,
		OrderStatus: status,
		PayStatus:   payStatus,
	}, nil
}

// ---------- history ----------
func (h *Handler) GetOrderHistory(userID int) ([]entity.OrderSummary, error) {
	rows, err := h.DB.Query(`
SELECT o.id, DATE_FORMAT(o.created_at, '%Y-%m-%d %H:%i:%s'), o.total_price, o.status,
       COALESCE((SELECT MAX(p.total_payment) FROM payments p WHERE p.order_id=o.id),0) AS paid
FROM orders o
WHERE o.user_id=?
ORDER BY o.id DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []entity.OrderSummary
	for rows.Next() {
		var s entity.OrderSummary
		var paid float64
		if err := rows.Scan(&s.ID, &s.CreatedAt, &s.TotalPrice, &s.Status, &paid); err != nil {
			return nil, err
		}
		s.Payment = paid
		out = append(out, s)
	}
	return out, rows.Err()
}
