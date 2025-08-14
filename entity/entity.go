package entity

type User struct {
	ID    int
	Email string
}

type Size struct {
	ID     int
	Name   string
	Width  int
	Length int
}

type Supplier struct {
	ID       int
	Name     string
	Location string
}

type Product struct {
	ID         int
	Name       string
	Price      float64
	SizeID     int
	SizeName   string
	SupplierID int
}

type CartItem struct {
	ProductID int
	Name      string
	SizeName  string
	Price     float64
	Quantity  int
	SubTotal  float64
}

type OrderSummary struct {
	ID         int
	CreatedAt  string
	TotalPrice float64
	Status     string
	Payment    float64
}
