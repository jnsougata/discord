package discord

import "errors"

type EmbedImage struct {
	URL      string `json:"url"`
	Height   int    `json:"height"`
	Width    int    `json:"width"`
	ProxyURL string `json:"proxy_url"`
}

type EmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

type EmbedFooter struct {
	Text    string `json:"text"`
	IconURL string `json:"icon_url"`
}

type EmbedAuthor struct {
	Name    string `json:"name"`
	IconURL string `json:"icon_url"`
}

type Embed struct {
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Url         string       `json:"url"`
	Timestamp   string       `json:"timestamp"`
	Color       int          `json:"color"`
	Footer      EmbedFooter  `json:"footer"`
	Author      EmbedAuthor  `json:"author"`
	Image       EmbedImage   `json:"image"`
	Thumbnail   EmbedImage   `json:"thumbnail"`
	Fields      []EmbedField `json:"fields"`
}

func (e *Embed) marshal() (map[string]interface{}, error) {
	embed := map[string]interface{}{}
	if e.Title != "" && len(e.Title) <= 256 {
		embed["title"] = e.Title
	} else {
		return nil, errors.New("`Title` can not be empty or longer than 256 characters")
	}
	if e.Description != "" && len(e.Description) <= 4096 {
		embed["description"] = e.Description
	} else {
		return nil, errors.New("`Description` can not be empty or longer than 4096 characters")
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
	} else {
		return nil, errors.New("`Footer.Text` can not be empty or longer than 2048 characters")
	}
	if e.Footer.IconURL != "" {
		embed["footer"].(map[string]interface{})["icon_url"] = e.Footer.IconURL
	}
	embed["author"] = map[string]interface{}{}
	if e.Author.Name != "" && len(e.Author.Name) <= 256 {
		embed["author"].(map[string]interface{})["name"] = e.Author.Name
	} else {
		return nil, errors.New("`Author.Name` can not be empty or longer than 256 characters")
	}
	if e.Author.IconURL != "" {
		embed["author"].(map[string]interface{})["icon_url"] = e.Author.IconURL
	}
	embed["image"] = map[string]interface{}{}
	if e.Image.URL != "" {
		embed["image"].(map[string]interface{})["url"] = e.Image.URL
	}
	if e.Image.Height > 0 {
		embed["image"].(map[string]interface{})["height"] = e.Image.Height
	}
	if e.Image.Width > 0 {
		embed["image"].(map[string]interface{})["width"] = e.Image.Width
	}
	embed["thumbnail"] = map[string]interface{}{}
	if e.Thumbnail.URL != "" {
		embed["thumbnail"].(map[string]interface{})["url"] = e.Thumbnail.URL
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
	embed["fields"] = []EmbedField{}
	for _, field := range e.Fields {
		if len(field.Name) <= 256 && len(field.Value) <= 1024 {
			embed["fields"] = append(embed["fields"].([]EmbedField), field)
		} else {
			return nil, errors.New(
				"`Field.Name` and `Field.Value` " +
					"can not be empty or longer than 256 and 1024 characters respectively")
		}
	}
	return embed, nil
}
