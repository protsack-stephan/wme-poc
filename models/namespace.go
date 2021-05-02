package models

type Namespace struct {
	ID         string    `json:"id,omitempty"`
	Name       string    `json:"name,omitempty"`
	Identifier int       `json:"identifier,omitempty"`
	InLanguage *Language `json:"inLanguage,omitempty"`
	IsPartOf   *Project  `json:"isPartOf,omitempty"`
}
