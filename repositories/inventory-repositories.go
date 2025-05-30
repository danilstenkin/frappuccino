package repositories

import (
	"database/sql"
	"fmt"
	"frappuccino/db"
	"frappuccino/models"
	"log"
	"strconv"
)

func GetInventoryItems() ([]models.InventoryItemResponce, error) {
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

	var items []models.InventoryItemResponce

	for rows.Next() {
		var item models.InventoryItemResponce
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

func CreateInventoryItems(item models.InventoryItem) (int, error) {
	dbConn, err := db.InitDB()
	if err != nil {
		return 0, fmt.Errorf("Не удалось подключится к базе данных, %v", err)
	}
	defer dbConn.Close()

	query := `INSERT INTO inventory (name, quantity, unit, price_per_unit) VALUES ($1, $2, $3, $4) RETURNING id`

	var id int

	err = dbConn.QueryRow(query, item.Name, item.Quantity, item.Unit, item.PricePerUnit).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("не удалось создать элемент инвентаря: %v", err)
	}

	return id, nil
}

func GetInventoryItemByID(idstr string) (models.InventoryItemResponce, error) {
	idInt, err := strconv.Atoi(idstr)
	if err != nil {
		return models.InventoryItemResponce{}, fmt.Errorf("ошибка при преобразовании ID: %v", err)
	}

	dbConn, err := db.InitDB()
	if err != nil {
		log.Println("Не удалось подключиться к БД:", err)
		return models.InventoryItemResponce{}, fmt.Errorf("не удалось подключиться к базе данных: %v", err)
	}
	defer dbConn.Close()

	var item models.InventoryItemResponce

	query := `SELECT id, name, quantity, unit, price_per_unit, last_updated FROM inventory WHERE id = $1`
	err = dbConn.QueryRow(query, idInt).Scan(&item.ID, &item.Name, &item.Quantity, &item.Unit, &item.PricePerUnit, &item.LastUpdated)

	if err == sql.ErrNoRows {
		return models.InventoryItemResponce{}, fmt.Errorf("инвентарь с таким ID не найден")
	} else if err != nil {
		return models.InventoryItemResponce{}, fmt.Errorf("ошибка при получении данных: %v", err)
	}

	return item, nil
}

func UpdateInventoryItem(idStr string, item models.InventoryItem) error {
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("неправильный формат ID: %v", err)
	}

	dbConn, err := db.InitDB()
	if err != nil {
		return fmt.Errorf("не удалось подключиться к БД: %v", err)
	}

	defer dbConn.Close()

	query := `UPDATE inventory SET name=$1, quantity=$2, unit=$3, price_per_unit=$4, last_updated=NOW() WHERE id=$5`

	result, err := dbConn.Exec(query, item.Name, item.Quantity, item.Unit, item.PricePerUnit, idInt)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении элемента меню: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка получения количества обновленных строк: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("элемент меню с ID %v не найден", idInt)
	}

	return nil
}

func DeleteInventoryItem(idstr string) error {
	idInt, err := strconv.Atoi(idstr)
	if err != nil {
		return fmt.Errorf("ошибка при преобразовании ID: %v", err)
	}

	// Подключаемся к базе данных
	dbConn, err := db.InitDB()
	if err != nil {
		log.Println("Не удалось подключиться к БД:", err)
		return fmt.Errorf("не удалось подключиться к базе данных: %v", err)
	}
	defer dbConn.Close()

	deleteIngredientsQuery := `DELETE FROM menu_item_ingredients WHERE ingredient_id = $1`
	_, err = dbConn.Exec(deleteIngredientsQuery, idInt)
	if err != nil {
		return fmt.Errorf("не удалось удалить зависимости из menu_item_ingredients: %v", err)
	}

	query := `DELETE FROM inventory WHERE id = $1`
	result, err := dbConn.Exec(query, idInt)
	if err != nil {
		return fmt.Errorf("ошибка при удалении элемента инвентаря: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка получения количества удалённых строк: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("элемент инвентаря с ID %v не найден", idInt)
	}

	return nil
}
