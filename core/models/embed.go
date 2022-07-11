package models

type Embed struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Url         string `json:"url"`
	Timestamp   string `json:"timestamp"`
	Color       int    `json:"color"`

	Footer struct {
		Text string `json:"text"`
		Icon string `json:"icon_url"`
	} `json:"footer"`

	Author struct {
		Name    string `json:"name"`
		IconUrl string `json:"icon_url"`
	} `json:"author"`

	Image struct {
		Url    string `json:"url"`
		Height int    `json:"height"`
		Width  int    `json:"width"`
	} `json:"image"`

	Thumbnail struct {
		Url    string `json:"url"`
		Height int    `json:"height"`
		Width  int    `json:"width"`
	} `json:"thumbnail"`

	Fields []struct {
		Name   string `json:"name"`
		Value  string `json:"value"`
		Inline bool   `json:"inline"`
	} `json:"fields"`
}
