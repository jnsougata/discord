package discord

import (
	"fmt"
	"strconv"
)

var groupBucket = map[string]interface{}{}
var subcommandBucket = map[string]interface{}{}
var autocompleteBucket = map[string]interface{}{}

type SubCommand struct {
	Name        string
	Description string
	Options     []Option
	Execute     func(bot Bot, ctx Context, options ResolvedOptions)
}

func (scmd *SubCommand) marshal() map[string]interface{} {
	body := map[string]interface{}{
		"type":        1,
		"name":        scmd.Name,
		"description": scmd.Description,
		"options":     []map[string]interface{}{},
	}
	for _, o := range scmd.Options {
		body["options"] = append(body["options"].([]map[string]interface{}), o.marshal())
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
	body := map[string]interface{}{
		"type":        2,
		"name":        scg.Name,
		"description": scg.Description,
		"options":     []map[string]interface{}{},
	}
	for _, subcommand := range scg.subcommands {
		body["options"] = append(body["options"].([]map[string]interface{}), subcommand.marshal())
	}
	return body
}

// Command is a base type for all discord application commands
type Command struct {
	uniqueId         string
	Type             int    // defaults to chat input
	Name             string // must be less than 32 characters
	Description      string // must be less than 100 characters
	Options          []Option
	DMPermission     bool         // default: false
	Permissions      []Permission // default: send_messages
	GuildId          int64
	subcommands      []SubCommand
	subcommandGroups []SubcommandGroup
	Execute          func(bot Bot, ctx Context, options ResolvedOptions)
	AutocompleteTask func(bot Bot, ctx Context, choices ...Choice)
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
	cmd.uniqueId = assignId()
	body := map[string]interface{}{
		"name":          cmd.Name,
		"description":   cmd.Description,
		"options":       []map[string]interface{}{},
		"dm_permission": cmd.DMPermission,
	}
	switch cmd.Type {
	case CommandTypes.Message:
		body["type"] = CommandTypes.Message
	case CommandTypes.User:
		body["type"] = CommandTypes.User
	case CommandTypes.Slash:
		body["type"] = CommandTypes.Slash
	default:
		body["type"] = CommandTypes.Slash
		cmd.Type = CommandTypes.Slash
	}
	for _, option := range cmd.Options {
		body["options"] = append(body["options"].([]map[string]interface{}), option.marshal())
	}
	for _, subcommand := range cmd.subcommands {
		body["Options"] = append(body["options"].([]map[string]interface{}), subcommand.marshal())
		subcommandBucket[cmd.uniqueId] = map[string]interface{}{subcommand.Name: subcommand.Execute}
	}
	for _, subcommandGroup := range cmd.subcommandGroups {
		body["Options"] = append(body["Options"].([]map[string]interface{}), subcommandGroup.marshal())
		for _, subcommand := range subcommandGroup.subcommands {
			groupBucket[cmd.uniqueId] = map[string]interface{}{
				fmt.Sprintf(`%s_%s`, subcommandGroup.Name, subcommand.Name): subcommand.Execute,
			}
		}
	}
	if cmd.AutocompleteTask != nil {
		autocompleteBucket[cmd.uniqueId] = cmd.AutocompleteTask
	}
	if len(cmd.Permissions) > 0 {
		p := 0
		for _, permission := range cmd.Permissions {
			p |= int(permission)
			body["default_member_permissions"] = strconv.Itoa(p)
		}
	} else {
		body["default_member_permissions"] = strconv.Itoa(1 << 11)
	}
	return body, cmd.Execute, cmd.GuildId
}
