package repositories

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"frappuccino/db"
	"frappuccino/models"
	"log"
	"strconv"
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

// DELETE CUSTOMER ------------------------------------------------------------------------------

func DeleteCustomer(idStr string) error {
	const logPrefix = "[DeleteCustomer]"

	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("%s incorrect ID: %v", logPrefix, err)
		return fmt.Errorf("Error converting ID: %v", err)
	}

	dbConn, err := db.InitDB()
	if err != nil {
		log.Printf("%s Failed to connect to DB: %v", logPrefix, err)
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	log.Printf("%s Connection to the database was successful", logPrefix)

	log.Printf("%s Remove item from Customers with ID= %d", logPrefix, idInt)
	query := `DELETE FROM customers WHERE id = $1`
	result, err := dbConn.Exec(query, idInt)
	if err != nil {
		log.Printf("%s Error while deleting from customers: %v", logPrefix, err)
		return fmt.Errorf("failed to delete customers: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("%s Error getting number of deleted rows: %v", logPrefix, err)
		return fmt.Errorf("error getting number of deleted rows: %v", err)
	}
	if rowsAffected == 0 {
		log.Printf("%s Customers with ID %d not found", logPrefix, idInt)
		return fmt.Errorf("Customers with ID %v not found", idInt)
	}

	log.Printf("%s Customers with ID %d successfully removed", logPrefix, idInt)
	return nil
}
