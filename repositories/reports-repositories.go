package repositories

import (
	"database/sql"
	"fmt"
	"strings"
)

func GetOrderedItemsByDay(db *sql.DB, month string) ([]map[string]int64, error) {
	// Преобразуем текстовое название месяца в номер (март → 3 и т.п.)
	monthNumber, err := monthNameToNumber(month)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT EXTRACT(DAY FROM o.order_date)::int AS day, SUM(oi.quantity)::bigint AS total
		FROM order_items oi
		JOIN orders o ON o.id = oi.order_id
		WHERE DATE_PART('month', o.order_date)::int = $1
		GROUP BY day
		ORDER BY day
	`

	rows, err := db.Query(query, monthNumber)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []map[string]int64
	for rows.Next() {
		var day int
		var total int64
		if err := rows.Scan(&day, &total); err != nil {
			return nil, err
		}
		result = append(result, map[string]int64{fmt.Sprintf("%d", day): total})
	}
	return result, nil
}

func GetOrderedItemsByMonth(db *sql.DB, year string) ([]map[string]int64, error) {
	query := `
		SELECT 
			TO_CHAR(o.order_date, 'Month') AS month, 
			DATE_PART('month', o.order_date) AS month_num,
			SUM(oi.quantity)::bigint AS total
		FROM order_items oi
		JOIN orders o ON o.id = oi.order_id
		WHERE EXTRACT(YEAR FROM o.order_date)::text = $1
		GROUP BY month, month_num
		ORDER BY month_num`

	rows, err := db.Query(query, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []map[string]int64
	for rows.Next() {
		var month string
		var monthNum int
		var total int64
		if err := rows.Scan(&month, &monthNum, &total); err != nil {
			return nil, err
		}
		result = append(result, map[string]int64{strings.ToLower(strings.TrimSpace(month)): total})
	}
	return result, nil
}

func monthNameToNumber(name string) (int, error) {
	months := map[string]int{
		"january": 1, "february": 2, "march": 3,
		"april": 4, "may": 5, "june": 6,
		"july": 7, "august": 8, "september": 9,
		"october": 10, "november": 11, "december": 12,
	}
	num, ok := months[strings.ToLower(strings.TrimSpace(name))]
	if !ok {
		return 0, fmt.Errorf("invalid month: %s", name)
	}
	return num, nil
}
