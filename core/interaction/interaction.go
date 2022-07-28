package interaction

import (
	"encoding/json"
	"github.com/jnsougata/disgo/core/member"
	"github.com/jnsougata/disgo/core/user"
)

type Option struct {
	Name    string      `json:"name"`
	Type    int         `json:"type"`
	Value   interface{} `json:"value"`
	Options []Option    `json:"options"`
	Focused bool        `json:"focused"`
}

type Data struct {
	Id       string                 `json:"id"`
	Name     string                 `json:"name"`
	Type     int                    `json:"type"`
	Resolved map[string]interface{} `json:"resolved"`
	Options  []Option               `json:"options"`
	GuildId  string                 `json:"guild_id"`
	TargetId string                 `json:"target_id"`
}

type Interaction struct {
	Id             string        `json:"id"`
	ApplicationId  string        `json:"application_id"`
	Type           int           `json:"type"`
	Data           Data          `json:"data"`
	GuildId        string        `json:"guild_id"`
	ChannelId      string        `json:"channel_id"`
	Member         member.Member `json:"member"`
	User           user.User     `json:"user"`
	Token          string        `json:"token"`
	Version        int           `json:"version"`
	Message        interface{}   `json:"message"`
	AppPermissions string        `json:"app_permissions"`
	Locale         string        `json:"locale"`
	GuildLocale    string        `json:"guild_locale"`
}

func FromData(payload interface{}) *Interaction {
	i := &Interaction{}
	data, _ := json.Marshal(payload)
	_ = json.Unmarshal(data, i)
	return i
}
