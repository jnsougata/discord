package discord

import (
	"errors"
	"fmt"
)

var groupBucket = map[string]interface{}{}
var subcommandBucket = map[string]interface{}{}
var autocompleteBucket = map[string]interface{}{}

type CommandKind int

const (
	SlashCommand   CommandKind = 1
	UserCommand    CommandKind = 2
	MessageCommand CommandKind = 3
)

type SubCommand struct {
	Name        string
	Description string
	options     []Option
	Execute     func(bot Bot, ctx Context, options ResolvedOptions)
}

func (sub *SubCommand) AddOption(option Option) error {
	if len(sub.options) == 25 {
		return errors.New("application command can only have max 25 options")
	}
	if option.Name == "" {
		return errors.New("option {name} must be set")
	}
	if len(option.Name) > 32 {
		return errors.New(fmt.Sprintf("option (%s) {name} must be less than 32 characters", option.Name))
	}
	if option.Description == "" {
		return errors.New(fmt.Sprintf("option (%s) {description} must be set", option.Name))
	}
	if len(option.Description) > 100 {
		return errors.New(fmt.Sprintf("option (%s) {description} must be less than 100 characters", option.Name))
	}
	if int(option.Type) == 0 {
		return errors.New(fmt.Sprintf("option (%s) {type} must be set", option.Name))
	}
	switch option.Type {
	case StringOption:
		if option.MaxLength < option.MinLength {
			return errors.New("option {MaxLength} must be greater than {MinLength}")
		}
		if option.MinLength < 0 || option.MaxLength < 0 {
			return errors.New("option {MinLength} and {MaxLength} must be greater than 0")
		}
		if option.MaxLength > 6000 || option.MinLength > 6000 {
			return errors.New("option length must be less than equals to 6000")
		}
		if option.MaxLength == 0 {
			option.MaxLength = 6000
		}
		if option.MinLength == 0 {
			option.MinLength = 1
		}
	case NumberOption:
		if option.MaxValue < option.MinValue {
			return errors.New("option {MaxValue} must be greater than {MinValue}")
		}
	case IntegerOption:
		if option.MaxValue < option.MinValue {
			return errors.New("option {MaxValue} must be greater than {MinValue}")
		}
	}
	if !(option.Type == StringOption || option.Type == IntegerOption || option.Type == NumberOption) {
		if len(option.Choices) > 0 {
			return errors.New(
				"option {choices} can only be used with {string} or {integer} or {number} type options")
		}
		if option.AutoComplete {
			return errors.New(
				"option {autocomplete} can only be used with {string} or {integer} or {number} type options")
		}
	}
	sub.options = append(sub.options, option)
	return nil
}

func (sub *SubCommand) marshal() map[string]interface{} {
	body := map[string]interface{}{}
	if sub.Name == "" || sub.Description == "" {
		panic("Both command {name} or {description} must be set")
	}
	if len(sub.Name) > 32 {
		panic(fmt.Sprintf("Command (%s) {name} must be less than 32 characters", sub.Name))
	}
	if len(sub.Description) > 100 {
		panic(fmt.Sprintf("Command (%s) {description} must be less than 100 characters", sub.Name))
	}
	body["type"] = 1
	body["name"] = sub.Name
	body["description"] = sub.Description
	body["options"] = []map[string]interface{}{}
	for _, option := range sub.options {
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
	Type              CommandKind // defaults to chat input
	Name              string      // must be less than 32 characters
	Description       string      // must be less than 100 characters
	options           []Option
	DMPermission      bool       // default: false
	MemberPermissions Permission // default: send_messages
	GuildId           int64
	subcommands       []SubCommand
	subcommandGroups  []SubcommandGroup
	Execute           func(bot Bot, ctx Context, options ResolvedOptions)
	AutocompleteTask  func(bot Bot, ctx Context, choices ...Choice)
}

func (cmd *Command) AddOption(option Option) error {
	if len(cmd.options) == 25 {
		return errors.New("application command can only have max 25 options")
	}
	if option.Name == "" {
		return errors.New("option {name} must be set")
	}
	if len(option.Name) > 32 {
		return errors.New(fmt.Sprintf("option (%s) {name} must be less than 32 characters", option.Name))
	}
	if option.Description == "" {
		return errors.New(fmt.Sprintf("option (%s) {description} must be set", option.Name))
	}
	if len(option.Description) > 100 {
		return errors.New(fmt.Sprintf("option (%s) {description} must be less than 100 characters", option.Name))
	}
	if int(option.Type) == 0 {
		return errors.New(fmt.Sprintf("option (%s) {type} must be set", option.Name))
	}
	switch option.Type {
	case StringOption:
		if option.MaxLength < option.MinLength {
			return errors.New("option {MaxLength} must be greater than {MinLength}")
		}
		if option.MinLength < 0 || option.MaxLength < 0 {
			return errors.New("option {MinLength} and {MaxLength} must be greater than 0")
		}
		if option.MaxLength > 6000 || option.MinLength > 6000 {
			return errors.New("option length must be less than equals to 6000")
		}
		if option.MaxLength == 0 {
			option.MaxLength = 6000
		}
		if option.MinLength == 0 {
			option.MinLength = 1
		}
	case NumberOption:
		if option.MaxValue < option.MinValue {
			return errors.New("option {MaxValue} must be greater than {MinValue}")
		}
	case IntegerOption:
		if option.MaxValue < option.MinValue {
			return errors.New("option {MaxValue} must be greater than {MinValue}")
		}
	}
	if !(option.Type == StringOption || option.Type == IntegerOption || option.Type == NumberOption) {
		if len(option.Choices) > 0 {
			return errors.New(
				"option {choices} can only be used with {string} or {integer} or {number} type options")
		}
		if option.AutoComplete {
			return errors.New(
				"option {autocomplete} can only be used with {string} or {integer} or {number} type options")
		}
	}
	cmd.options = append(cmd.options, option)
	return nil
}

func (cmd *Command) SubCommand(subcommand SubCommand) {
	cmd.subcommands = append(cmd.subcommands, subcommand)
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
	case UserCommand:
		body["type"] = int(UserCommand)
	case SlashCommand:
		body["type"] = int(SlashCommand)
	default:
		body["type"] = int(SlashCommand)
		cmd.Type = SlashCommand
	}
	cmd.uniqueId = assignId()
	body["name"] = cmd.Name
	switch cmd.Type {
	case SlashCommand:
		if cmd.Description == "" {
			panic("Command {description} must be set")
		}
		body["description"] = cmd.Description
		body["options"] = []map[string]interface{}{}
		if len(cmd.options) > 0 && len(cmd.subcommands) > 0 {
			panic("Command cannot have both options and Subcommands")
		} else if len(cmd.options) > 0 {
			for _, option := range cmd.options {
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
		if len(cmd.options) > 0 {
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
		if len(cmd.options) > 0 {
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
