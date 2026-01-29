package models

type Employee struct {
	ID       int   `json:"id"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Address  string `json:"address"`
	IsActive *bool  `json:"is_active,omitempty"`
}
