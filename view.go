package disgo

import (
	"log"
)

var callbackTasks = map[string]interface{}{}
var timeoutTasks = map[string][]interface{}{}

type Button struct {
	Style    int    // default: 1 (blue) More: 2 (grey), 3 (green), 4 (red), 5 (link)
	Label    string // default: "Button"
	Emoji    PartialEmoji
	URL      string // only for style 5 (link)
	Disabled bool
	CustomId string // filled internally
	OnClick  func(bot BotUser, comp Context)
}

func (b *Button) Marshal() map[string]interface{} {
	b.CustomId = AssignId("")
	if b.OnClick != nil {
		callbackTasks[b.CustomId] = b.OnClick
	}
	btn := map[string]interface{}{
		"type":      2,
		"custom_id": b.CustomId,
	}
	if b.Style != 0 {
		btn["style"] = b.Style
	} else {
		btn["style"] = 1
	}
	if b.Label != "" {
		btn["label"] = b.Label
	} else {
		btn["label"] = "Button"
	}
	if b.Emoji.Id != "" {
		btn["emoji"] = b.Emoji
	}
	if b.URL != "" && b.Style == 5 {
		btn["url"] = b.URL
	}
	if b.Disabled {
		btn["disabled"] = true
	}
	return btn
}

type SelectOption struct {
	Label       string // max 100 characters
	Value       string // default: ""
	Description string // max 100 characters
	Emoji       PartialEmoji
	Default     bool // default: false
}

func (so *SelectOption) ToComponent() map[string]interface{} {
	op := map[string]interface{}{}
	if so.Label != "" && len(so.Label) <= 100 {
		op["label"] = so.Label
	} else {
		panic("Name of the option can contain max 100 characters and must not be empty")
	}
	op["value"] = so.Value
	if len(so.Description) <= 100 {
		op["description"] = so.Description
	} else {
		panic("Description of the option can contain max 100 characters")
	}
	if so.Emoji.Id != "" {
		op["emoji"] = so.Emoji
	}
	if so.Default {
		op["default"] = true
	}
	return op
}

type SelectMenu struct {
	CustomId    string         // filled internally
	Options     []SelectOption // max 25 options
	Placeholder string         // max 100 characters
	MinValues   int            // default: 0
	MaxValues   int            // default: 1
	Disabled    bool
	OnSelection func(bot BotUser, comp Context, values ...string)
}

func (s *SelectMenu) ToComponent() map[string]interface{} {
	s.CustomId = AssignId("")
	if s.OnSelection != nil {
		callbackTasks[s.CustomId] = s.OnSelection
	}
	menu := map[string]interface{}{"type": 3, "custom_id": s.CustomId}
	if s.Placeholder != "" {
		menu["placeholder"] = s.Placeholder
	}
	if s.MinValues != 0 && s.MinValues > 25 {
		s.MinValues = 25
	}
	if s.MinValues < 0 {
		s.MinValues = 0
	}
	if s.MaxValues != 0 && s.MaxValues > 25 {
		s.MaxValues = 25
	}
	if s.MaxValues < 0 || s.MaxValues == 0 {
		s.MaxValues = 1
	}
	menu["min_values"] = s.MinValues
	menu["max_values"] = s.MaxValues
	if s.Disabled {
		menu["disabled"] = true
	}
	if len(s.Options) > 25 {
		s.Options = s.Options[:25]
	}
	menu["options"] = []map[string]interface{}{}
	for _, option := range s.Options {
		menu["options"] = append(menu["options"].([]map[string]interface{}), option.ToComponent())
	}
	return menu
}

type ActionRow struct {
	Buttons    []Button // max 5 buttons
	SelectMenu SelectMenu
}

type View struct {
	Timeout    float64     // default: 15 * 60 seconds
	ActionRows []ActionRow // max 5 rows
	OnTimeout  func(bot BotUser, ctx Context)
}

func (v *View) AddRow(row ActionRow) {
	if len(v.ActionRows) < 5 {
		v.ActionRows = append(v.ActionRows, row)
	}
}

func (v *View) AddButtons(buttons ...Button) {
	if len(v.ActionRows) < 5 {
		row := ActionRow{Buttons: buttons}
		v.ActionRows = append([]ActionRow{row}, v.ActionRows...)
	}
}

func (v *View) AddSelectMenu(menu SelectMenu) {
	if len(v.ActionRows) < 5 {
		row := ActionRow{SelectMenu: menu}
		v.ActionRows = append([]ActionRow{row}, v.ActionRows...)
	}
}

func (v *View) ToComponent() []interface{} {
	const timeout = 14.98 * 60
	if v.Timeout == 0 || v.Timeout > timeout {
		v.Timeout = timeout
	}
	var undo = map[string]bool{}
	var c []interface{}
	if len(v.ActionRows) > 5 {
		v.ActionRows = v.ActionRows[:5]
	}
	for _, row := range v.ActionRows {
		num := 0
		tmp := map[string]interface{}{
			"type":       1,
			"components": []interface{}{},
		}
		for _, button := range row.Buttons {
			if num < 5 {
				undo[button.CustomId] = true
				if v.OnTimeout != nil {
					timeoutTasks[button.CustomId] = []interface{}{v.Timeout, v.OnTimeout}
				}
				tmp["components"] = append(tmp["components"].([]interface{}), button.Marshal())
				num++
			}
		}
		if len(row.SelectMenu.Options) > 0 {
			if num == 0 {
				undo[row.SelectMenu.CustomId] = true
				if v.OnTimeout != nil {
					timeoutTasks[row.SelectMenu.CustomId] = []interface{}{v.Timeout, v.OnTimeout}
				}
				tmp["components"] = append(tmp["components"].([]interface{}), row.SelectMenu.ToComponent())
			} else {
				log.Println("Single ActionRow can contain either 1x SelectMenu or max 5x Buttons")
			}
		}
		if len(undo) > 0 {
			c = append(c, tmp)
			go ScheduleDeletion(v.Timeout, callbackTasks, undo)
		}
	}
	return c
}
