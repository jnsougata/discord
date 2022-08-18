package discord

import (
	"fmt"
)

var groupBucket = map[string]interface{}{}
var subcommandBucket = map[string]interface{}{}
var autocompleteBucket = map[string]interface{}{}

type CommandType int

const (
	SlashCommand   CommandType = 1
	UserCommand    CommandType = 2
	MessageCommand CommandType = 3
)

type optionType int

type optionTypes struct {
	String      optionType
	Integer     optionType
	Boolean     optionType
	User        optionType
	Channel     optionType
	Role        optionType
	Mentionable optionType
	Number      optionType
	Attachment  optionType
}

var OptionTypes = optionTypes{
	String:      optionType(3),
	Integer:     optionType(4),
	Boolean:     optionType(5),
	User:        optionType(6),
	Channel:     optionType(7),
	Role:        optionType(8),
	Mentionable: optionType(9),
	Number:      optionType(10),
	Attachment:  optionType(11),
}

type channelType int

type channelTypes struct {
	GuildText          channelType
	DMText             channelType
	GuildVoice         channelType
	GroupDM            channelType
	GuildCategory      channelType
	GuildNews          channelType
	GuildNewsThread    channelType
	GuildPublicThread  channelType
	GuildPrivateThread channelType
	GuildStageVoice    channelType
	GuildDirectory     channelType
	GuildForum         channelType
}

var ChannelTypes = channelTypes{
	GuildText:          channelType(0),
	DMText:             channelType(1),
	GuildVoice:         channelType(2),
	GroupDM:            channelType(3),
	GuildCategory:      channelType(4),
	GuildNews:          channelType(5),
	GuildNewsThread:    channelType(10),
	GuildPublicThread:  channelType(11),
	GuildPrivateThread: channelType(12),
	GuildStageVoice:    channelType(13),
	GuildDirectory:     channelType(14),
	GuildForum:         channelType(15),
}

type Choice struct {
	Name  string      `json:"name"`
	Value interface{} `json:"Value"` // same type as type of option
}

type Option struct {
	Name         string     `json:"name"`
	Type         optionType `json:"type"`
	Value        any        `json:"Value"`   // available only during option parsing
	Focused      bool       `json:"focused"` // available only during option parsing
	Description  string
	Required     bool
	MinLength    int           // allowed for: StringOption
	MaxLength    int           // allowed for: StringOption
	MinValue     int64         // allowed for: IntegerOption, NumberOption
	MaxValue     int64         // allowed for: IntegerOption, NumberOption
	AutoComplete bool          // allowed for: StringOption, NumberOption, IntegerOption
	ChannelTypes []channelType // allowed for: ChannelOption
	Choices      []Choice      // allowed for: StringOption, IntegerOption, NumberOption
}

func (o *Option) marshal() map[string]interface{} {
	if o.Value != nil {
		panic("Option field {Value} must not be set while creating an option")
	}
	if o.Focused {
		panic("Option field {focused} must not be set while creating an option")
	}
	body := map[string]interface{}{}
	if o.Name == "" || o.Description == "" {
		panic("Both command {name} or {description} must be set")
	}
	if len(o.Name) > 32 {
		panic(fmt.Sprintf("Command (%s) {name} must be less than 32 characters", o.Name))
	}
	if len(o.Description) > 100 {
		panic(fmt.Sprintf("Command (%s) {description} must be less than 100 characters", o.Name))
	}
	body["name"] = o.Name
	body["description"] = o.Description
	body["type"] = o.Type
	body["required"] = o.Required
	switch o.Type {
	case OptionTypes.String:
		if o.MinLength > 0 && o.MaxLength > 0 && o.MinLength < o.MaxLength {
			body["min_length"] = o.MinLength
			body["max_length"] = o.MaxLength
		}
		if len(o.Choices) > 0 {
			body["choices"] = o.Choices
		} else if o.AutoComplete {
			body["auto_complete"] = true
		}
	case OptionTypes.Integer:
		body["min_value"] = o.MinValue
		body["max_value"] = o.MaxValue
		if len(o.Choices) > 0 {
			body["choices"] = o.Choices
		} else if o.AutoComplete {
			body["auto_complete"] = true
		}
	case OptionTypes.Number:
		body["min_value"] = o.MinValue
		body["max_value"] = o.MaxValue
		if len(o.Choices) > 0 {
			body["choices"] = o.Choices
		} else if o.AutoComplete {
			body["auto_complete"] = true
		}
	case OptionTypes.Channel:
		body["channel_types"] = []int{}
		for _, chanType := range o.ChannelTypes {
			body["channel_types"] = append(body["channel_types"].([]int), int(chanType))
		}
	}
	return body
}

