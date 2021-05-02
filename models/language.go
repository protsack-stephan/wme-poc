package models

type Language struct {
	ID            string `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	AlternateName string `json:"alternateName,omitempty"`
	Identifier    string `json:"identifier,omitempty"`
}
