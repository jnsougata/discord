package discord

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
	state         *state
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
	m.state = u.state
	m.System = u.System
	m.Banner = u.Banner
	m.Name = u.Username
	m.MfaEnabled = u.MfaEnabled
	m.Discriminator = u.Discriminator
	m.Email = u.Email
	m.Color = u.Color
	m.Locale = u.Locale
	m.Verified = u.Verified
	m.Flags = u.Flags
	m.PremiumType = u.PremiumType
	m.PublicFlags = u.PublicFlags
}

func (m *Member) HasPermissions(permissions ...Permission) bool {
	p, _ := strconv.Atoi(m.Permissions)
	return Permissions.Check(Permission(p), permissions...)
}
