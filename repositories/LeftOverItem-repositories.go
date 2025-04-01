package repositories

import (
	"database/sql"
	"fmt"
	"frappuccino/models"
)

func GetLeftOvers(db *sql.DB, sortBy string, page, pageSize int) ([]models.LeftOverItem, int, error) {
	if sortBy != "price" && sortBy != "quantity" {
		sortBy = "quantity"
	}

	// Считаем общее количество записей
	var total int
	err := db.QueryRow("SELECT COUNT(*) FROM inventory").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	query := fmt.Sprintf(`
		SELECT name, quantity, price_per_unit
		FROM inventory
		ORDER BY %s DESC
		LIMIT $1 OFFSET $2
	`, sortBy)

	rows, err := db.Query(query, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []models.LeftOverItem
	for rows.Next() {
		var item models.LeftOverItem
		if err := rows.Scan(&item.Name, &item.Quantity, &item.Price); err != nil {
			return nil, 0, err
		}
		items = append(items, item)
	}

	totalPages := (total + pageSize - 1) / pageSize
	return items, totalPages, nil
}
