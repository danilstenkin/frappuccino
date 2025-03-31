package handlers

import (
	"encoding/json"
	"frappuccino/db"
	"frappuccino/models"
	"frappuccino/repositories"
	"net/http"
)

func GetOrderedItemsByPeriod(w http.ResponseWriter, r *http.Request) {
	period := r.URL.Query().Get("period")
	month := r.URL.Query().Get("month")
	year := r.URL.Query().Get("year")

	if period != "day" && period != "month" {
		http.Error(w, `{"error": "Invalid or missing period. Use 'day' or 'month'."}`, http.StatusBadRequest)
		return
	}
	dbConn, err := db.InitDB()
	if err != nil {
		http.Error(w, `{"error": "Failed to connect to DB"}`, http.StatusInternalServerError)
		return
	}
	defer dbConn.Close()

	var response models.OrderedItemByPeriodResponse
	response.Period = period

	if period == "day" {
		if month == "" {
			http.Error(w, `{"error": "'month' parameter is required for period=day"}`, http.StatusBadRequest)
			return
		}
		response.Month = month
		data, err := repositories.GetOrderedItemsByDay(dbConn, month)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		response.OrderedItems = data
	} else if period == "month" {
		if year == "" {
			year = "2025"
		}
		response.Year = year
		data, err := repositories.GetOrderedItemsByMonth(dbConn, year)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		response.OrderedItems = data
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
