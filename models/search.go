package models

type SearchMenuItem struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Relevance   float64 `json:"relevance"`
}

type SearchOrder struct {
	ID           int      `json:"id"`
	CustomerName string   `json:"customer_name"`
	Items        []string `json:"items"`
	Total        float64  `json:"total"`
	Relevance    float64  `json:"relevance"`
}

type SearchResponse struct {
	MenuItems    []SearchMenuItem `json:"menu_items,omitempty"`
	Orders       []SearchOrder    `json:"orders,omitempty"`
	TotalMatches int              `json:"total_matches"`
}
