package discord

type optionKind int

const (
	stringOption      optionKind = 3
	integerOption     optionKind = 4
	booleanOption     optionKind = 5
	userOption        optionKind = 6
	channelOption     optionKind = 7
	roleOption        optionKind = 8
	mentionableOption optionKind = 9
	numberOption      optionKind = 10
	attachmentOption  optionKind = 11
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

type option struct {
	Name         string     `json:"name"`
	Type         optionKind `json:"type"`
	Description  string
	Required     bool
	MinLength    int           // allowed for: stringOption
	MaxLength    int           // allowed for: stringOption
	MinValue     int           // allowed for: integerOption
	MaxValue     int           // allowed for: integerOption
	MaxValueNum  float64       // allowed for: numberOption
	MinValueNum  float64       // allowed for: numberOption
	AutoComplete bool          // allowed for: stringOption, numberOption, integerOption
	ChannelTypes []ChannelKind // allowed for: channelOption
	Choices      []Choice      // allowed for: stringOption, integerOption, numberOption
	Value        any           `json:"Value"`   // available only during option parsing
	Focused      bool          `json:"focused"` // available only during option parsing
}

func (o *option) marshal() map[string]interface{} {
	body := map[string]interface{}{}
	body["name"] = o.Name
	body["type"] = o.Type
	body["required"] = o.Required
	body["description"] = o.Description
	switch o.Type {
	case stringOption:
		body["min_length"] = o.MinLength
		body["max_length"] = o.MaxLength
		if len(o.Choices) > 0 {
			body["choices"] = o.Choices
		} else if o.AutoComplete {
			body["auto_complete"] = true
		}
	case integerOption:
		body["min_value"] = o.MinValue
		body["max_value"] = o.MaxValue
		if len(o.Choices) > 0 {
			body["choices"] = o.Choices
		} else if o.AutoComplete {
			body["auto_complete"] = true
		}
	case numberOption:
		body["min_value"] = o.MinValueNum
		body["max_value"] = o.MaxValueNum
		if len(o.Choices) > 0 {
			body["choices"] = o.Choices
		} else if o.AutoComplete {
			body["auto_complete"] = true
		}
	case channelOption:
		body["channel_types"] = []int{}
		for _, c := range o.ChannelTypes {
			body["channel_types"] = append(body["channel_types"].([]int), int(c))
		}
	}
	return body
}
