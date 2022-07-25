package domain

// ContentItemCategory maps the schema for the table that stores the content item category
type ContentItemCategory struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	IconURL string `json:"icon_url"`
}
