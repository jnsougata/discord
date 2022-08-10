package disgo

type Guild struct {
	Id                          string                   `json:"id"`
	Name                        string                   `json:"name"`
	Owner                       bool                     `json:"owner"`
	OwnerID                     string                   `json:"owner_id"`
	Permissions                 string                   `json:"permissions"`
	Region                      string                   `json:"region"`
	AfkChannelID                string                   `json:"afk_channel_id"`
	AfkTimeout                  int                      `json:"afk_timeout"`
	WidgetEnabled               bool                     `json:"widget_enabled"`
	WidgetChannelID             string                   `json:"widget_channel_id"`
	VerificationLevel           int                      `json:"verification_level"`
	DefaultMessageNotifications int                      `json:"default_message_notifications"`
	ExplicitContentFilter       int                      `json:"explicit_content_filter"`
	Emojis                      []Emoji                  `json:"emojis"`
	Features                    []string                 `json:"features"`
	MFALevel                    int                      `json:"mfa_level"`
	ApplicationID               string                   `json:"application_id"`
	SystemChannelID             string                   `json:"system_channel_id"`
	SystemChannelFlags          int                      `json:"system_channel_flags"`
	RulesChannelID              string                   `json:"rules_channel_id"`
	MaxPresences                int                      `json:"max_presences"`
	MaxMembers                  int                      `json:"max_members"`
	VanityURLCode               string                   `json:"vanity_url_code"`
	Description                 string                   `json:"description"`
	PremiumTier                 int                      `json:"premium_tier"`
	PremiumSubscriptionCount    int                      `json:"premium_subscription_count"`
	PreferredLocale             string                   `json:"preferred_locale"`
	PublicUpdatesChannelID      string                   `json:"public_updates_channel_id"`
	MaxVideoChannelUsers        int                      `json:"max_video_channel_users"`
	ApproximateMemberCount      int                      `json:"approximate_member_count"`
	ApproximatePresenceCount    int                      `json:"approximate_presence_count"`
	WelcomeScreen               map[string]interface{}   `json:"welcome_screen_enabled"`
	NSFWLevel                   int                      `json:"nsfw_level"`
	Stickers                    map[string]interface{}   `json:"stickers"`
	PremiumProgressBarEnabled   bool                     `json:"premium_progress_bar_enabled"`
	JoinedAT                    string                   `json:"joined_at"`
	Large                       bool                     `json:"large"`
	MemberCount                 int                      `json:"member_count"`
	VoiceStates                 []map[string]interface{} `json:"voice_states"`
	Presences                   []map[string]interface{} `json:"presences"`
	Threads                     []map[string]interface{} `json:"threads"`
	StageInstances              []map[string]interface{} `json:"stage_instances"`
	Unavailable                 bool                     `json:"unavailable"`
	GuildScheduledEvents        []map[string]interface{} `json:"guild_scheduled_events"`
	token                       string
	clientId                    string
	Icon                        Asset
	Banner                      Asset
	Splash                      Asset //    `json:"splash"`
	DiscoverySplash             Asset //    `json:"discovery_splash"`
	Members                     map[string]*Member
	Channels                    map[string]*Channel
	Roles                       map[string]*Role
	Me                          *Member
}

func (guild *Guild) unmarshalMembers(objs []interface{}) {
	var members = map[string]*Member{}
	for _, o := range objs {
		uo := Converter{payload: o, token: guild.token}.Member()
		uo.GuildId = guild.Id
		members[uo.Id] = uo
	}
	guild.Members = members
	guild.Me = guild.Members[guild.clientId]
}

func (guild *Guild) unmarshalRoles(objs []interface{}) {
	var roles = map[string]*Role{}
	for _, o := range objs {
		uo := Converter{payload: o, token: guild.token}.Role()
		uo.GuildId = guild.Id
		roles[uo.Id] = uo
	}
	guild.Roles = roles
}

func (guild *Guild) unmarshalChannels(objs []interface{}) {
	var channels = map[string]*Channel{}
	for _, o := range objs {
		uo := Converter{payload: o, token: guild.token}.Channel()
		channels[uo.Id] = uo
	}
	guild.Channels = channels
}
