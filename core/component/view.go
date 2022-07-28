package component

import (
	"github.com/jnsougata/disgo/core/emoji"
	"github.com/jnsougata/disgo/core/user"
	"github.com/jnsougata/disgo/core/utils"
	"log"
)

const (
	BlueButton  = 1
	GreyButton  = 2
	GreenButton = 3
	RedButton   = 4
	LinkButton  = 5
)

var CallbackTasks = map[string]interface{}{}
var TimeoutTasks = map[string][]interface{}{}

type Button struct {
	Style    int
	Label    string
	Emoji    emoji.Partial
	URL      string
	Disabled bool
	CustomId string
	OnClick  func(bot user.User, cc Context)
}

func (b *Button) ToComponent() map[string]interface{} {
	b.CustomId = utils.AssignId("")
	if b.OnClick != nil {
		CallbackTasks[b.CustomId] = b.OnClick
	}
	btn := map[string]interface{}{
		"type":      2,
		"custom_id": b.CustomId,
	}
	if b.Style != 0 {
		btn["style"] = b.Style
	} else {
		btn["style"] = BlueButton
	}
	if b.Label != "" {
		btn["label"] = b.Label
	} else {
		btn["label"] = "Button"
	}
	if b.Emoji.Id != "" {
		btn["emoji"] = b.Emoji
	}
	if b.URL != "" && b.Style == LinkButton {
		btn["url"] = b.URL
	}
	if b.Disabled {
		btn["disabled"] = true
	}
	return btn
}

type SelectOption struct {
	Label       string
	Value       string
	Description string
	Emoji       emoji.Partial
	Default     bool
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
	Type        int
	CustomId    string
	Options     []SelectOption
	Placeholder string
	MinValues   int
	MaxValues   int
	Disabled    bool
	OnSelection func(bot user.User, cc Context, values ...string)
}

func (s *SelectMenu) ToComponent() map[string]interface{} {
	s.CustomId = utils.AssignId("")
	if s.OnSelection != nil {
		CallbackTasks[s.CustomId] = s.OnSelection
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
	Buttons    []Button
	SelectMenu SelectMenu
}

type View struct {
	Timeout    float64
	ActionRows []ActionRow
	OnTimeout  func(bot user.User, interaction Context)
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
					TimeoutTasks[button.CustomId] = []interface{}{v.Timeout, v.OnTimeout}
				}
				tmp["components"] = append(tmp["components"].([]interface{}), button.ToComponent())
				num++
			}
		}
		if len(row.SelectMenu.Options) > 0 {
			if num == 0 {
				undo[row.SelectMenu.CustomId] = true
				if v.OnTimeout != nil {
					TimeoutTasks[row.SelectMenu.CustomId] = []interface{}{v.Timeout, v.OnTimeout}
				}
				tmp["components"] = append(tmp["components"].([]interface{}), row.SelectMenu.ToComponent())
			} else {
				log.Println("Single ActionRow can contain either 1x SelectMenu or max 5x Buttons")
			}
		}
		if len(undo) > 0 {
			c = append(c, tmp)
			go utils.ScheduleDeletion(v.Timeout, CallbackTasks, undo)
		}
	}
	return c
}
