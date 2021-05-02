package models

type ArticleBody struct {
	ID      string `json:"id,omitempty"`
	HTML    string `json:"html,omitempty"`
	Wiktext string `json:"wikitext,omitempty"`
}
