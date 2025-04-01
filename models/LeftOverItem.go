package models

type LeftOverItem struct {
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

type LeftOverResponse struct {
	CurrentPage int            `json:"currentPage"`
	HasNextPage bool           `json:"hasNextPage"`
	PageSize    int            `json:"pageSize"`
	TotalPages  int            `json:"totalPages"`
	Data        []LeftOverItem `json:"data"`
}
