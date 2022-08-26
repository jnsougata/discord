package discord

import (
	"errors"
	"fmt"
	"strconv"
)

var groupBucket = map[string]interface{}{}
var subcommandBucket = map[string]interface{}{}
var autocompleteBucket = map[string]interface{}{}

type SubCommand struct {
	Name        string
	Description string
	options     []option
	Execute     func(bot Bot, ctx Context, options ResolvedOptions)
}

func (scmd *SubCommand) OptionSTRING(
	name string,
	description string,
	required bool,
	minLength int,
	maxLength int,
	autocomplete bool,
	choices ...Choice,
) error {
	if len(scmd.options) == 25 {
		return errors.New("application command can only have max 25 options")
	}
	err := checkNameDesc(name, description)
	if err != nil {
		return err
	}
	if minLength < 0 || maxLength < 0 {
		return errors.New("option {MinLength} and {MaxLength} must be greater than 0")
	}
	if minLength > 6000 || maxLength > 6000 {
		return errors.New("option length must be less than equals to 6000")
	}
	if minLength == 0 {
		minLength = 1
	}
	if maxLength == 0 {
		maxLength = 6000
	}
	if maxLength < minLength {
		return errors.New("option {maxLength} must be greater than {minLength}")
	}
	if len(choices) > 0 && autocomplete {
		return errors.New("option {choices} can only be used with {autocomplete} disabled")
	}
	scmd.options = append(scmd.options, option{
		Type:         stringOption,
		Name:         name,
		Description:  description,
		MinLength:    minLength,
		MaxLength:    maxLength,
		AutoComplete: autocomplete,
		Choices:      choices,
		Required:     required,
	})
	return nil
}

func (scmd *SubCommand) OptionINTEGER(
	name string,
	description string,
	required bool,
	minValue int,
	maxValue int,
	autocomplete bool,
	choices ...Choice,
) error {
	if len(scmd.options) == 25 {
		return errors.New("application command can only have max 25 options")
	}
	err := checkNameDesc(name, description)
	if err != nil {
		return err
	}
	if minValue > maxValue {
		return errors.New("option {minValue} must be less than {maxValue}")
	}
	if len(choices) > 0 && autocomplete {
		return errors.New("option {choices} can only be used with {autocomplete} disabled")
	}
	scmd.options = append(scmd.options, option{
		Type:         integerOption,
		Name:         name,
		Description:  description,
		MinValue:     minValue,
		MaxValue:     maxValue,
		AutoComplete: autocomplete,
		Choices:      choices,
		Required:     required,
	})
	return nil
}

func (scmd *SubCommand) OptionBOOLEAN(name string, description string, required bool) error {
	if len(scmd.options) == 25 {
		return errors.New("application command can only have max 25 options")
	}
	err := checkNameDesc(name, description)
	if err != nil {
		return err
	}
	scmd.options = append(scmd.options, option{
		Type:        booleanOption,
		Name:        name,
		Description: description,
		Required:    required,
	})
	return nil
}

func (scmd *SubCommand) OptionUSER(name string, description string, required bool) error {
	if len(scmd.options) == 25 {
		return errors.New("application command can only have max 25 options")
	}
	err := checkNameDesc(name, description)
	if err != nil {
		return err
	}
	scmd.options = append(scmd.options, option{
		Type:        userOption,
		Name:        name,
		Description: description,
		Required:    required,
	})
	return nil
}

