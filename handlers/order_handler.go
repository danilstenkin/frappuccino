package handlers

import (
	"encoding/json"
	"frappuccino/models"
	"frappuccino/repositories"
	"net/http"
	"strings"
)

func GetOrdersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	orders, err := repositories.GetOrders()
	if err != nil {
		http.Error(w, "Ошибка при получении заказов: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var order models.Order

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&order)
	if err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if order.CustomerID == 0 {
		http.Error(w, "Неверные данные заказа", http.StatusBadRequest)
		return
	}

	id, err := repositories.CreateOrder(order)
	if err != nil {
		http.Error(w, "Ошибка при создании заказа: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := map[string]int{"id": id}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func GetOrderByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/orders/")

	if id == "" {
		http.Error(w, "ID not found", http.StatusBadRequest)
		return
	}

	order, err := repositories.GetOrderById(id)
	if err != nil {
		http.Error(w, "Ошибка при получении заказа: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

func UpdateOrderHandler(w http.ResponseWriter, r *http.Request) {
	idAndStatusTrim := strings.TrimPrefix(r.URL.Path, "/orders/")
	idStatus := strings.Split(idAndStatusTrim, "/")
	if len(idStatus) != 2 {
		http.Error(w, "Request should contain id and status", http.StatusBadRequest)
		return
	}
	id := idStatus[0]
	status := idStatus[1]
	if id == "" {
		http.Error(w, "ID не указан", http.StatusBadRequest)
		return
	}

	// Ожидаем {"status": "completed"}
	var data struct {
		Status string `json:"status"`
	}
	data.Status = status

	err := repositories.UpdateOrderStatus(id, data.Status)
	if err != nil {
		http.Error(w, "Ошибка при обновлении: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Статус обновлён и история записана"}`))
}
