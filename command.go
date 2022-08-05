package disgo

import (
	"fmt"
)

var subcommandBucket = map[string]interface{}{}
var groupBucket = map[string]interface{}{}

type CommandType int

const (
	SlashCommand   CommandType = 1
	UserCommand    CommandType = 2
	MessageCommand CommandType = 3
)

type OptionType int

const (
	SubCommandType      OptionType = 1
	SupCommandGroupType OptionType = 2
	StringOption        OptionType = 3
	IntegerOption       OptionType = 4
	BooleanOption       OptionType = 5
	UserOption          OptionType = 6
	ChannelOption       OptionType = 7
	RoleOption          OptionType = 8
	MentionableOption   OptionType = 9
	NumberOption        OptionType = 10
	AttachmentOption    OptionType = 11
)

type ChannelType int

const (
	GuildText          ChannelType = 0
	DMText             ChannelType = 1
	GuildVoice         ChannelType = 2
	GroupDM            ChannelType = 3
	GuildCategory      ChannelType = 4
	GuildNews          ChannelType = 5
	GuildNewsThread    ChannelType = 10
	GuildPublicThread  ChannelType = 11
	GuildPrivateThread ChannelType = 12
	GuildStageVoice    ChannelType = 13
	GuildDirectory     ChannelType = 14
	GuildForum         ChannelType = 15
)

type CommandOption struct {
	Name         string
	Description  string
	Type         OptionType
	Required     bool
	MinLength    int           // allowed for: StringOption
	MaxLength    int           // allowed for: StringOption
	MinValue     int64         // allowed for: IntegerOption, NumberOption
	MaxValue     int64         // allowed for: IntegerOption, NumberOption
	AutoComplete bool          // allowed for: StringOption, NumberOption, IntegerOption
	ChannelTypes []ChannelType // allowed for: ChannelOption
	Choices      []Choice      // allowed for: StringOption, IntegerOption, NumberOption
}

func (co *CommandOption) marshal() map[string]interface{} {
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
	body["required"] = co.Required
	switch int(co.Type) {
	case 3:
		body["min_length"] = co.MinLength
		body["max_length"] = co.MaxLength
		if len(co.Choices) > 0 {
			body["choices"] = co.Choices
		} else if co.AutoComplete {
			body["auto_complete"] = true
		}
	case 4:
		body["min_value"] = co.MinValue
		body["max_value"] = co.MaxValue
		if len(co.Choices) > 0 {
			body["choices"] = co.Choices
		} else if co.AutoComplete {
			body["auto_complete"] = true
		}
	case 10:
		body["min_value"] = co.MinValue
		body["max_value"] = co.MaxValue
		if len(co.Choices) > 0 {
			body["choices"] = co.Choices
		} else if co.AutoComplete {
			body["auto_complete"] = true
		}
	case 7:
		for _, channelType := range co.ChannelTypes {
			body["channel_types"] = append(body["channel_types"].([]int), int(channelType))
		}
	}
	//body["options"] = co.Options
	return body
}

type SubCommand struct {
	Name        string
	Description string
	Options     []CommandOption
	handler     func(bot BotUser, ctx Context, options ...SlashCommandOption)
}

func (sc *SubCommand) Handler(handler func(bot BotUser, ctx Context, options ...SlashCommandOption)) {
	sc.handler = handler
}

func (sc *SubCommand) marshal() map[string]interface{} {
	body := map[string]interface{}{}
	if sc.Name == "" || sc.Description == "" {
		panic("Both command {name} or {description} must be set")
	}
	if len(sc.Name) > 32 {
		panic(fmt.Sprintf("Command (%s) {name} must be less than 32 characters", sc.Name))
	}
	if len(sc.Description) > 100 {
		panic(fmt.Sprintf("Command (%s) {description} must be less than 100 characters", sc.Name))
	}
	body["type"] = 1
	body["name"] = sc.Name
	body["description"] = sc.Description
	for _, option := range sc.Options {
		body["options"] = append(body["options"].([]map[string]interface{}), option.marshal())
	}
	return body
}

type SubcommandGroup struct {
	Name        string
	Description string
	subcommands []SubCommand
}

func (scg *SubcommandGroup) Subcommands(subcommands ...SubCommand) {
	for _, subcommand := range subcommands {
		scg.subcommands = append(scg.subcommands, subcommand)
	}
}

