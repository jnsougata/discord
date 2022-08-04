package disgo

type Intent int

const (
	GuildsIntent                Intent = 1 << 0
	GuildMembersIntent          Intent = 1 << 1
	GuildBansIntent             Intent = 1 << 2
	GuildEmojisIntent           Intent = 1 << 3
	GuildIntegrationsIntent     Intent = 1 << 4
	GuildWebhooksIntent         Intent = 1 << 5
	GuildInvitesIntent          Intent = 1 << 6
	GuildVoiceStatesIntent      Intent = 1 << 7
	GuildPresencesIntent        Intent = 1 << 8
	GuildMessagesIntent         Intent = 1 << 9
	GuildMessageReactionsIntent Intent = 1 << 10
	GuildMessageTypingIntent    Intent = 1 << 11
	DMIntent                    Intent = 1 << 12
	DMReactionsIntent           Intent = 1 << 13
	DMTypingIntent              Intent = 1 << 14
	MessageContentIntent        Intent = 1 << 15
	GuildScheduleEventsIntent   Intent = 1 << 16
	AutoModConfigIntent         Intent = 1 << 17
	AutoModExecuteIntent        Intent = 1 << 18
)

type Intents struct {
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

func (i Intents) Build(intents ...Intent) int {
	base := 0
	for _, intent := range intents {
		base |= int(intent)
	}
	return base
}

func (i Intents) All() int {
	base := 0
	for i := 0; i < 19; i++ {
		base |= 1 << i
	}
	return base
}
