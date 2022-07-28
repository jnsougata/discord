package command

import (
	"fmt"
	"github.com/jnsougata/disgo/core/interaction"
	"github.com/jnsougata/disgo/core/user"
)

const (
	SlashCommandType    = 1
	UserCommandType     = 2
	MessageCommandType  = 3
	SubCommandType      = 1
	SubCommandGroupType = 2
	StringType          = 3
	IntegerType         = 4
	BooleanType         = 5
	UserType            = 6
	ChannelType         = 7
	RoleType            = 8
	MentionableType     = 9
	NumberType          = 10
	AttachmentType      = 11
)

const (
	GuildTextChannel   = 0
	DMChannel          = 1
	GuildVoiceChannel  = 2
	GroupDMChannel     = 3
	GuildCategory      = 4
	GuildNews          = 5
	GuildNewsThread    = 10
	GuildPublicThread  = 11
	GuildPrivateThread = 12
	GuildStageVoice    = 13
	GuildDirectory     = 14
	GuildForum         = 15
)

type ApplicationCommand struct {
	Type                     int
	Name                     string
	Description              string
	Options                  []Option
	DMPermission             bool
	DefaultMemberPermissions int
	GuildId                  int64
	Handler                  func(bot user.User, ctx Context, options ...interaction.Option)
}

func (cmd *ApplicationCommand) ToData() (
	map[string]interface{},
	func(bot user.User, ctx Context, options ...interaction.Option),
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
	if cmd.DefaultMemberPermissions != 0 {
		body["default_member_permissions"] = cmd.DefaultMemberPermissions
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
	Type         int      `json:"type"`
	Required     bool     `json:"required,omitempty"`
	MinLength    int      `json:"min_length,omitempty"`
	MaxLength    int      `json:"max_length,omitempty"`
	MinValue     int64    `json:"min_value,omitempty"`
	MaxValue     int64    `json:"max_value,omitempty"`
	AutoComplete bool     `json:"auto_complete,omitempty"`
	ChannelTypes []int    `json:"channel_types,omitempty"`
	Options      []Option `json:"options,omitempty"`
	Choices      []Choice `json:"choices,omitempty"`
}

type Choice struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}
