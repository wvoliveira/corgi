package model

type Health struct {
	Required    bool   `json:"required"`
	Status      string `json:"status"`
	Component   string `json:"component"`
	Description string `json:"description"`
}