type SubCommand struct {
	Name        string
	Description string
	Options     []Option
	Execute     func(bot Bot, ctx Context, options ResolvedOptions)
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
	body["options"] = []map[string]interface{}{}
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

// Command is a base type for all discord application commands
type Command struct {
	uniqueId          string
	Type              CommandType // defaults to chat input
	Name              string      // must be less than 32 characters
	Description       string      // must be less than 100 characters
	Options           []Option
	DMPermission      bool       // default: false
	MemberPermissions Permission // default: send_messages
	GuildId           int64
	subcommands       []SubCommand
	subcommandGroups  []SubcommandGroup
	Execute           func(bot Bot, ctx Context, options ResolvedOptions)
	AutocompleteTask  func(bot Bot, ctx Context, choices ...Choice)
}

func (cmd *Command) SubCommands(subcommands ...SubCommand) {
	for _, subcommand := range subcommands {
		cmd.subcommands = append(cmd.subcommands, subcommand)
	}
}

func (cmd *Command) SubcommandGroups(subcommandGroups ...SubcommandGroup) {
	for _, subcommandGroup := range subcommandGroups {
		cmd.subcommandGroups = append(cmd.subcommandGroups, subcommandGroup)
	}
}

func (cmd *Command) marshal() (
	map[string]interface{}, func(bot Bot, ctx Context, options ResolvedOptions), int64) {
	body := map[string]interface{}{}
	switch cmd.Type {
	case MessageCommand:
		body["type"] = int(MessageCommand)
		cmd.Type = MessageCommand
	case UserCommand:
		body["type"] = int(UserCommand)
		cmd.Type = UserCommand
	default:
		body["type"] = int(SlashCommand)
		cmd.Type = SlashCommand
	}
	cmd.uniqueId = assignId()
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
	switch cmd.Type {
	case SlashCommand:
		if cmd.Description == "" {
			panic("Command {description} must be set")
		}
		body["description"] = cmd.Description
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
				subcommandBucket[cmd.uniqueId] = map[string]interface{}{subcommand.Name: subcommand.Execute}
			}
			for _, subcommandGroup := range cmd.subcommandGroups {
				body["options"] = append(body["options"].([]map[string]interface{}), subcommandGroup.marshal())
				for _, subcommand := range subcommandGroup.subcommands {
					groupBucket[cmd.uniqueId] = map[string]interface{}{
						fmt.Sprintf(`%s_%s`, subcommandGroup.Name, subcommand.Name): subcommand.Execute,
					}
				}
			}
		}
		if cmd.AutocompleteTask != nil {
			autocompleteBucket[cmd.uniqueId] = cmd.AutocompleteTask
		}
	case UserCommand:
		if cmd.Description != "" {
			panic("Command {description} must not be set for user commands")
		}
		if len(cmd.Options) > 0 {
			panic("Command cannot have options for user commands")
		}
		if len(cmd.subcommands) > 0 {
			panic("Command cannot have subcommands for user commands")
		}
		if len(cmd.subcommandGroups) > 0 {
			panic("Command cannot have subcommand groups for user commands")
		}
	case MessageCommand:
		if cmd.Description != "" {
			panic("Command {description} must not be set for message commands")
		}
		if len(cmd.Options) > 0 {
			panic("Command cannot have options for message commands")
		}
		if len(cmd.subcommands) > 0 {
			panic("Command cannot have subcommands for user commands")
		}
		if len(cmd.subcommandGroups) > 0 {
			panic("Command cannot have subcommand groups for user commands")
		}
	}
	body["dm_permission"] = cmd.DMPermission
	switch int(cmd.MemberPermissions) {
	case 0:
		body["default_member_permissions"] = 1 << 11
	default:
		body["default_member_permissions"] = int(cmd.MemberPermissions)
	}
	return body, cmd.Execute, cmd.GuildId
}
