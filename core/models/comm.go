package models

// type Option map[string]interface{}

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

type SlashCommand struct {
	Name                     string   `json:"name"`
	Description              string   `json:"description"`
	Options                  []Option `json:"options,omitempty"`
	DefaultMemberPermissions int      `json:"default_member_permissions,omitempty"`
	DMPermission             bool     `json:"dm_permission,omitempty"`
	TestGuildId              int64    `json:"guild_id,string"`
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
	ChannelTypes []string `json:"channel_types,omitempty"`
	Choices      []Choice `json:"choices,omitempty"`
}

type Choice struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}