func (scg *SubcommandGroup) marshal() map[string]interface{} {
	body := map[string]interface{}{}
	if scg.Name == "" || scg.Description == "" {
		panic("Both command {name} or {description} must be set")
	}
	if len(scg.Name) > 32 {
		panic(fmt.Sprintf("Command (%s) {name} must be less than 32 characters", scg.Name))
	}
	if len(scg.Description) > 100 {
		panic(fmt.Sprintf("Command (%s) {description} must be less than 100 characters", scg.Name))
	}
	body["type"] = 2
	body["name"] = scg.Name
	body["description"] = scg.Description
	for _, subcommand := range scg.subcommands {
		body["options"] = append(body["options"].([]map[string]interface{}), subcommand.marshal())
	}
	return body
}

// ApplicationCommand is a base type for all discord application commands
type ApplicationCommand struct {
	Type              CommandType
	Name              string // must be less than 32 characters
	Description       string // must be less than 100 characters
	Options           []CommandOption
	DMPermission      bool // default: false
	MemberPermissions int  // default: send_messages
	GuildId           int64
	uniqueId          string
	handler           func(bot BotUser, ctx Context, options ...SlashCommandOption)
	autocomplete      func(bot BotUser, ctx Context, choices ...Choice)
	subcommands       []SubCommand
	subcommandGroups  []SubcommandGroup
}

func (cmd *ApplicationCommand) Handler(handler func(bot BotUser, ctx Context, options ...SlashCommandOption)) {
	cmd.handler = handler
}

func (cmd *ApplicationCommand) AutoCompleteHandler(handler func(bot BotUser, ctx Context, choices ...Choice)) {
	cmd.autocomplete = handler
}

func (cmd *ApplicationCommand) SubCommands(subcommands ...SubCommand) {
	for _, subcommand := range subcommands {
		cmd.subcommands = append(cmd.subcommands, subcommand)
	}
}

func (cmd *ApplicationCommand) SubcommandGroups(subcommandGroups ...SubcommandGroup) {
	for _, subcommandGroup := range subcommandGroups {
		cmd.subcommandGroups = append(cmd.subcommandGroups, subcommandGroup)
	}
}

func (cmd *ApplicationCommand) marshal() (
	map[string]interface{},
	func(bot BotUser, ctx Context, options ...SlashCommandOption),
	int64) {
	body := map[string]interface{}{}
	switch int(cmd.Type) {
	case 3:
		body["type"] = 3
		cmd.Type = CommandType(3)
	case 2:
		body["type"] = 2
		cmd.Type = CommandType(2)
	default:
		body["type"] = 1
		cmd.Type = CommandType(1)
	}
	cmd.uniqueId = assignId("")
	if cmd.Name == "" {
		panic("Command {name} must be set")
	}
	if len(cmd.Name) > 32 {
		panic(fmt.Sprintf("Command (%s) {name} must be less than 32 characters", cmd.Name))
	}
	if len(cmd.Description) > 100 {
		panic(fmt.Sprintf("Command (%s) {description} must be less than 100 characters", cmd.Name))
	}
	body["name"] = cmd.Name
	if cmd.Type == SlashCommand && cmd.Description == "" {
		panic("Command {description} must be set for command " + cmd.Name)
	} else if cmd.Type != SlashCommand && cmd.Description != "" {
		panic("Command {description} is only allowed for SlashCommand")
	} else {
		body["description"] = cmd.Description
	}
	body["dm_permission"] = cmd.DMPermission
	switch cmd.MemberPermissions {
	case 0:
		body["default_member_permissions"] = 1 << 11
	default:
		body["default_member_permissions"] = cmd.MemberPermissions
	}
	if int(cmd.Type) == 1 {
		body["options"] = []map[string]interface{}{}
		if len(cmd.Options) > 0 && len(cmd.subcommands) > 0 {
			panic("Command cannot have both options and Subcommands")
		} else if len(cmd.Options) > 0 {
			for _, option := range cmd.Options {
				body["options"] = append(body["options"].([]map[string]interface{}), option.marshal())
			}
		} else if len(cmd.subcommands) > 0 || len(cmd.subcommandGroups) > 0 {
			for _, subcommand := range cmd.subcommands {
				body["options"] = append(body["options"].([]map[string]interface{}), subcommand.marshal())
				subcommandBucket[cmd.uniqueId] = map[string]interface{}{subcommand.Name: subcommand.handler}
			}
			for _, subcommandGroup := range cmd.subcommandGroups {
				body["options"] = append(body["options"].([]map[string]interface{}), subcommandGroup.marshal())
				for _, subcommand := range subcommandGroup.subcommands {
					groupBucket[cmd.uniqueId] = map[string]interface{}{
						fmt.Sprintf(`%s_%s`, subcommandGroup.Name, subcommand.Name): subcommand.handler,
					}
				}
			}
		}
	}
	return body, cmd.handler, cmd.GuildId
}

type Choice struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"` // same type as type of option
}
