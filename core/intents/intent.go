package intents

const (
	GuildsIntent                = 1 << 0
	GuildMembersIntent          = 1 << 1
	GuildBansIntent             = 1 << 2
	GuildEmojisIntent           = 1 << 3
	GuildIntegrationsIntent     = 1 << 4
	GuildWebhooksIntent         = 1 << 5
	GuildInvitesIntent          = 1 << 6
	GuildVoiceStatesIntent      = 1 << 7
	GuildPresencesIntent        = 1 << 8
	GuildMessagesIntent         = 1 << 9
	GuildMessageReactionsIntent = 1 << 10
	GuildMessageTypingIntent    = 1 << 11
	DMIntent                    = 1 << 12
	DMReactionsIntent           = 1 << 13
	DMTypingIntent              = 1 << 14
	MessageContentIntent        = 1 << 15
	GuildScheduleEventsIntent   = 1 << 16
	AutoModConfigIntent         = 1 << 17
	AutoModExecuteIntent        = 1 << 18
)

type Intent struct {
	Guilds                bool
	GuildMembers          bool
	GuildBans             bool
	GuildEmojis           bool
	GuildIntegrations     bool
	GuildWebhooks         bool
	GuildInvites          bool
	GuildVoiceStates      bool
	GuildPresences        bool
	GuildMessages         bool
	GuildMessageReactions bool
	GuildMessageTyping    bool
	DM                    bool
	DMReactions           bool
	DMTyping              bool
	MessageContent        bool
	GuildScheduleEvents   bool
	AutoModConfig         bool
	AutoModExecute        bool
}

func (i Intent) Custom() int {
	base := 0
	if i.Guilds {
		base |= GuildsIntent
	}
	if i.GuildMembers {
		base |= GuildMembersIntent
	}
	if i.GuildBans {
		base |= GuildBansIntent
	}
	if i.GuildEmojis {
		base |= GuildEmojisIntent
	}
	if i.GuildIntegrations {
		base |= GuildIntegrationsIntent
	}
	if i.GuildWebhooks {
		base |= GuildWebhooksIntent
	}
	if i.GuildInvites {
		base |= GuildInvitesIntent
	}
	if i.GuildVoiceStates {
		base |= GuildVoiceStatesIntent
	}
	if i.GuildPresences {
		base |= GuildPresencesIntent
	}
	if i.GuildMessages {
		base |= GuildMessagesIntent
	}
	if i.GuildMessageReactions {
		base |= GuildMessageReactionsIntent
	}
	if i.GuildMessageTyping {
		base |= GuildMessageTypingIntent
	}
	if i.DM {
		base |= DMIntent
	}
	if i.DMReactions {
		base |= DMReactionsIntent
	}
	if i.DMTyping {
		base |= DMTypingIntent
	}
	if i.MessageContent {
		base |= MessageContentIntent
	}
	if i.GuildScheduleEvents {
		base |= GuildScheduleEventsIntent
	}
	if i.AutoModConfig {
		base |= AutoModConfigIntent
	}
	if i.AutoModExecute {
		base |= AutoModExecuteIntent
	}
	return base
}

func All() int {
	return GuildsIntent | GuildMembersIntent | GuildBansIntent | GuildEmojisIntent | GuildIntegrationsIntent | GuildWebhooksIntent | GuildInvitesIntent | GuildVoiceStatesIntent | GuildPresencesIntent | GuildMessagesIntent | GuildMessageReactionsIntent | GuildMessageTypingIntent | DMIntent | DMReactionsIntent | DMTypingIntent | MessageContentIntent | AutoModConfigIntent | AutoModExecuteIntent
}

func Basic() int {
	return GuildsIntent | GuildBansIntent | GuildEmojisIntent | GuildIntegrationsIntent | GuildWebhooksIntent | GuildInvitesIntent | GuildVoiceStatesIntent | GuildMessagesIntent | GuildMessageReactionsIntent | GuildMessageTypingIntent | DMIntent | DMReactionsIntent | DMTypingIntent | AutoModConfigIntent | AutoModExecuteIntent | GuildScheduleEventsIntent
}

func Exclusive() int {
	return GuildMembersIntent | GuildPresencesIntent | MessageContentIntent
}
