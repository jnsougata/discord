package discord

type OptionKind int

const (
	StringOption      OptionKind = 3
	IntegerOption     OptionKind = 4
	BooleanOption     OptionKind = 5
	UserOption        OptionKind = 6
	ChannelOption     OptionKind = 7
	RoleOption        OptionKind = 8
	MentionableOption OptionKind = 9
	NumberOption      OptionKind = 10
	AttachmentOption  OptionKind = 11
)

type ChannelKind int

type channelKinds struct {
	Text          ChannelKind
	DM            ChannelKind
	Voice         ChannelKind
	GroupDM       ChannelKind
	Category      ChannelKind
	News          ChannelKind
	NewsThread    ChannelKind
	PublicThread  ChannelKind
	PrivateThread ChannelKind
	StageVoice    ChannelKind
	Directory     ChannelKind
	Forum         ChannelKind
}

var ChannelKinds = channelKinds{
	Text:          ChannelKind(0),
	DM:            ChannelKind(1),
	Voice:         ChannelKind(2),
	GroupDM:       ChannelKind(3),
	Category:      ChannelKind(4),
	News:          ChannelKind(5),
	NewsThread:    ChannelKind(10),
	PublicThread:  ChannelKind(11),
	PrivateThread: ChannelKind(12),
	StageVoice:    ChannelKind(13),
	Directory:     ChannelKind(14),
	Forum:         ChannelKind(15),
}

type Choice struct {
	Name  string      `json:"name"`
	Value interface{} `json:"Value"` // same type as type of option
}

type Option struct {
	Name         string     `json:"name"`
	Type         OptionKind `json:"type"`
	Description  string
	Required     bool
	MinLength    int           // allowed for: StringOption
	MaxLength    int           // allowed for: StringOption
	MinValue     int64         // allowed for: IntegerOption, NumberOption
	MaxValue     int64         // allowed for: IntegerOption, NumberOption
	AutoComplete bool          // allowed for: StringOption, NumberOption, IntegerOption
	ChannelTypes []ChannelKind // allowed for: ChannelOption
	Choices      []Choice      // allowed for: StringOption, IntegerOption, NumberOption
	Value        any           `json:"Value"`   // available only during option parsing
	Focused      bool          `json:"focused"` // available only during option parsing
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
