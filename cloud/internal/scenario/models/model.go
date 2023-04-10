package models

type MakeCoffeeRequest struct {
	ID   string     `json:"id"`
	Type CoffeeType `json:"type"`

	Sugar uint `json:"sugar"`
}
