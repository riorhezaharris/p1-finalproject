package handler

import (
	"errors"
	"fmt"
	"p1finalproject/entity"
)

func (h *Handler) GetUser(email string, password string) (entity.User, error) {
	query := fmt.Sprintf(`SELECT * 
	FROM users 
	WHERE users.email = '%s' && users.password = '%s';`, email, password)

	// Run the query
	var result entity.User
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
		err = rows.Scan(&result.UserId, &result.Email, &result.Password)
		if err != nil {
			return result, err
		}
	} else {
		return result, errors.New("invalid user credentials")
	}
	return result, nil
}
