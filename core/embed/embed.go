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

func (e *Embed) ToBody() map[string]interface{} {
	embed := map[string]interface{}{}
	if e.Title != "" && len(e.Title) <= 256 {
		embed["title"] = e.Title
	}
	if e.Description != "" && len(e.Description) <= 4096 {
		embed["description"] = e.Description
	}
	if e.Url != "" {
		embed["url"] = e.Url
	}
	if e.Timestamp != "" {
		embed["timestamp"] = e.Timestamp
	}
	embed["color"] = e.Color
	embed["footer"] = map[string]interface{}{}
	if e.Footer.Text != "" && len(e.Footer.Text) <= 2048 {
		embed["footer"].(map[string]interface{})["text"] = e.Footer.Text
	}
	if e.Footer.IconUrl != "" {
		embed["footer"].(map[string]interface{})["icon_url"] = e.Footer.IconUrl
	}
	embed["author"] = map[string]interface{}{}
	if e.Author.Name != "" && len(e.Author.Name) <= 256 {
		embed["author"].(map[string]interface{})["name"] = e.Author.Name
	}
	if e.Author.IconUrl != "" {
		embed["author"].(map[string]interface{})["icon_url"] = e.Author.IconUrl
	}
	embed["image"] = map[string]interface{}{}
	if e.Image.Url != "" {
		embed["image"].(map[string]interface{})["url"] = e.Image.Url
	}
	if e.Image.Height > 0 {
		embed["image"].(map[string]interface{})["height"] = e.Image.Height
	}
	if e.Image.Width > 0 {
		embed["image"].(map[string]interface{})["width"] = e.Image.Width
	}
	embed["thumbnail"] = map[string]interface{}{}
	if e.Thumbnail.Url != "" {
		embed["thumbnail"].(map[string]interface{})["url"] = e.Thumbnail.Url
	}
	if e.Thumbnail.Height > 0 {
		embed["thumbnail"].(map[string]interface{})["height"] = e.Thumbnail.Height
	}
	if e.Thumbnail.Width > 0 {
		embed["thumbnail"].(map[string]interface{})["width"] = e.Thumbnail.Width
	}
	if len(e.Fields) > 25 {
		e.Fields = e.Fields[:25]
	}
	embed["fields"] = []Field{}
	for _, field := range e.Fields {
		if len(field.Name) <= 256 && len(field.Value) <= 1024 {
			embed["fields"] = append(embed["fields"].([]Field), field)
		}
	}
	return embed
}

func FromData(payload interface{}) *Embed {
	embed := &Embed{}
	data, _ := json.Marshal(payload)
	_ = json.Unmarshal(data, embed)
	return embed
}
