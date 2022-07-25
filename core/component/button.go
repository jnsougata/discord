package component

import "github.com/jnsougata/disgo/core/emoji"

type Button struct {
	Type     string        `json:"type"`
	Style    string        `json:"style"`
	Label    string        `json:"label"`
	Emoji    emoji.Partial `json:"emoji"`
	CustomID string        `json:"custom_id"`
	URL      string        `json:"url"`
	Disabled bool          `json:"disabled"`
}
