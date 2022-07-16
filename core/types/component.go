package types

type Component struct {
}

type Button struct {
	Type     string       `json:"type"`
	Style    string       `json:"style"`
	Label    string       `json:"label"`
	Emoji    PartialEmoji `json:"emoji"`
	CustomID string       `json:"custom_id"`
	URL      string       `json:"url"`
	Disabled bool         `json:"disabled"`
}
