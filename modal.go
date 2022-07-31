package main

import (
	"github.com/jnsougata/disgo/bot"
)

type TextInput struct {
	CustomId    string `json:"custom_id"`   // filled internally
	Label       string `json:"label"`       // required default: "Text Input"
	Style       int    `json:"style"`       // 1 for short, 2 for long default: 1
	Value       string `json:"value"`       // default: ""
	Placeholder string `json:"placeholder"` // max 100 chars
	MinLength   int    `json:"min_length"`  // default: 0 upto 4000
	MaxLength   int    `json:"max_length"`  // default: 0 upto 4000
	Required    bool   `json:"required"`    // default: false
}

func (inp *TextInput) ToComponent() map[string]interface{} {
	inp.CustomId = AssignId(inp.CustomId)
	field := map[string]interface{}{
		"type":      4,
		"custom_id": inp.CustomId,
	}
	if inp.Label != "" {
		field["label"] = inp.Label
	} else {
		field["label"] = "Text Input"
	}
	if inp.Style != 0 {
		field["style"] = inp.Style
	} else {
		field["style"] = 1
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
	} else if inp.MaxLength > 4000 {
		field["max_length"] = 4000
	} else {
		field["max_length"] = 1
	}
	if inp.Required {
		field["required"] = true
	}
	return field
}

type Modal struct {
	Title       string
	CustomId    string // filled internally
	Fields      []TextInput
	SelectMenus []SelectMenu
}

func (m *Modal) OnSubmit(handler func(bot bot.User, ctx Context)) {
	m.CustomId = AssignId(m.CustomId)
	CallbackTasks[m.CustomId] = handler
}

func (m *Modal) Marshal() map[string]interface{} {
	modal := map[string]interface{}{}
	modal["title"] = m.Title
	modal["custom_id"] = AssignId(m.CustomId)
	modal["components"] = []map[string]interface{}{}
	if len(m.Fields) > 0 {
		for _, field := range m.Fields {
			row := map[string]interface{}{
				"type":       1,
				"components": []map[string]interface{}{},
			}
			row["components"] = append(row["components"].([]map[string]interface{}), field.ToComponent())
			modal["components"] = append(modal["components"].([]map[string]interface{}), row)
		}
	}
	if len(m.SelectMenus) > 0 {
		for _, menu := range m.SelectMenus {
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
