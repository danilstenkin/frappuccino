package handlers

import (
	"encoding/json"
	"frappuccino/db"
	"frappuccino/models"
	"frappuccino/repositories"
	"frappuccino/utils" // Импортируем utils для проверки валидности
	"log"
	"net/http"
)

func CreateMenuItemHandler(w http.ResponseWriter, r *http.Request) {
	// Только POST-запросы
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var item models.MenuItem

	// Декодируем JSON из тела запроса в структуру MenuItem
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Валидация обязательных полей
	if item.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	if item.Price <= 0 {
		http.Error(w, "Price must be greater than 0", http.StatusBadRequest)
		return
	}

	// Проверка правильности размера (size) - он должен быть из перечисления item_size
	validSizes := []string{"small", "medium", "large"}
	if !utils.IsValidSize(validSizes, item.Size) {
		http.Error(w, "Invalid size", http.StatusBadRequest)
		return
	}

	// Подключаемся к базе данных
	dbConn, err := db.InitDB()
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	defer dbConn.Close()

	// Создаем новый элемент меню в базе данных
	id, err := repositories.CreateMenuItem(dbConn, item)
	if err != nil {
		http.Error(w, "Could not create menu item: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем успешный ответ с ID нового элемента
	response := map[string]int{"id": id}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func GetMenuItemsHandler(w http.ResponseWriter, r *http.Request) {
	// Шаг 1: Получаем данные из базы данных
	items, err := repositories.GetMenuItems()
	if err != nil {
		http.Error(w, "Не удалось получить элементы меню: "+err.Error(), http.StatusInternalServerError)
		log.Println("Ошибка при получении элементов меню:", err)
		return
	}

	// Шаг 2: Отправляем ответ с элементами меню в формате JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}
