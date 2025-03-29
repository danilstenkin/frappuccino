package models

type CustomersResponse struct {
	ID          int                    `json:"id"`
	Name        string                 `json:"name"`
	Preferences map[string]interface{} `json:"preferences"`
}

type Customers struct {
	Name        string                 `json:"name"`
	Preferences map[string]interface{} `json:"preferences"`
}
