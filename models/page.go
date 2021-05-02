package models

type Page struct {
	ID             string     `json:"id,omitempty"`
	Name           string     `json:"name,omitempty"`
	Identifier     int        `json:"identifier,omitempty"`
	Version        int        `json:"version,omitempty"`
	DateModified   string     `json:"dateModified,omitempty"`
	URL            string     `json:"url,omitempty"`
	Namespace      *Namespace `json:"namespace,omitempty"`
	InLanguage     *Language  `json:"inLanguage,omitempty"`
	MainEntity     *QID       `json:"mainEntity,omitempty"`
	ArticleBody    string     `json:"articleBody,omitempty"`
	EncodingFormat string     `json:"encodingFormat,omitempty"`
	IsPartOf       *Project   `json:"isPartOf,omitempty"`
}
