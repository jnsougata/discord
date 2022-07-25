package command

const (
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

type SlashCommand struct {
	Name                     string   `json:"name"`
	Description              string   `json:"description"`
	Options                  []Option `json:"options,omitempty"`
	DefaultMemberPermissions int      `json:"default_member_permissions,omitempty"`
	DMPermission             bool     `json:"dm_permission,omitempty"`
	TestGuildId              int64    `json:"guild_id,string"`
}

func New(
	name string, description string,
	defaultMemberPermissions int, dmInvoke bool,
	guildId int64, options ...Option) SlashCommand {
	return SlashCommand{
		Name:                     name,
		Description:              description,
		Options:                  options,
		DMPermission:             dmInvoke,
		TestGuildId:              guildId,
		DefaultMemberPermissions: defaultMemberPermissions,
	}
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
