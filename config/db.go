package config

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

func InitDb() (*sql.DB, error) {
	// Load env
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env not found, using process env")
	}
	user := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	database := os.Getenv("DATABASE")

	// Open
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4,utf8",
		user, password, host, port, database)

	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("open mysql: %w", err)
	}
	db.SetConnMaxLifetime(10 * time.Minute)
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("ping mysql: %w", err)
	}

	// Seed all
	if err := seed(db); err != nil {
		return nil, fmt.Errorf("seeding failed: %w", err)
	}

	return db, nil
}

// ====================================
// SEEDERS
// ====================================

func seed(db *sql.DB) error {
	// 0) Users (2 akun) — pakai UPSERT by email
	if err := ensureUser(db, "demo@store.local", "demo123"); err != nil {
		return err
	}
	if err := ensureUser(db, "alice@store.local", "alice123"); err != nil {
		return err
	}

	// 1) Sizes
	if err := ensureSizes(db); err != nil {
		return err
	}

	// 2) Suppliers
	if err := ensureSuppliers(db); err != nil {
		return err
	}

	// 3) Products
	if err := ensureProducts(db); err != nil {
		return err
	}

	// 4) Unique key carts (abaikan error jika sudah ada)
	_, _ = db.Exec(`ALTER TABLE carts ADD UNIQUE KEY uq_user_product (user_id, product_id)`)

	return nil
}

// ---------- users ----------
func ensureUser(db *sql.DB, email, plain string) error {
	var cnt int
	if err := db.QueryRow(`SELECT COUNT(*) FROM users WHERE email=?`, email).Scan(&cnt); err != nil {
		return fmt.Errorf("cek user %s: %w", email, err)
	}
	if cnt > 0 {
		return nil
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password %s: %w", email, err)
	}
	if _, err := db.Exec(`INSERT INTO users (email, password) VALUES (?, ?)`, email, string(hashed)); err != nil {
		return fmt.Errorf("insert user %s: %w", email, err)
	}
	return nil
}

// ---------- sizes ----------
func ensureSizes(db *sql.DB) error {
	var n int
	if err := db.QueryRow(`SELECT COUNT(*) FROM sizes`).Scan(&n); err != nil {
		return fmt.Errorf("cek sizes: %w", err)
	}
	if n > 0 {
		return nil
	}
	_, err := db.Exec(`
		INSERT INTO sizes (name, width, length) VALUES
		('S', 48, 65),
		('M', 51, 68),
		('L', 54, 71),
		('XL', 57, 74)
	`)
	if err != nil {
		return fmt.Errorf("seed sizes: %w", err)
	}
	return nil
}

// ---------- suppliers ----------
func ensureSuppliers(db *sql.DB) error {
	var n int
	if err := db.QueryRow(`SELECT COUNT(*) FROM suppliers`).Scan(&n); err != nil {
		return fmt.Errorf("cek suppliers: %w", err)
	}
	if n > 0 {
		return nil
	}
	_, err := db.Exec(`
		INSERT INTO suppliers (name, location) VALUES
		('Alpha Textile', 'Bandung'),
		('Beta Apparel', 'Jakarta')
	`)
	if err != nil {
		return fmt.Errorf("seed suppliers: %w", err)
	}
	return nil
}

// ---------- products ----------
func ensureProducts(db *sql.DB) error {
	var n int
	if err := db.QueryRow(`SELECT COUNT(*) FROM products`).Scan(&n); err != nil {
		return fmt.Errorf("cek products: %w", err)
	}
	if n > 0 {
		return nil
	}

	sizeID := func(name string) (int, error) {
		var id int
		if err := db.QueryRow(`SELECT id FROM sizes WHERE name=? LIMIT 1`, name).Scan(&id); err != nil {
			return 0, fmt.Errorf("size '%s' tidak ditemukan: %w", name, err)
		}
		return id, nil
	}
	supplierID := func(name string) (int, error) {
		var id int
		if err := db.QueryRow(`SELECT id FROM suppliers WHERE name=? LIMIT 1`, name).Scan(&id); err != nil {
			return 0, fmt.Errorf("supplier '%s' tidak ditemukan: %w", name, err)
		}
		return id, nil
	}

	sM, err := sizeID("M"); if err != nil { return err }
	sL, err := sizeID("L"); if err != nil { return err }
	sXL, err := sizeID("XL"); if err != nil { return err }

	spA, err := supplierID("Alpha Textile"); if err != nil { return err }
	spB, err := supplierID("Beta Apparel");  if err != nil { return err }

	type row struct{ name string; price float64; size, supp int }
	rows := []row{
		{"T-Shirt Basic Putih", 79000, sM,  spA},
		{"T-Shirt Basic Hitam", 79000, sL,  spA},
		{"Hoodie Fleece Abu",   199000, sL, spB},
		{"Kemeja Oxford Biru",  159000, sM, spB},
		{"Celana Chino Khaki",  179000, sL, spA},
		{"Crewneck Navy",       149000, sXL, spB},
	}
	stmt, err := db.Prepare(`INSERT INTO products (name, price, size_id, supplier_id) VALUES (?,?,?,?)`)
	if err != nil {
		return fmt.Errorf("prepare insert products: %w", err)
	}
	defer stmt.Close()

	for _, r := range rows {
		if _, err := stmt.Exec(r.name, r.price, r.size, r.supp); err != nil {
			return fmt.Errorf("insert product '%s': %w", r.name, err)
		}
	}
	return nil
}

// ====================================

var ErrNotFound = errors.New("not found")
