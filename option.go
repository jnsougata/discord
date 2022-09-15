package discord

type Choice struct {
	Name  string      `json:"name"`
	Value interface{} `json:"Value"` // same type as type of Option
}

type Option struct {
	Name         string `json:"name"`
	Type         int    `json:"type"`
	Description  string
	Required     bool
	MinLength    int      // allowed for: string Option Range(0-6000)
	MaxLength    int      // allowed for: string Option Range(1-6000)
	MinValue     int      // allowed for: integer Option
	MaxValue     int      // allowed for: integer Option
	MaxValueNum  float64  // allowed for: number Option
	MinValueNum  float64  // allowed for: number Option
	AutoComplete bool     // allowed for: string Option, number Option, integer Option
	ChannelTypes []int    // allowed for: channel Option
	Choices      []Choice // allowed for: string Option, number Option, integer Option
	Value        any      `json:"Value"`   // do not use this field
	Focused      bool     `json:"focused"` // do not use this field
}

func (o *Option) marshal() map[string]interface{} {
	body := map[string]interface{}{}
	body["name"] = o.Name
	body["type"] = o.Type
	body["required"] = o.Required
	body["description"] = o.Description
	switch o.Type {
	case OptionTypes.String:
		body["min_length"] = o.MinLength
		body["max_length"] = o.MaxLength
		if len(o.Choices) > 0 {
			body["choices"] = o.Choices
		} else if o.AutoComplete {
			body["auto_complete"] = true
		}
	case OptionTypes.Integer:
		body["min_value"] = o.MinValue
		body["max_value"] = o.MaxValue
		if len(o.Choices) > 0 {
			body["choices"] = o.Choices
		} else if o.AutoComplete {
			body["auto_complete"] = true
		}
	case OptionTypes.Number:
		body["min_value"] = o.MinValueNum
		body["max_value"] = o.MaxValueNum
		if len(o.Choices) > 0 {
			body["choices"] = o.Choices
		} else if o.AutoComplete {
			body["auto_complete"] = true
		}
	case OptionTypes.Channel:
		body["channel_types"] = []int{}
		for _, c := range o.ChannelTypes {
			body["channel_types"] = append(body["channel_types"].([]int), c)
		}
	}
	return body
}
