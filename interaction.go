package disgo

type InteractionData struct {
	Id       string                 `json:"id"`
	Name     string                 `json:"name"`
	Type     int                    `json:"type"`
	Resolved map[string]interface{} `json:"resolved"`
	Options  []Option               `json:"options"`
	GuildId  string                 `json:"guild_id"`
	TargetId string                 `json:"target_id"`
}

type Interaction struct {
	Id             string          `json:"id"`
	ApplicationId  string          `json:"application_id"`
	Type           int             `json:"type"`
	Data           InteractionData `json:"data"`
	GuildId        string          `json:"guild_id"`
	ChannelId      string          `json:"channel_id"`
	Member         Member          `json:"member"`
	User           User            `json:"user"`
	Token          string          `json:"token"`
	Version        int             `json:"version"`
	Message        interface{}     `json:"message"`
	AppPermissions string          `json:"app_permissions"`
	Locale         string          `json:"locale"`
	GuildLocale    string          `json:"guild_locale"`
}
