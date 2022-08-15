package discord

import "errors"

type ButtonStyle int

const (
	BlueButton  ButtonStyle = 1
	GreyButton  ButtonStyle = 2
	GreenButton ButtonStyle = 3
	RedButton   ButtonStyle = 4
	Linkbutton  ButtonStyle = 5
)

var callbackTasks = map[string]interface{}{}
var timeoutTasks = map[string][]interface{}{}

type Button struct {
	customId string
	Style    ButtonStyle
	Label    string
	Emoji    PartialEmoji
	URL      string
	Disabled bool
	OnClick  func(bot Bot, comp Context)
}

func (b *Button) marshal() (map[string]interface{}, error) {
	b.customId = assignId()
	if b.OnClick != nil {
		callbackTasks[b.customId] = b.OnClick
	}
	btn := map[string]interface{}{
		"type":      2,
		"custom_id": b.customId,
	}
	if int(b.Style) != 0 {
		btn["style"] = int(b.Style)
	} else {
		btn["style"] = int(BlueButton)
	}
	if b.Label != "" {
		btn["label"] = b.Label
	} else {
		return nil, errors.New("button label can not be empty")
	}
	if b.Emoji.Id != "" {
		btn["emoji"] = b.Emoji
	}
	if b.URL != "" && b.Style == Linkbutton {
		btn["url"] = b.URL
	}
	if b.Disabled {
		btn["disabled"] = true
	}
	return btn, nil
}

type SelectOption struct {
	Label       string // max 100 characters
	Value       string // default: ""
	Description string // max 100 characters
	Emoji       PartialEmoji
	Default     bool // default: false
}

func (so *SelectOption) marshal() (map[string]interface{}, error) {
	op := map[string]interface{}{}
	if so.Label != "" && len(so.Label) <= 100 {
		op["label"] = so.Label
	} else {
		return nil, errors.New("name of the option can contain max 100 characters and must not be empty")
	}
	op["value"] = so.Value
	if len(so.Description) <= 100 {
		op["description"] = so.Description
	} else {
		return nil, errors.New("description of the select option can contain max 100 characters")
	}
	if so.Emoji.Id != "" {
		op["emoji"] = so.Emoji
	}
	if so.Default {
		op["default"] = true
	}
	return op, nil
}

type SelectMenu struct {
	customId    string
	Options     []SelectOption // max 25 options
	Placeholder string         // max 100 characters
	MinValues   int            // default: 0
	MaxValues   int            // default: 1
	Disabled    bool
	OnSelection func(bot Bot, comp Context, values ...string)
}

func (s *SelectMenu) marshal() (map[string]interface{}, error) {
	s.customId = assignId()
	if s.OnSelection != nil {
		callbackTasks[s.customId] = s.OnSelection
	}
	menu := map[string]interface{}{"type": 3, "custom_id": s.customId}
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
		op, err := option.marshal()
		if err != nil {
			return nil, err
		}
		menu["options"] = append(menu["options"].([]map[string]interface{}), op)
	}
	return menu, nil
}

type actionRow struct {
	Buttons    []Button // max 5 buttons
	SelectMenu SelectMenu
}

type View struct {
	rows      []actionRow
	Timeout   float64 // default: 15 * 60 seconds
	OnTimeout func(bot Bot, ctx Context)
}

func (v *View) AddButtons(buttons ...Button) error {
	if len(v.rows) < 5 {
		if len(buttons) <= 5 {
			row := actionRow{Buttons: buttons}
			v.rows = append(v.rows, row)
			return nil
		} else {
			row := actionRow{Buttons: buttons[:5]}
			v.rows = append(v.rows, row)
			return errors.New("you can add max 5 buttons at a time")
		}
	} else {
		return errors.New("view can contain max 5 rows")
	}
}

func (v *View) AddSelectMenu(menu SelectMenu) error {
	if len(v.rows) < 5 {
		if len(menu.Options) == 0 {
			return errors.New("select menu must have at least one option")
		}
		row := actionRow{SelectMenu: menu}
		v.rows = append(v.rows, row)
		return nil
	} else {
		return errors.New("view can contain max 5 rows")
	}
}

func (v *View) marshal() ([]interface{}, error) {
	const timeout = 14.98 * 60
	if v.Timeout == 0 || v.Timeout > timeout {
		v.Timeout = timeout
	}
	var undo = map[string]bool{}
	var c []interface{}
	for _, row := range v.rows {
		hasButton := len(row.Buttons) > 0
		hasSelect := len(row.SelectMenu.Options) > 0
		tmp := map[string]interface{}{
			"type":       1,
			"components": []interface{}{},
		}
		for _, button := range row.Buttons {
			undo[button.customId] = true
			if v.OnTimeout != nil {
				timeoutTasks[button.customId] = []interface{}{v.Timeout, v.OnTimeout}
			}
			btn, err := button.marshal()
			if err != nil {
				return nil, err
			}
			tmp["components"] = append(tmp["components"].([]interface{}), btn)
		}
		if !hasButton {
			undo[row.SelectMenu.customId] = true
			if v.OnTimeout != nil {
				timeoutTasks[row.SelectMenu.customId] = []interface{}{v.Timeout, v.OnTimeout}
			}
			menu, err := row.SelectMenu.marshal()
			if err != nil {
				return nil, err
			}
			tmp["components"] = append(tmp["components"].([]interface{}), menu)
		}
		if len(undo) > 0 {
			c = append(c, tmp)
			go scheduleDeletion(v.Timeout, callbackTasks, undo)
		}
		if !(hasButton || hasSelect) {
			return nil, errors.New("view must contain at least one button or select menu")
		}
	}
	return c, nil
}
