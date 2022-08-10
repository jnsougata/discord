package disgo

import "strconv"

type Member struct {
	Nickname      string           `json:"nick"`
	AvatarHash    string           `json:"avatar"`
	JoinedAt      string           `json:"joined_at"`
	PremiumSince  string           `json:"premium_since"`
	Deaf          bool             `json:"deaf"`
	Mute          bool             `json:"mute"`
	Pending       bool             `json:"pending"`
	Permissions   string           `json:"permissions"`
	TimeoutExpiry string           `json:"communication_disabled_until"`
	GuildId       string           `json:"guild_id"`
	Roles         map[string]*Role `json:"roles"`
	token         string
	Id            string
	Name          string
	Discriminator string
	Avatar        Asset
	Bot           bool
	System        bool
	MfaEnabled    bool
	Banner        Asset
	Color         int
	Locale        string
	Verified      bool
	Email         string
	Flags         int
	PremiumType   int
	PublicFlags   int
}

func (m *Member) fillUser(u *User) {
	m.Id = u.Id
	m.Bot = u.Bot
	m.token = u.token
	m.System = u.System
	m.Banner = u.Banner
	m.Name = u.Username
	m.MfaEnabled = u.MfaEnabled
	m.Discriminator = u.Discriminator
}

func (m *Member) HasPermissions(permissions ...Permission) bool {
	total, _ := strconv.Atoi(m.Permissions)
	return Permissions().check(Permission(total), permissions...)
}
