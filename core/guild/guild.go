package guild

import "encoding/json"

type Guild struct {
	ID                          string                   `json:"id"`
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
	Roles                       []map[string]interface{} `json:"roles"`
	Emojis                      []map[string]interface{} `json:"emojis"`
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
}

func FromData(payload interface{}) *Guild {
	guild := &Guild{}
	data, _ := json.Marshal(payload)
	_ = json.Unmarshal(data, guild)
	return guild
}
