package component

import (
	"github.com/jnsougata/disgo/core/user"
	"github.com/jnsougata/disgo/core/utils"
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
	inp.CustomId = utils.AssignId(inp.CustomId)
	field := map[string]interface{}{
		"type":      4,
		"custom_id": inp.CustomId,
	}
	if inp.Label != "" {
		field["label"] = inp.Label
	}
	if inp.Style != 0 {
		field["style"] = inp.Style
	} else {
		field["style"] = ShortInput
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

func (m *Modal) OnSubmit(handler func(bot user.User, interaction Interaction)) {
	m.CustomId = utils.AssignId(m.CustomId)
	CallbackTasks[m.CustomId] = handler
}

func (m *Modal) ToBody() map[string]interface{} {
	modal := map[string]interface{}{}
	modal["title"] = m.Title
	modal["custom_id"] = utils.AssignId(m.CustomId)
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
