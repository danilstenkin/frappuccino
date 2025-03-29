package handlers

import (
	"encoding/json"
	"frappuccino/db"
	"frappuccino/models"
	"frappuccino/repositories"
	"log"
	"net/http"
	"time"
)

// Обработчик для получения количества заказанных позиций
func GetNumberOfOrderedItemsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request to get number of ordered items")

	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		log.Printf("Invalid method: %s. Expected GET.", r.Method)
		return
	}

	query := r.URL.Query()
	startDateStr := query.Get("startDate")
	endDateStr := query.Get("endDate")

	var startDate, endDate *time.Time

	// Парсим дату начала
	if startDateStr != "" {
		t, err := time.Parse("02.01.2006", startDateStr)
		if err != nil {
			response := models.ErrorResponse{Error: "Invalid startDate format. Use DD.MM.YYYY"}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}
		startDate = &t
	}

	// Парсим дату конца
	if endDateStr != "" {
		t, err := time.Parse("02.01.2006", endDateStr)
		if err != nil {
			response := models.ErrorResponse{Error: "Invalid endDate format. Use DD.MM.YYYY"}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}
		endDate = &t
	}

	// Получаем данные из репозитория
	dbConn, err := db.InitDB()
	if err != nil {
		log.Printf("Failed to connect to DB: %v", err)
		http.Error(w, "Failed to connect to the database", http.StatusInternalServerError)
		return
	}
	defer dbConn.Close()

	items, err := repositories.NumberOfOrderedItems(dbConn, startDate, endDate)
	if err != nil {
		response := models.ErrorResponse{Error: err.Error()}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Если данных нет
	if len(items) == 0 {
		response := models.ErrorResponse{Error: "No data found for the given date range"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Создаем результат
	result := make(map[string]int)
	for _, item := range items {
		result[item.Name] = item.Quantity
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
