package discord

import (
	"errors"
	"strings"
)

type embedImage struct {
	URL      string `json:"url"`
	Height   int    `json:"height"`
	Width    int    `json:"width"`
	ProxyURL string `json:"proxy_url"`
}

type embedField struct {
	Name   string `json:"name"`
	Value  string `json:"Value"`
	Inline bool   `json:"inline"`
}

type embedFooter struct {
	Text    string `json:"text"`
	IconURL string `json:"icon_url"`
}

type embedAuthor struct {
	Name    string `json:"name"`
	IconURL string `json:"icon_url"`
}

type Embed struct {
	Title       string       `json:"title"`
	Description string       `json:"description"`
	URL         string       `json:"url"`
	Timestamp   string       `json:"timestamp"`
	Color       int          `json:"color"`
	Footer      embedFooter  `json:"footer"`
	Author      embedAuthor  `json:"author"`
	Image       embedImage   `json:"image"`
	Thumbnail   embedImage   `json:"thumbnail"`
	Fields      []embedField `json:"fields"`
}

func (e *Embed) SetImage(url string) error {
	e.Image.URL = url
	if url != "" {
		if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") || strings.HasPrefix(url, "attachment://") {
			e.Image.URL = url
			return nil
		} else {
			return errors.New("`Image.URL` must be a valid URL of following schemas http:// or https:// or attachment://")
		}
	} else {
		return errors.New("`url` can not be empty for embed image")
	}
}

func (e *Embed) SetThumbnail(url string) error {
	if url != "" {
		if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") || strings.HasPrefix(url, "attachment://") {
			e.Thumbnail.URL = url
			return nil
		} else {
			return errors.New("`Thumbnail.URL` must be a valid URL of following schemas http:// or https:// or attachment://")
		}
	} else {
		return errors.New("`url` can not be empty for embed thumbnail")
	}
}

func (e *Embed) SetAuthor(name string, iconURL string) error {
	if len(name) <= 256 {
		e.Author.Name = name
	} else {
		return errors.New("`name` can not be longer than 256 characters")
	}
	if iconURL != "" {
		if strings.HasPrefix(iconURL, "attachment://") || strings.HasPrefix(iconURL, "https://") {
			e.Author.IconURL = iconURL
		} else {
			return errors.New("`iconURL` must be a valid URL of following schemas https:// or attachment://")
		}
	}
	return nil
}

func (e *Embed) SetFooter(text string, iconURL string) error {
	if len(text) <= 2048 {
		e.Footer.Text = text
	} else {
		return errors.New("`text` can not be longer than 2048 characters")
	}
	if iconURL != "" {
		if strings.HasPrefix(iconURL, "http://") || strings.HasPrefix(iconURL, "https://") || strings.HasPrefix(iconURL, "attachment://") {
			e.Footer.IconURL = iconURL
		} else {
			return errors.New("`iconURL` must be a valid URL of following schemas http:// or https:// or attachment://")
		}
	}
	return nil
}
func (e *Embed) SetField(name string, value string, inline bool) error {
	if len(e.Fields) >= 25 {
		return errors.New("single embed can not have more than 25 fields")
	}
	if len(name) > 256 {
		return errors.New("`field name` can not be longer than 256 characters")
	}
	if len(value) > 1024 {
		return errors.New("`field Value` can not be longer than 1024 characters")
	}
	if name == "" {
		return errors.New("`field name` can not be empty")
	}
	if value == "" {
		return errors.New("`field Value` can not be empty")
	}
	e.Fields = append(e.Fields, embedField{Name: name, Value: value, Inline: inline})
	return nil
}

func (e *Embed) marshal() (map[string]interface{}, error) {
	embed := map[string]interface{}{}
	if len(e.Title) <= 256 {
		embed["title"] = e.Title
	} else {
		return nil, errors.New("`Embed.Title` can not be longer than 256 characters")
	}
	if len(e.Description) <= 4096 {
		embed["description"] = e.Description
	} else {
		return nil, errors.New("`Embed.Description` can not be longer than 4096 characters")
	}
	if e.URL != "" {
		if strings.HasPrefix(e.URL, "http://") || strings.HasPrefix(e.URL, "https://") {
			embed["url"] = e.URL
		} else {
			return nil, errors.New("`Embed.URL` must be a valid URL of following schemas http:// or https://")
		}
	}
	embed["color"] = e.Color
	embed["timestamp"] = e.Timestamp
	embed["footer"] = map[string]interface{}{"text": e.Footer.Text, "icon_url": e.Footer.IconURL}
	embed["author"] = map[string]interface{}{"name": e.Author.Name, "icon_url": e.Author.IconURL}
	embed["image"] = map[string]interface{}{"url": e.Image.URL}
	embed["thumbnail"] = map[string]interface{}{"url": e.Thumbnail.URL}
	embed["fields"] = []map[string]interface{}{}
	for _, field := range e.Fields {
		embed["fields"] = append(
			embed["fields"].([]map[string]interface{}),
			map[string]interface{}{"name": field.Name, "value": field.Value, "inline": field.Inline})
	}
	return embed, nil
}
