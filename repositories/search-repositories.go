package repositories

import (
	"database/sql"
	"frappuccino/models"
)

func SearchMenuItems(db *sql.DB, query string, minPrice, maxPrice float64) ([]models.SearchMenuItem, error) {
	rows, err := db.Query(`
		SELECT id, name, description, price,
		       ts_rank(to_tsvector('english', name || ' ' || description), plainto_tsquery('english', $1)) AS relevance
		FROM menu_items
		WHERE to_tsvector('english', name || ' ' || description) @@ plainto_tsquery('english', $1)
		  AND price >= $2 AND price <= $3
		ORDER BY relevance DESC
	`, query, minPrice, maxPrice)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.SearchMenuItem
	for rows.Next() {
		var item models.SearchMenuItem
		if err := rows.Scan(&item.ID, &item.Name, &item.Description, &item.Price, &item.Relevance); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func SearchOrders(db *sql.DB, query string, minPrice, maxPrice float64) ([]models.SearchOrder, error) {
	rows, err := db.Query(`
		SELECT DISTINCT o.id, c.name, o.total_amount,
		  ts_rank(
		    setweight(to_tsvector('english', c.name), 'A') ||
		    setweight(to_tsvector('english', string_agg(mi.name, ' ')), 'B'),
		    plainto_tsquery('english', $1)
		  ) AS relevance
		FROM orders o
		JOIN customers c ON o.customer_id = c.id
		JOIN order_items oi ON oi.order_id = o.id
		JOIN menu_items mi ON mi.id = oi.menu_item_id
		WHERE to_tsvector('english', c.name || ' ' || mi.name) @@ plainto_tsquery('english', $1)
		  AND o.total_amount >= $2 AND o.total_amount <= $3
		GROUP BY o.id, c.name, o.total_amount
		ORDER BY relevance DESC;
	`, query, minPrice, maxPrice)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.SearchOrder
	for rows.Next() {
		var order models.SearchOrder
		if err := rows.Scan(&order.ID, &order.CustomerName, &order.Total, &order.Relevance); err != nil {
			return nil, err
		}

		itemRows, err := db.Query(`
			SELECT mi.name
			FROM order_items oi
			JOIN menu_items mi ON mi.id = oi.menu_item_id
			WHERE oi.order_id = $1
		`, order.ID)
		if err != nil {
			return nil, err
		}
		var items []string
		for itemRows.Next() {
			var name string
			if err := itemRows.Scan(&name); err != nil {
				return nil, err
			}
			items = append(items, name)
		}
		itemRows.Close()
		order.Items = items
		results = append(results, order)
	}
	return results, nil
}
