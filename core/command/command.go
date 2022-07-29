package command

import (
	"fmt"
	"github.com/jnsougata/disgo/core/interaction"
	"github.com/jnsougata/disgo/core/user"
)

type ApplicationCommand struct {
	Type              int    // 0: slash command, 1: user command, 2: message command
	Name              string // must be less than 32 characters
	Description       string // must be less than 100 characters
	Options           []Option
	DMPermission      bool // default: false
	MemberPermissions int  // default: send_messages
	GuildId           int64
	Handler           func(bot user.Bot, ctx Context, options ...interaction.Option)
}

func (cmd *ApplicationCommand) ToData() (
	map[string]interface{},
	func(bot user.Bot, ctx Context, options ...interaction.Option),
	int64) {
	body := map[string]interface{}{}
	switch cmd.Type {
	case 0:
		body["type"] = 1
	case 1:
		body["type"] = 1
	case 2:
		body["type"] = 2
	case 3:
		body["type"] = 3
	}
	if cmd.Name == "" || cmd.Description == "" {
		panic("Both command {name} or {description} must be set")
	}
	if len(cmd.Name) > 32 {
		panic(fmt.Sprintf("Command (%s) {name} must be less than 32 characters", cmd.Name))
	}
	if len(cmd.Description) > 100 {
		panic(fmt.Sprintf("Command (%s) {description} must be less than 100 characters", cmd.Name))
	}
	body["name"] = cmd.Name
	body["description"] = cmd.Description
	body["dm_permission"] = cmd.DMPermission
	if cmd.MemberPermissions != 0 {
		body["default_member_permissions"] = cmd.MemberPermissions
	} else {
		body["default_member_permissions"] = 1 << 11
	}
	if cmd.Type == 1 {
		body["options"] = cmd.Options
	}
	return body, cmd.Handler, cmd.GuildId
}

type Option struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Type         int      `json:"type"` // 3: string, 4: integer, 5: boolean, 6: user, 7: channel, 8: role, 9: mentionable, 10: number, 11: attachment
	Required     bool     `json:"required,omitempty"`
	MinLength    int      `json:"min_length,omitempty"`    // for type 3 only
	MaxLength    int      `json:"max_length,omitempty"`    // for type 3 only
	MinValue     int64    `json:"min_value,omitempty"`     // for type 4 and 10 only
	MaxValue     int64    `json:"max_value,omitempty"`     // for type 4 and 10 only
	AutoComplete bool     `json:"auto_complete,omitempty"` // for type 1 and
	ChannelTypes []int    `json:"channel_types,omitempty"` // 0: guild text channel, 1: DM channel, 2: guild voice channel, 3: group DM channel, 4: guild category, 5: guild news, 10: guild news thread, 11: guild public thread, 12: guild private thread, 13: guild stage voice, 14: guild directory, 15: guild forum
	Options      []Option `json:"options,omitempty"`       // for type 1 and 2 only
	Choices      []Choice `json:"choices,omitempty"`       // for type 3 or 4 or 10 only
}

type Choice struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"` // same type as type of option
}
