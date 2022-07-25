package embed

import "encoding/json"

type Image struct {
	Url    string `json:"url"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

type Field struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

type Footer struct {
	Text    string `json:"text"`
	IconUrl string `json:"icon_url"`
}

type Author struct {
	Name    string `json:"name"`
	IconUrl string `json:"icon_url"`
}

type Embed struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Url         string  `json:"url"`
	Timestamp   string  `json:"timestamp"`
	Color       int     `json:"color"`
	Footer      Footer  `json:"footer"`
	Author      Author  `json:"author"`
	Image       Image   `json:"image"`
	Thumbnail   Image   `json:"thumbnail"`
	Fields      []Field `json:"fields"`
}

func FromData(payload interface{}) *Embed {
	embed := &Embed{}
	data, _ := json.Marshal(payload)
	_ = json.Unmarshal(data, embed)
	return embed
}
