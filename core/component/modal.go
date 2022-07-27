package component

import (
	"github.com/jnsougata/disgo/core/user"
	"log"
)

const (
	LongInput  = 2
	ShortInput = 1
)

type TextInput struct {
	CustomId    string `json:"custom_id"`
	Label       string `json:"label"`
	Style       int    `json:"style"`
	Value       string `json:"value"`
	Placeholder string `json:"placeholder"`
	MinLength   int    `json:"min_length"`
	MaxLength   int    `json:"max_length"`
	Required    bool   `json:"required"`
}

func (inp *TextInput) ToComponent() map[string]interface{} {
	if inp.CustomId == "" {
		log.Println("CustomId is required for each TextInput")
		return map[string]interface{}{}
	}
	if Ids[inp.CustomId] {
		log.Println("CustomId must be unique for each TextInput")
		return map[string]interface{}{}
	} else {
		Ids[inp.CustomId] = true
	}
	field := map[string]interface{}{
		"type":      4,
		"custom_id": inp.CustomId,
	}
	if inp.Label != "" {
		field["label"] = inp.Label
	}
	if inp.Style != 0 {
		field["style"] = inp.Style
	}
	if inp.Value != "" {
		field["value"] = inp.Value
	}
	if inp.Placeholder != "" && len(inp.Placeholder) <= 100 {
		field["placeholder"] = inp.Placeholder
	}
	field["min_length"] = inp.MinLength
	if inp.MaxLength > 0 && inp.MaxLength <= 4000 {
		field["max_length"] = inp.MaxLength
	}
	if inp.Required {
		field["required"] = true
	}
	return field
}

type Modal struct {
	Title       string
	CustomId    string
	Fields      []TextInput
	SelectMenus []SelectMenu
}

func (i *Modal) Callback(handler func(bot user.User, interaction Interaction)) {
	if i.CustomId == "" {
		log.Println("CustomId is required for each Modal")
		return
	}
	CallbackFactory[i.CustomId] = handler
}

func (i *Modal) ToBody() map[string]interface{} {
	modal := map[string]interface{}{}
	modal["title"] = i.Title
	modal["custom_id"] = i.CustomId
	modal["components"] = []map[string]interface{}{}
	if len(i.Fields) > 0 {
		for _, field := range i.Fields {
			row := map[string]interface{}{
				"type":       1,
				"components": []map[string]interface{}{},
			}
			row["components"] = append(row["components"].([]map[string]interface{}), field.ToComponent())
			modal["components"] = append(modal["components"].([]map[string]interface{}), row)
		}
	}
	if len(i.SelectMenus) > 0 {
		for _, menu := range i.SelectMenus {
			row := map[string]interface{}{
				"type":       1,
				"components": []map[string]interface{}{},
			}
			row["components"] = append(row["components"].([]map[string]interface{}), menu.ToComponent())
			modal["components"] = append(modal["components"].([]map[string]interface{}), row)
		}
	}
	return map[string]interface{}{"type": 9, "data": modal}
}
