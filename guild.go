package disgo

import (
	"encoding/json"
)

type Guild struct {
	Id                          string                   `json:"id"`
	Name                        string                   `json:"name"`
	Icon                        string                   `json:"icon"`
	IconHash                    string                   `json:"icon_hash"`
	Splash                      string                   `json:"splash"`
	DiscoverySplash             string                   `json:"discovery_splash"`
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
	Roles                       map[string]Role          `json:"x_roles"`
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
	Banner                      string                   `json:"banner"`
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
	Members                     map[string]Member        `json:"x_members"`
	Channels                    map[string]Channel       `json:"x_channels"`
	JoinedAT                    string                   `json:"joined_at"`
	Large                       bool                     `json:"large"`
	MemberCount                 int                      `json:"member_count"`
	VoiceStates                 []map[string]interface{} `json:"voice_states"`
	Presences                   []map[string]interface{} `json:"presences"`
	Threads                     []map[string]interface{} `json:"threads"`
	StageInstances              []map[string]interface{} `json:"stage_instances"`
	Unavailable                 bool                     `json:"unavailable"`
	GuildScheduledEvents        []map[string]interface{} `json:"guild_scheduled_events"`
}

func unmarshalGuild(payload interface{}) *Guild {
	guild := &Guild{}
	data, _ := json.Marshal(payload)
	_ = json.Unmarshal(data, guild)
	return guild
}

func (guild *Guild) unmarshalMembers(objs []interface{}) {
	var members = map[string]Member{}
	for _, o := range objs {
		uo := unmarshalMember(o)
		uo.GuildId = guild.Id
		members[uo.User.Id] = *uo
	}
	guild.Members = members
}

func (guild *Guild) unmarshalRoles(objs []interface{}) {
	var roles = map[string]Role{}
	for _, o := range objs {
		uo := unmarshalRole(o)
		uo.GuildId = guild.Id
		roles[uo.Id] = *uo
	}
	guild.Roles = roles
}

func (guild *Guild) unmarshalChannels(objs []interface{}) {
	var chs = map[string]Channel{}
	for _, o := range objs {
		uo := unmarshalChannel(o)
		chs[uo.Id] = *uo
	}
	guild.Channels = chs
}
