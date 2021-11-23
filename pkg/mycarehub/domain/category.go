package domain

// ContentItemCategory maps the schema for the table that stores the content item category
type ContentItemCategory struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	IconURL string `json:"icon_id"`
}

// WagtailImage maps the schema for the table that stores the wagtail images
type WagtailImage struct {
	ID      int    `json:"icon_id"`
	IconURL string `json:"file"`
}
