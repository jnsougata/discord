package component

import (
	"fmt"
	"github.com/jnsougata/disgo/core/emoji"
	"log"
)

const (
	ActionRowType        = 1
	ButtonType           = 2
	SelectMenuType       = 3
	PrimaryButtonStyle   = 1
	SecondaryButtonStyle = 2
	SuccessButtonStyle   = 3
	DangerButtonStyle    = 4
	LinkButtonStyle      = 5
)

type SelectOption struct {
	Label       string        `json:"label"`
	Value       string        `json:"value"`
	Description string        `json:"description"`
	Emoji       emoji.Partial `json:"emoji"`
	Default     bool          `json:"default"`
}

type Button struct {
	Type     int           `json:"type"`
	Style    int           `json:"style"`
	Label    string        `json:"label"`
	Emoji    emoji.Partial `json:"emoji,omitempty"`
	CustomId string        `json:"custom_id"`
	URL      string        `json:"url,omitempty"`
	Disabled bool          `json:"disabled"`
}

type SelectMenu struct {
	Type        int            `json:"type"`
	CustomId    string         `json:"custom_id"`
	Options     []SelectOption `json:"options"`
	Placeholder string         `json:"placeholder"`
	MinValues   int            `json:"min_values"`
	MaxValues   int            `json:"max_values"`
	Disabled    bool           `json:"disabled"`
}

type ActionRow struct {
	Buttons     []Button
	SelectMenus []SelectMenu
}

type View struct {
	ActionRows []ActionRow
}

func (v *View) ToComponent() []interface{} {
	var c []interface{}
	ids := map[string]bool{}
	if len(v.ActionRows) > 0 && len(v.ActionRows) <= 5 {
		for _, row := range v.ActionRows {
			numButtons := 0
			tmp := map[string]interface{}{
				"type":       1,
				"components": []interface{}{},
			}
			for _, button := range row.Buttons {
				numButtons++
				if button.CustomId == "" && button.Style != LinkButtonStyle {
					log.Println(
						fmt.Sprintf("CustomId must be provided with non-link button `%s`", button.Label))
				} else if _, ok := ids[button.CustomId]; !ok {
					if numButtons <= 5 {
						ids[button.CustomId] = true
						tmp["components"] = append(tmp["components"].([]interface{}), button)
					} else {
						log.Println("An Action Row can either contain max 5x Buttons")
					}
				} else {
					log.Println(
						fmt.Sprintf("CustomId `%s` already used with a previous component", button.CustomId))
				}
			}
			for _, selectMenu := range row.SelectMenus {
				if selectMenu.CustomId == "" {
					log.Println("CustomId must be provided with select menu")
				} else if selectMenu.MaxValues > 25 {
					log.Println("MaxValues must be less than or equals to 25")
				} else if selectMenu.MinValues > selectMenu.MaxValues {
					log.Println("MinValues must be less than or equals to MaxValues")
				} else if selectMenu.MinValues < 0 {
					log.Println("MinValues must be greater than or equals to 0")
				} else if _, ok := ids[selectMenu.CustomId]; !ok {
					if numButtons == 0 {
						ids[selectMenu.CustomId] = true
						tmp["components"] = append(tmp["components"].([]interface{}), selectMenu)
					} else {
						log.Println("An Action Row can contain one of these: (1x SelectMenu) or (max 5x Buttons)")
					}
				} else {
					log.Println(
						fmt.Sprintf("CustomId `%s` already used with a previous component", selectMenu.CustomId))
				}
			}
			if len(tmp["components"].([]interface{})) > 0 {
				c = append(c, tmp)
			}

		}
	} else {
		log.Println("ActionRows must contain between 1 and 5 rows")
	}
	return c
}
