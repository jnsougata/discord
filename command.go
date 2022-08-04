package disgo

import (
	"fmt"
)

type CommandOption struct {
	Name         string          `json:"name"`
	Description  string          `json:"description"`
	Type         int             `json:"type"` // 3: string, 4: integer, 5: boolean, 6: user, 7: channel, 8: role, 9: mentionable, 10: number, 11: attachment
	Required     bool            `json:"required,omitempty"`
	MinLength    int             `json:"min_length,omitempty"`    // for type 3 only
	MaxLength    int             `json:"max_length,omitempty"`    // for type 3 only
	MinValue     int64           `json:"min_value,omitempty"`     // for type 4 and 10 only
	MaxValue     int64           `json:"max_value,omitempty"`     // for type 4 and 10 only
	AutoComplete bool            `json:"auto_complete,omitempty"` // for type 1 and
	ChannelTypes []int           `json:"channel_types,omitempty"` // 0: guild text channel, 1: DM channel, 2: guild voice channel, 3: group DM channel, 4: guild category, 5: guild news, 10: guild news thread, 11: guild public thread, 12: guild private thread, 13: guild stage voice, 14: guild directory, 15: guild forum
	Options      []CommandOption `json:"options,omitempty"`       // for type 1 and 2 only
	Choices      []Choice        `json:"choices,omitempty"`       // for type 3 and 4 and 10 only
}

func (co *CommandOption) Marshal() map[string]interface{} {
	body := map[string]interface{}{}
	if co.Name == "" || co.Description == "" {
		panic("Both command {name} or {description} must be set")
	}
	if len(co.Name) > 32 {
		panic(fmt.Sprintf("Command (%s) {name} must be less than 32 characters", co.Name))
	}
	if len(co.Description) > 100 {
		panic(fmt.Sprintf("Command (%s) {description} must be less than 100 characters", co.Name))
	}
	body["name"] = co.Name
	body["description"] = co.Description
	body["type"] = co.Type
	switch co.Type {
	case 3:
		body["min_length"] = co.MinLength
		body["max_length"] = co.MaxLength
		body["choices"] = co.Choices
	case 4:
		body["min_value"] = co.MinValue
		body["max_value"] = co.MaxValue
	case 10:
		body["min_value"] = co.MinValue
		body["max_value"] = co.MaxValue

	}

	return body
}

// ApplicationCommand is a base type for all discord application commands
type ApplicationCommand struct {
	Type              int    // 1: slash command, 2: user command, 3: message command
	Name              string // must be less than 32 characters
	Description       string // must be less than 100 characters
	Options           []CommandOption
	DMPermission      bool // default: false
	MemberPermissions int  // default: send_messages
	GuildId           int64
	handler           func(bot BotUser, ctx Context, options ...SlashCommandOption)
	autocomplete      func(bot BotUser, ctx Context, choices ...Choice)
}

func (cmd *ApplicationCommand) Handler(handler func(bot BotUser, ctx Context, options ...SlashCommandOption)) {
	cmd.handler = handler
}

func (cmd *ApplicationCommand) AutoCompleteHandler(handler func(bot BotUser, ctx Context, choices ...Choice)) {
	cmd.autocomplete = handler
}

func (cmd *ApplicationCommand) marshal() (
	map[string]interface{},
	func(bot BotUser, ctx Context, options ...SlashCommandOption),
	int64) {
	body := map[string]interface{}{}
	switch cmd.Type {
	case 3:
		body["type"] = 3
	case 2:
		body["type"] = 2
	default:
		body["type"] = 1
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
	switch cmd.MemberPermissions {
	case 0:
		body["default_member_permissions"] = 1 << 11
	default:
		body["default_member_permissions"] = cmd.MemberPermissions
	}
	if cmd.Type == 1 {
		body["options"] = cmd.Options
	}
	return body, cmd.handler, cmd.GuildId
}

type Choice struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"` // same type as type of option
}
