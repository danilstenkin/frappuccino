package repositories

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"frappuccino/db"
	"frappuccino/models"
	"log"
)

// CREATE CUSTOMER -----------------------------------------------------------------------------
func CreateCustomers(db *sql.DB, person models.Customers) (int, error) {
	preferencesJSON, err := json.Marshal(person.Preferences)
	if err != nil {
		return 0, fmt.Errorf("couldn not seralize customers_preference: %v", err)
	}

	query := `INSERT INTO customers (name, preferences) VALUES ($1, $2) RETURNING id`

	var id int

	err = db.QueryRow(query, person.Name, preferencesJSON).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("could not insert customer: %v", err)
	}

	return id, nil
}

// GET CUSTOMERS ----------------------------------------------------------------------------------

func GetCustomers() ([]models.CustomersResponse, error) {
	const logPrefix = "[GetCustomers]"

	dbConn, err := db.InitDB()
	if err != nil {
		log.Printf("%s Failed to connect to DB: %v", logPrefix, err)
		return nil, fmt.Errorf("failed to connect to the database: %v", err)
	}
	defer dbConn.Close()

	rows, err := dbConn.Query("SELECT id, name, preferences FROM customers")
	if err != nil {
		log.Printf("%s Failed to get customers: %v", logPrefix, err)
		return nil, fmt.Errorf("failed to get customers: %v", err)
	}
	defer rows.Close()

	var persons []models.CustomersResponse

	for rows.Next() {
		var person models.CustomersResponse
		var preferences sql.NullString

		err := rows.Scan(&person.ID, &person.Name, &preferences)
		if err != nil {
			return nil, fmt.Errorf("ошибка при сканировании строки: %v", err)
		}

		if preferences.Valid {
			err = json.Unmarshal([]byte(preferences.String), &person.Preferences)
			if err != nil {
				return nil, fmt.Errorf("ошибка при сканировании metadata: %v", err)
			}
		} else {
			person.Preferences = nil
		}

		persons = append(persons, person)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при итерации по строкам: %v", err)
	}

	return persons, nil
}
