package discord

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
