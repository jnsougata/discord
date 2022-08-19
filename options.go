package discord

type OptionType int

const (
	StringOption      OptionType = 3
	IntegerOption     OptionType = 4
	BooleanOption     OptionType = 5
	UserOption        OptionType = 6
	ChannelOption     OptionType = 7
	RoleOption        OptionType = 8
	MentionableOption OptionType = 9
	NumberOption      OptionType = 10
	AttachmentOption  OptionType = 11
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

type Choice struct {
	Name  string      `json:"name"`
	Value interface{} `json:"Value"` // same type as type of option
}

type Option struct {
	Name         string     `json:"name"`
	Type         OptionType `json:"type"`
	Value        any        `json:"Value"`   // available only during option parsing
	Focused      bool       `json:"focused"` // available only during option parsing
	Description  string
	Required     bool
	MinLength    int           // allowed for: StringOption
	MaxLength    int           // allowed for: StringOption
	MinValue     int64         // allowed for: IntegerOption, NumberOption
	MaxValue     int64         // allowed for: IntegerOption, NumberOption
	AutoComplete bool          // allowed for: StringOption, NumberOption, IntegerOption
	ChannelTypes []ChannelType // allowed for: ChannelOption
	Choices      []Choice      // allowed for: StringOption, IntegerOption, NumberOption
}

func (o *Option) marshal() map[string]interface{} {
	body := map[string]interface{}{}
	body["name"] = o.Name
	body["type"] = o.Type
	body["required"] = o.Required
	body["description"] = o.Description
	switch o.Type {
	case StringOption:
		body["min_length"] = o.MinLength
		body["max_length"] = o.MaxLength
		if len(o.Choices) > 0 {
			body["choices"] = o.Choices
		} else if o.AutoComplete {
			body["auto_complete"] = true
		}
	case IntegerOption:
		body["min_value"] = o.MinValue
		body["max_value"] = o.MaxValue
		if len(o.Choices) > 0 {
			body["choices"] = o.Choices
		} else if o.AutoComplete {
			body["auto_complete"] = true
		}
	case NumberOption:
		body["min_value"] = o.MinValue
		body["max_value"] = o.MaxValue
		if len(o.Choices) > 0 {
			body["choices"] = o.Choices
		} else if o.AutoComplete {
			body["auto_complete"] = true
		}
	case ChannelOption:
		body["channel_types"] = []int{}
		for _, c := range o.ChannelTypes {
			body["channel_types"] = append(body["channel_types"].([]int), int(c))
		}
	}
	return body
}
