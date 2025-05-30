package models

type InventoryItem struct {
	ID           int     `json:"-"`
	Name         string  `json:"name"`
	Quantity     int     `json:"quantity"`
	Unit         string  `json:"unit"`
	PricePerUnit float64 `json:"price_per_unit"`
	LastUpdated  string  `json:"-"`
}

type InventoryItemResponce struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Quantity     int     `json:"quantity"`
	Unit         string  `json:"unit"`
	PricePerUnit float64 `json:"price_per_unit"`
	LastUpdated  string  `json:"-"`
}
