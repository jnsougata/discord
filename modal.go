package discord

import "errors"

type TextInput struct {
	CustomId    string `json:"custom_id"`
	Label       string `json:"label"`       // required default: "Text Input"
	Style       int    `json:"style"`       // 1 for short, 2 for long default: 1
	Value       string `json:"value"`       // default: ""
	Placeholder string `json:"placeholder"` // max 100 chars
	MinLength   int    `json:"min_length"`  // default: 0 upto 4000
	MaxLength   int    `json:"max_length"`  // default: 0 upto 4000
	Required    bool   `json:"required"`    // default: false
}

func (inp *TextInput) marshal() (map[string]interface{}, error) {
	if inp.CustomId == "" {
		return nil, errors.New("`CustomId` can not be empty for TextInput")
	}
	field := map[string]interface{}{
		"type":      4,
		"custom_id": inp.CustomId,
	}
	if inp.Label != "" {
		field["label"] = inp.Label
	} else {
		return nil, errors.New("`Label` can not be empty for TextInput")
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
	} else {
		field["placeholder"] = inp.Placeholder[:100]
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
	return field, nil
}

type Modal struct {
	customId    string
	Title       string
	Fields      []TextInput
	SelectMenus []SelectMenu
}

func (m *Modal) OnSubmit(handler func(bot Bot, ctx Context)) {
	m.customId = assignId()
	callbackTasks[m.customId] = handler
}

func (m *Modal) marshal() (map[string]interface{}, error) {
	modal := map[string]interface{}{}
	modal["title"] = m.Title
	if m.customId != "" {
		modal["custom_id"] = m.customId
	} else {
		modal["custom_id"] = assignId()
	}
	modal["components"] = []map[string]interface{}{}
	if len(m.Fields) > 0 {
		for _, field := range m.Fields {
			row := map[string]interface{}{
				"type":       1,
				"components": []map[string]interface{}{},
			}
			fieldValue, err := field.marshal()
			if err == nil {
				row["components"] = append(row["components"].([]map[string]interface{}), fieldValue)
				modal["components"] = append(modal["components"].([]map[string]interface{}), row)
			} else {
				return nil, err
			}
		}
	}
	if len(m.SelectMenus) > 0 {
		for _, menu := range m.SelectMenus {
			row := map[string]interface{}{
				"type":       1,
				"components": []map[string]interface{}{},
			}
			row["components"] = append(row["components"].([]map[string]interface{}), menu.marshal())
			modal["components"] = append(modal["components"].([]map[string]interface{}), row)
		}
	}
	return map[string]interface{}{"type": 9, "data": modal}, nil
}
