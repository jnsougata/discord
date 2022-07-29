package member

import (
	"encoding/json"
	"github.com/jnsougata/disgo/user"
)

type Member struct {
	User          user.User `json:"user"`
	Nickname      string    `json:"nick"`
	Roles         []string  `json:"roles"`
	JoinedAt      string    `json:"joined_at"`
	PremiumSince  string    `json:"premium_since"`
	Deaf          bool      `json:"deaf"`
	Mute          bool      `json:"mute"`
	Pending       bool      `json:"pending"`
	Permissions   int       `json:"permissions"`
	TimeoutExpiry string    `json:"communication_disabled_until"`
	GuildId       string    `json:"guild_id"`
}

func Unmarshal(payload interface{}) *Member {
	m := &Member{}
	data, _ := json.Marshal(payload)
	_ = json.Unmarshal(data, m)
	return m
}
