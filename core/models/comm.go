package models

type JSONMap map[string]interface{}

type SlashCommand struct {
	Name                     string    `json:"name"`
	Description              string    `json:"description"`
	Options                  []JSONMap `json:"options"`
	DefaultMemberPermissions int       `json:"default_member_permissions"`
	DMPermission             bool      `json:"dm_permission"`
	TestGuildId              int64     `json:"guild_id,string"`
}

type Option struct{}

func (o Option) String(
	name string, description string, required bool,
	minLength int, maxLength int, autoComplete bool, choices ...Choice) JSONMap {
	return JSONMap{
		"type":          3,
		"name":          name,
		"description":   description,
		"required":      required,
		"min_length":    minLength,
		"max_length":    maxLength,
		"auto_complete": autoComplete,
		"choices":       choices,
	}
}

func (o Option) Integer(
	name string, description string, required bool,
	minValue int64, maxValue int64, autoComplete bool, choices ...Choice) JSONMap {
	return JSONMap{
		"type":          4,
		"name":          name,
		"description":   description,
		"required":      required,
		"min_value":     minValue,
		"max_value":     maxValue,
		"auto_complete": autoComplete,
		"choices":       choices,
	}
}

func (o Option) Number(
	name string, description string, required bool,
	minValue float64, maxValue float64, autoComplete bool, choices ...Choice) JSONMap {
	return JSONMap{
		"type":          10,
		"name":          name,
		"description":   description,
		"required":      required,
		"min_value":     minValue,
		"max_value":     maxValue,
		"auto_complete": autoComplete,
		"choices":       choices,
	}
}

func (o Option) Boolean(
	name string, description string, required bool) JSONMap {
	return JSONMap{
		"type":        5,
		"name":        name,
		"description": description,
		"required":    required,
	}
}

func (o Option) User(
	name string, description string, required bool) JSONMap {
	return JSONMap{
		"type":        6,
		"name":        name,
		"description": description,
		"required":    required,
	}
}

func (o Option) Channel(
	name string, description string, required bool, channelTypes ...int64) JSONMap {
	return JSONMap{
		"type":          7,
		"name":          name,
		"description":   description,
		"required":      required,
		"channel_types": channelTypes,
	}
}

func (o Option) Role(
	name string, description string, required bool) JSONMap {
	return JSONMap{
		"type":        8,
		"name":        name,
		"description": description,
		"required":    required,
	}
}

func (o Option) Mentionable(
	name string, description string, required bool) JSONMap {
	return JSONMap{
		"type":        9,
		"name":        name,
		"description": description,
		"required":    required,
	}
}

func (o Option) Attachment(
	name string, description string, required bool) JSONMap {
	return JSONMap{
		"type":        11,
		"name":        name,
		"description": description,
		"required":    required,
	}
}

type Choice struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}
