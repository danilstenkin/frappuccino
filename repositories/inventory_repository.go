package repositories

import (
	"fmt"
	"frappuccino/db"
	"frappuccino/models"
)

func GetInventoryItems() ([]models.InventoryItem, error) {
	dbConn, err := db.InitDB()
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к БД: %v", err)
	}
	defer dbConn.Close()

	rows, err := dbConn.Query(`SELECT id, name, quantity, unit, price_per_unit, last_updated FROM inventory`)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить инвентарь: %v", err)
	}
	defer rows.Close()

	var items []models.InventoryItem

	for rows.Next() {
		var item models.InventoryItem
		err := rows.Scan(&item.ID, &item.Name, &item.Quantity, &item.Unit, &item.PricePerUnit, &item.LastUpdated)
		if err != nil {
			return nil, fmt.Errorf("ошибка при сканировании строки: %v", err)
		}
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при итерации по строкам: %v", err)
	}

	return items, nil
}
