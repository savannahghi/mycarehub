package domain

// Content aggregates all content details into one payload that is returned from an API and
// rendered on the front end
type Content struct {
	Meta  Meta          `json:"meta"`
	Items []ContentItem `json:"items"`
}

// Meta holds the information that shows the total count of items returned from the API
type Meta struct {
	TotalCount int `json:"total_count"`
}

// ContentItem holds all the information necessary relating to content item
type ContentItem struct {
	ID                  int                `json:"id"`
	Meta                ContentMeta        `json:"meta"`
	Title               string             `json:"title"`
	Date                string             `json:"date"`
	Intro               string             `json:"intro"`
	Author              Author             `json:"author"`
	AuthorName          string             `json:"author_name"`
	ItemType            string             `json:"item_type"`
	TimeEstimateSeconds int                `json:"time_estimate_seconds"`
	Body                string             `json:"body"`
	TagNames            []string           `json:"tag_names"`
	HeroImage           HeroImage          `json:"hero_image"`
	HeroImageRendition  HeroImageRendition `json:"hero_image_rendition"`
	LikeCount           int                `json:"like_count"`
	BookmarkCount       int                `json:"bookmark_count"`
	ViewCount           int                `json:"view_count"`
	ShareCount          int                `json:"share_count"`
	Documents           []Documents        `json:"documents"`
	CategoryDetails     []CategoryDetails  `json:"category_details"`
}

// HeroImage ...
type HeroImage struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

// HeroImageRendition ...
type HeroImageRendition struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Alt    string `json:"alt"`
}

// CategoryDetails holds all information regarding a certain item's category.
// It will be used to determine which content to show on the UI based on the category
// ID provided to the API
type CategoryDetails struct {
	ID           int    `json:"category_id"`
	CategoryName string `json:"category_name"`
	CategoryIcon string `json:"category_icon"`
}

// Documents ...
type Documents struct {
	ID       int          `json:"id"`
	Meta     AuthorMeta   `json:"meta"`
	Document DocumentData `json:"document"`
}

// DocumentData holds the information regarding a document
type DocumentData struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

// ContentMeta ...
type ContentMeta struct {
	ContentType       string `json:"type"`
	ContentDetailURL  string `json:"detail_url"`
	ContentHTMLURL    string `json:"html_url"`
	Slug              string `json:"slug"`
	ShowInMenus       bool   `json:"show_in_menus"`
	SEOTitle          string `json:"seo_title"`
	SearchDescription string `json:"search_description"`
	FirstPublishedAt  string `json:"first_published_at"`
	Locale            string `json:"locale"`
}

// Author ...
type Author struct {
	ID         string     `json:"id"`
	Meta       AuthorMeta `json:"meta"`
	AuthorName string     `json:"author_name"`
}

// AuthorMeta ...
type AuthorMeta struct {
	Type string `json:"type"`
}
