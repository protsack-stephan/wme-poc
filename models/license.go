package models

type License struct {
	ID         string `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	Identifier string `json:"identifier,omitempty"`
}
