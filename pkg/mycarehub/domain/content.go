package domain

// Content aggregates all content details into one payload that is returned from an API and
// rendered on the front end
type Content struct {
	Meta  Meta          `json:"meta"`
	Items []ContentItem `json:"items"`
}

// Meta holds the information that shows the total count of items returned from the API
// The total count displayed is irrespective of pagination
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
	Documents           []Document         `json:"documents"`
	CategoryDetails     []CategoryDetail   `json:"category_details"`
	FeaturedMedia       []FeaturedMedia    `json:"featured_media"`
	GalleryImages       []GalleryImage     `json:"gallery_images"`
}

// GalleryImage contains details about images that can be featured on a gallery
type GalleryImage struct {
	ID    int         `json:"id"`
	Image ImageDetail `json:"image"`
}

// ImageDetail contains more information about an image
type ImageDetail struct {
	ID    int       `json:"id"`
	Title string    `json:"title"`
	Meta  ImageMeta `json:"meta"`
}

// ImageMeta holds more information about an Image
type ImageMeta struct {
	Type             string `json:"type"`
	ImageDetailURL   string `json:"detail_url"`
	ImageDownloadURL string `json:"download_url"`
}

// FeaturedMedia ...
type FeaturedMedia struct {
	ID        int     `json:"id"`
	URL       string  `json:"url"`
	Title     string  `json:"title"`
	Type      string  `json:"type"`
	Width     int     `json:"width"`
	Height    int     `json:"height"`
	Thumbnail string  `json:"thumbnail"`
	Duration  float64 `json:"duration"`
}

// HeroImage contains details about the hero image i.e The title
// This is the oversized image displayed at the top of a content
type HeroImage struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

// HeroImageRendition contains more details about the hero image. These details will be used
// by the frontend to get the actual image and render it on the app
type HeroImageRendition struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Alt    string `json:"alt"`
}

// CategoryDetail holds all information regarding a certain item's category.
// It will be used to determine which content to show on the UI based on the category
// ID provided to the API
type CategoryDetail struct {
	ID           int    `json:"category_id"`
	CategoryName string `json:"category_name"`
	CategoryIcon string `json:"category_icon"`
}

// Document contains details about a document eg a PDF file
type Document struct {
	ID       int          `json:"id"`
	Meta     DocumentMeta `json:"meta"`
	Document DocumentData `json:"document"`
}

// DocumentMeta represents a list of properties that are associated with the document
type DocumentMeta struct {
	Type                string `json:"type"`
	DocumentDetailURL   string `json:"detail_url"`
	DocumentDownloadURL string `json:"download_url"`
}

// DocumentData holds the information regarding a document
type DocumentData struct {
	ID    int          `json:"id"`
	Title string       `json:"title"`
	Meta  DocumentMeta `json:"meta"`
}

// ContentMeta represents a list of properties that are associated with the content model
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

// Author models the details about an author
type Author struct {
	ID   string     `json:"id"`
	Meta AuthorMeta `json:"meta"`
}

// AuthorMeta holds the properties that are associated to the author model
type AuthorMeta struct {
	Type string `json:"type"`
}
