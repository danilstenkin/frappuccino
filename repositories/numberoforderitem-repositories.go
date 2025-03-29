package repositories

import (
	"database/sql"
	"frappuccino/models"
	"time"
)

// Функция для получения количества заказанных товаров
func NumberOfOrderedItems(db *sql.DB, startDate, endDate *time.Time) ([]models.OrderedItem, error) {
	// Выполняем SQL-запрос для получения количества заказанных позиций по заданным датам
	rows, err := db.Query(`
        SELECT mi.name, SUM(oi.quantity)
        FROM order_items oi
        JOIN orders o ON o.id = oi.order_id
        JOIN menu_items mi ON mi.id = oi.menu_item_id
        WHERE ($1::date IS NULL OR o.order_date >= $1)
          AND ($2::date IS NULL OR o.order_date <= $2)
        GROUP BY mi.name
    `, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.OrderedItem
	for rows.Next() {
		var item models.OrderedItem
		if err := rows.Scan(&item.Name, &item.Quantity); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}
