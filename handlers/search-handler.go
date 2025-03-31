package handlers

import (
	"encoding/json"
	"frappuccino/db"
	"frappuccino/models"
	"frappuccino/repositories"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func SearchReportsHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" {
		http.Error(w, `{"error": "Query (q) is required"}`, http.StatusBadRequest)
		return
	}

	filters := strings.Split(r.URL.Query().Get("filter"), ",")
	if len(filters) == 1 && filters[0] == "" {
		filters = []string{"all"}
	}

	minPrice, _ := strconv.ParseFloat(r.URL.Query().Get("minPrice"), 64)
	maxPrice, _ := strconv.ParseFloat(r.URL.Query().Get("maxPrice"), 64)
	if maxPrice == 0 {
		maxPrice = 1e6 // default upper limit
	}

	dbConn, err := db.InitDB()
	if err != nil {
		log.Println("DB connection error:", err)
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		return
	}
	defer dbConn.Close()

	var result models.SearchResponse
	for _, f := range filters {
		switch f {
		case "menu":
			items, err := repositories.SearchMenuItems(dbConn, q, minPrice, maxPrice)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			result.MenuItems = items
			result.TotalMatches += len(items)
		case "orders":
			orders, err := repositories.SearchOrders(dbConn, q, minPrice, maxPrice)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			result.Orders = orders
			result.TotalMatches += len(orders)
		case "all":
			items, _ := repositories.SearchMenuItems(dbConn, q, minPrice, maxPrice)
			orders, _ := repositories.SearchOrders(dbConn, q, minPrice, maxPrice)
			result.MenuItems = items
			result.Orders = orders
			result.TotalMatches = len(items) + len(orders)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
