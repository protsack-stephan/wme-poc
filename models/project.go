package models

type Project struct {
	ID         string       `json:"id,omitempty"`
	Name       string       `json:"name,omitempty"`
	Identifier string       `json:"identifier,omitempty"`
	URL        string       `json:"url,omitempty"`
	InLanguage *Language    `json:"inLanguage,omitempty"`
	IsPartOf   *Project     `json:"isPartOf,omitempty"`
	Namespaces []*Namespace `json:"namespaces,omitempty"`
}
