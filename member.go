package disgo

type Member struct {
	//User          User     `json:"user"`
	Nickname      string   `json:"nick"`
	AvatarHash    string   `json:"avatar"`
	Roles         []string `json:"roles"`
	JoinedAt      string   `json:"joined_at"`
	PremiumSince  string   `json:"premium_since"`
	Deaf          bool     `json:"deaf"`
	Mute          bool     `json:"mute"`
	Pending       bool     `json:"pending"`
	Permissions   int      `json:"permissions"`
	TimeoutExpiry string   `json:"communication_disabled_until"`
	GuildId       string   `json:"guild_id"`
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
	m.token = u.token
	m.Id = u.Id
	m.Name = u.Username
	m.Discriminator = u.Discriminator
	m.Bot = u.Bot
	m.System = u.System
	m.MfaEnabled = u.MfaEnabled
	m.Banner = u.Banner
}