func (scmd *SubCommand) OptionCHANNEL(name string, description string, required bool, channelTypes ...ChannelKind) error {
	if len(scmd.options) == 25 {
		return errors.New("application command can only have max 25 options")
	}
	err := checkNameDesc(name, description)
	if err != nil {
		return err
	}
	kinds := map[int]ChannelKind{
		0:  ChannelKinds.Text,
		1:  ChannelKinds.DM,
		2:  ChannelKinds.GroupDM,
		3:  ChannelKinds.Voice,
		4:  ChannelKinds.Category,
		5:  ChannelKinds.News,
		10: ChannelKinds.NewsThread,
		11: ChannelKinds.PublicThread,
		12: ChannelKinds.PrivateThread,
		13: ChannelKinds.StageVoice,
		14: ChannelKinds.Directory,
		15: ChannelKinds.Forum,
	}
	for _, ch := range channelTypes {
		if _, ok := kinds[int(ch)]; !ok {
			return errors.New(fmt.Sprintf("option {ChannelTypes} contains invalid channel kind (%d)", ch))
		}
	}
	scmd.options = append(scmd.options, option{
		Type:         channelOption,
		Name:         name,
		Description:  description,
		ChannelTypes: channelTypes,
		Required:     required,
	})
	return nil
}

func (scmd *SubCommand) OptionROLE(name string, description string, required bool) error {
	if len(scmd.options) == 25 {
		return errors.New("application command can only have max 25 options")
	}
	err := checkNameDesc(name, description)
	if err != nil {
		return err
	}
	scmd.options = append(scmd.options, option{
		Type:        roleOption,
		Name:        name,
		Description: description,
		Required:    required,
	})
	return nil
}

func (scmd *SubCommand) OptionMENTIONABLE(name string, description string, required bool) error {
	if len(scmd.options) == 25 {
		return errors.New("application command can only have max 25 options")
	}
	err := checkNameDesc(name, description)
	if err != nil {
		return err
	}
	scmd.options = append(scmd.options, option{
		Type:        mentionableOption,
		Name:        name,
		Description: description,
		Required:    required,
	})
	return nil
}

func (scmd *SubCommand) OptionNUMBER(
	name string,
	description string,
	required bool,
	minValue float64,
	maxValue float64,
	autocomplete bool,
	choices ...Choice,
) error {
	if len(scmd.options) == 25 {
		return errors.New("application command can only have max 25 options")
	}
	err := checkNameDesc(name, description)
	if err != nil {
		return err
	}
	if minValue > maxValue {
		return errors.New("option {minValue} must be less than {maxValue}")
	}
	if len(choices) > 0 && autocomplete {
		return errors.New("option {choices} can only be used with {autocomplete} disabled")
	}
	scmd.options = append(scmd.options, option{
		Type:         numberOption,
		Name:         name,
		Description:  description,
		MinValueNum:  minValue,
		MaxValueNum:  maxValue,
		AutoComplete: autocomplete,
		Choices:      choices,
		Required:     required,
	})
	return nil
}

func (scmd *SubCommand) OptionATTACHMENT(name string, description string, required bool) error {
	if len(scmd.options) == 25 {
		return errors.New("application command can only have max 25 options")
	}
	err := checkNameDesc(name, description)
	if err != nil {
		return err
	}
	scmd.options = append(scmd.options, option{
		Type:        attachmentOption,
		Name:        name,
		Description: description,
		Required:    required,
	})
	return nil
}

