package models

// Структура для представления заказанных товаров
type OrderedItem struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

// Структура для ошибки
type ErrorResponse struct {
	Error string `json:"error"`
}
