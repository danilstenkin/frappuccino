package handlers

import (
	"encoding/json"
	"fmt"
	"frappuccino/db"
	"frappuccino/models"
	"frappuccino/repositories"
	"log"
	"net/http"
	"strings"
)

func CreateCustomersHandlers(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request to create customer")

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		log.Printf("Invalid method: %s. Expected POST.", r.Method)
		return
	}

	var person models.Customers
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&person)
	if err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if person.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		log.Printf("Error: Name is required in the request body")
		return
	}

	dbConn, err := db.InitDB()
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	defer dbConn.Close()

	id, err := repositories.CreateCustomers(dbConn, person)
	if err != nil {
		http.Error(w, "Could not create customers: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Error creating customers: %v", err)
		return
	}

	log.Printf("Menu item created successfully with ID: %d", id)

	response := map[string]int{"id": id}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// DELETE -------------------------------------------------------------------------------------------------------

func DeleteCustomersHandler(w http.ResponseWriter, r *http.Request) {
	const logPrefix = "[DeleteMenuItemHandler]"

	id := strings.TrimPrefix(r.URL.Path, "/customers/")

	if id == "" {
		log.Printf("%s Missing customers ID in URL", logPrefix)
		http.Error(w, "ID not found", http.StatusBadRequest)
		return
	}

	err := repositories.DeleteCustomer(id)
	if err != nil {
		log.Printf("%s Failed to delete customers: %v", logPrefix, err)
		http.Error(w, fmt.Sprintf("Failed to delete customers: %v", err), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
