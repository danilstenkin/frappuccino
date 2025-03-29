package repositories

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"frappuccino/models"
)

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
