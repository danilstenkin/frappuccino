package handlers

import (
	"encoding/json"
	"frappuccino/db"
	"frappuccino/models"
	"frappuccino/repositories"
	"net/http"
	"strconv"
)

func GetLeftOversHandler(w http.ResponseWriter, r *http.Request) {
	sortBy := r.URL.Query().Get("sortBy")
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")

	page := 1
	pageSize := 10
	var err error

	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			http.Error(w, `{"error": "Invalid page"}`, http.StatusBadRequest)
			return
		}
	}

	if pageSizeStr != "" {
		pageSize, err = strconv.Atoi(pageSizeStr)
		if err != nil || pageSize < 1 {
			http.Error(w, `{"error": "Invalid pageSize"}`, http.StatusBadRequest)
			return
		}
	}

	dbConn, err := db.InitDB()
	if err != nil {
		http.Error(w, `{"error": "DB connection failed"}`, http.StatusInternalServerError)
		return
	}
	defer dbConn.Close()

	items, totalPages, err := repositories.GetLeftOvers(dbConn, sortBy, page, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := models.LeftOverResponse{
		CurrentPage: page,
		PageSize:    pageSize,
		TotalPages:  totalPages,
		HasNextPage: page < totalPages,
		Data:        items,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
