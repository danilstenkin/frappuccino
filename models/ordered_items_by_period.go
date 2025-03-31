package models

type OrderedItemByPeriodResponse struct {
	Period       string             `json:"period"`
	Month        string             `json:"month,omitempty"`
	Year         string             `json:"year,omitempty"`
	OrderedItems []map[string]int64 `json:"orderedItems"`
}