func (scmd *SubCommand) marshal() map[string]interface{} {
	body := map[string]interface{}{}
	if scmd.Name == "" || scmd.Description == "" {
		panic("Both command {name} or {description} must be set")
	}
	if len(scmd.Name) > 32 {
		panic(fmt.Sprintf("Command (%s) {name} must be less than 32 characters", scmd.Name))
	}
	if len(scmd.Description) > 100 {
		panic(fmt.Sprintf("Command (%s) {description} must be less than 100 characters", scmd.Name))
	}
	body["type"] = 1
	body["name"] = scmd.Name
	body["description"] = scmd.Description
	body["options"] = []map[string]interface{}{}
	for _, o := range scmd.options {
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
	uniqueId         string
	Type             commandKind // defaults to chat input
	Name             string      // must be less than 32 characters
	Description      string      // must be less than 100 characters
	options          []option
	DMPermission     bool         // default: false
	Permissions      []Permission // default: send_messages
	GuildId          int64
	subcommands      []SubCommand
	subcommandGroups []SubcommandGroup
	Execute          func(bot Bot, ctx Context, options ResolvedOptions)
	AutocompleteTask func(bot Bot, ctx Context, choices ...Choice)
}

func (cmd *Command) OptionSTRING(
	name string,
	description string,
	required bool,
	minLength int,
	maxLength int,
	autocomplete bool,
	choices ...Choice,
) error {
	if len(cmd.options) == 25 {
		return errors.New("application command can only have max 25 options")
	}
	err := checkNameDesc(name, description)
	if err != nil {
		return err
	}
	if minLength < 0 || maxLength < 0 {
		return errors.New("option {MinLength} and {MaxLength} must be greater than 0")
	}
	if minLength > 6000 || maxLength > 6000 {
		return errors.New("option length must be less than equals to 6000")
	}
	if minLength == 0 {
		minLength = 1
	}
	if maxLength == 0 {
		maxLength = 6000
	}
	if maxLength < minLength {
		return errors.New("option {maxLength} must be greater than {minLength}")
	}
	if len(choices) > 0 && autocomplete {
		return errors.New("option {choices} can only be used with {autocomplete} disabled")
	}
	cmd.options = append(cmd.options, option{
		Type:         stringOption,
		Name:         name,
		Description:  description,
		MinLength:    minLength,
		MaxLength:    maxLength,
		AutoComplete: autocomplete,
		Choices:      choices,
		Required:     required,
	})
	return nil
}

func (cmd *Command) OptionINTEGER(
	name string,
	description string,
	required bool,
	minValue int,
	maxValue int,
	autocomplete bool,
	choices ...Choice,
) error {
	if len(cmd.options) == 25 {
		return errors.New("application command can only have max 25 options")
	}
	err := checkNameDesc(name, description)
	if err != nil {
		return err
	}
	if minValue > maxValue {
		return errors.New("option {minValue} must be less than {maxValue}")
	}
	if len(choices) > 0 && autocomplete {
		return errors.New("option {choices} can only be used with {autocomplete} disabled")
	}
	cmd.options = append(cmd.options, option{
		Type:         integerOption,
		Name:         name,
		Description:  description,
		MinValue:     minValue,
		MaxValue:     maxValue,
		AutoComplete: autocomplete,
		Choices:      choices,
		Required:     required,
	})
	return nil
}

func (cmd *Command) OptionBOOLEAN(name string, description string, required bool) error {
	if len(cmd.options) == 25 {
		return errors.New("application command can only have max 25 options")
	}
	err := checkNameDesc(name, description)
	if err != nil {
		return err
	}
	cmd.options = append(cmd.options, option{
		Type:        booleanOption,
		Name:        name,
		Description: description,
		Required:    required,
	})
	return nil
}

func (cmd *Command) OptionUSER(name string, description string, required bool) error {
	if len(cmd.options) == 25 {
		return errors.New("application command can only have max 25 options")
	}
	err := checkNameDesc(name, description)
	if err != nil {
		return err
	}
	cmd.options = append(cmd.options, option{
		Type:        userOption,
		Name:        name,
		Description: description,
		Required:    required,
	})
	return nil
}

func (cmd *Command) OptionCHANNEL(name string, description string, required bool, channelTypes ...ChannelKind) error {
	if len(cmd.options) == 25 {
		return errors.New("application command can only have max 25 options")
	}
	err := checkNameDesc(name, description)
	if err != nil {
		return err
	}
	kinds := map[int]ChannelKind{
		0:  ChannelKinds.Text,
		1:  ChannelKinds.DM,
		2:  ChannelKinds.GroupDM,
		3:  ChannelKinds.Voice,
		4:  ChannelKinds.Category,
		5:  ChannelKinds.News,
		10: ChannelKinds.NewsThread,
		11: ChannelKinds.PublicThread,
		12: ChannelKinds.PrivateThread,
		13: ChannelKinds.StageVoice,
		14: ChannelKinds.Directory,
		15: ChannelKinds.Forum,
	}
	for _, ch := range channelTypes {
		if _, ok := kinds[int(ch)]; !ok {
			return errors.New(fmt.Sprintf("option {ChannelTypes} contains invalid channel kind (%d)", ch))
		}
	}
	cmd.options = append(cmd.options, option{
		Type:         channelOption,
		Name:         name,
		Description:  description,
		ChannelTypes: channelTypes,
		Required:     required,
	})
	return nil
}

func (cmd *Command) OptionROLE(name string, description string, required bool) error {
	if len(cmd.options) == 25 {
		return errors.New("application command can only have max 25 options")
	}
	err := checkNameDesc(name, description)
	if err != nil {
		return err
	}
	cmd.options = append(cmd.options, option{
		Type:        roleOption,
		Name:        name,
		Description: description,
		Required:    required,
	})
	return nil
}

func (cmd *Command) OptionMENTIONABLE(name string, description string, required bool) error {
	if len(cmd.options) == 25 {
		return errors.New("application command can only have max 25 options")
	}
	err := checkNameDesc(name, description)
	if err != nil {
		return err
	}
	cmd.options = append(cmd.options, option{
		Type:        mentionableOption,
		Name:        name,
		Description: description,
		Required:    required,
	})
	return nil
}

func (cmd *Command) OptionNUMBER(
	name string,
	description string,
	required bool,
	minValue float64,
	maxValue float64,
	autocomplete bool,
	choices ...Choice,
) error {
	if len(cmd.options) == 25 {
		return errors.New("application command can only have max 25 options")
	}
	err := checkNameDesc(name, description)
	if err != nil {
		return err
	}
	if minValue > maxValue {
		return errors.New("option {minValue} must be less than {maxValue}")
	}
	if len(choices) > 0 && autocomplete {
		return errors.New("option {choices} can only be used with {autocomplete} disabled")
	}
	cmd.options = append(cmd.options, option{
		Type:         numberOption,
		Name:         name,
		Description:  description,
		MinValueNum:  minValue,
		MaxValueNum:  maxValue,
		AutoComplete: autocomplete,
		Choices:      choices,
		Required:     required,
	})
	return nil
}

func (cmd *Command) OptionATTACHMENT(name string, description string, required bool) error {
	if len(cmd.options) == 25 {
		return errors.New("application command can only have max 25 options")
	}
	err := checkNameDesc(name, description)
	if err != nil {
		return err
	}
	cmd.options = append(cmd.options, option{
		Type:        attachmentOption,
		Name:        name,
		Description: description,
		Required:    required,
	})
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
	case CommandKinds.Message:
		body["type"] = int(CommandKinds.Message)
	case CommandKinds.User:
		body["type"] = int(CommandKinds.User)
	case CommandKinds.Slash:
		body["type"] = int(CommandKinds.Slash)
	default:
		body["type"] = int(CommandKinds.Slash)
		cmd.Type = CommandKinds.Slash
	}
	cmd.uniqueId = assignId()
	body["name"] = cmd.Name
	switch cmd.Type {
	case CommandKinds.Slash:
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
	case CommandKinds.User:
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
	case CommandKinds.Message:
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

func checkNameDesc(name string, desc string) error {
	if name == "" {
		return errors.New("option {name} must be set")
	}
	if len(name) > 32 {
		return errors.New(fmt.Sprintf("option (%s) {name} must be less than 32 characters", name))
	}
	if desc == "" {
		return errors.New(fmt.Sprintf("option (%s) {description} must be set", name))
	}
	if len(desc) > 100 {
		return errors.New(fmt.Sprintf("option (%s) {description} must be less than 100 characters", name))
	}
	return nil
}
