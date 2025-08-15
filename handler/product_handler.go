package handler

import (
	"fmt"
	"p1finalproject/entity"
)

func (h *Handler) GetProducts() ([]entity.Product, error) {
	query := `
	SELECT products.id, products.name, products.price, sizes.name, sizes.width, sizes.length, suppliers.name, suppliers.location
FROM products
JOIN suppliers ON products.supplier_id = suppliers.id
JOIN sizes ON products.size_id = sizes.id
ORDER BY products.id ASC;`

	// Run the query
	var result []entity.Product
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
		var row entity.Product
		err = rows.Scan(&row.ProductId, &row.Name, &row.Price, &row.SizeName, &row.SizeWidth, &row.SizeLength, &row.SupplierName, &row.SupplierLocation)
		if err != nil {
			return nil, err
		}
		result = append(result, row)
	}

	return result, nil
}

func (h *Handler) GetProductsById(productId int) (entity.Product, error) {
	query := fmt.Sprintf(`
	SELECT products.id, products.name, products.price, sizes.name, sizes.width, sizes.length, suppliers.name, suppliers.location
FROM products
JOIN suppliers ON products.supplier_id = suppliers.id
JOIN sizes ON products.size_id = sizes.id
WHERE products.id = %d
ORDER BY products.id ASC;`, productId)

	// Run the query
	var result entity.Product
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
		err = rows.Scan(&result.ProductId, &result.Name, &result.Price, &result.SizeName, &result.SizeWidth, &result.SizeLength, &result.SupplierName, &result.SupplierLocation)
		if err != nil {
			return result, err
		}
	}

	return result, nil
}
