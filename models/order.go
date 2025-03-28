package models

import "time"

type Order struct {
	ID                  int                    `json:"-"`
	CustomerID          int                    `json:"customer_id"`
	Status              string                 `json:"-"`
	SpecialInstructions map[string]interface{} `json:"special_instructions"`
	TotalAmount         float64                `json:"-"`
	OrderDate           time.Time              `json:"-"`
}

type OrderRespone struct {
	ID                  int                    `json:"id"`
	CustomerID          int                    `json:"customer_id"`
	Status              string                 `json:"status"`
	SpecialInstructions map[string]interface{} `json:"special_instructions"`
	TotalAmount         float64                `json:"total_amount"`
	OrderDate           time.Time              `json:"order_date"`
}
