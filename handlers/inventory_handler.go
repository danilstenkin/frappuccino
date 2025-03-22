package handlers

import (
	"encoding/json"
	"frappuccino/models"
	"frappuccino/repositories"
	"log"
	"net/http"
	"strings"
)

func GetInventoryHandler(w http.ResponseWriter, r *http.Request) {
	items, err := repositories.GetInventoryItems()
	if err != nil {
		http.Error(w, "Не удалось получить инвентарь: "+err.Error(), http.StatusInternalServerError)
		log.Println("Ошибка получения инвентаря:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func CreateInventoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request metgod", http.StatusMethodNotAllowed)
	}

	var item models.InventoryItem

	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
	}

	if item.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	if item.Quantity <= 0 {
		http.Error(w, "Quantity must be greater than 0", http.StatusBadRequest)
		return
	}

	if item.Unit == "" {
		http.Error(w, "Unit is required", http.StatusBadRequest)
		return
	}
	if item.PricePerUnit <= 0 {
		http.Error(w, "Price per unit must be greater than 0", http.StatusBadRequest)
		return
	}

	id, err := repositories.CreateInventoryItems(item)
	if err != nil {
		http.Error(w, "Could not create inventory item: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]int{"id": id}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func GetInventoryByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/inventory/")

	if id == "" {
		http.Error(w, "ID not found", http.StatusBadRequest)
		return
	}

	item, err := repositories.GetInventoryItemByID(id)
	if err != nil {
		http.Error(w, "Не удалось получить элемент инвентаря: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}
