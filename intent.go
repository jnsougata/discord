package discord

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

func Intents(intents ...Intent) Intent {
	base := 0
	if len(intents) > 0 {
		for _, intent := range intents {
			base |= int(intent)
		}
	} else {
		for i := 0; i < 19; i++ {
			base |= 1 << i
		}
	}
	return Intent(base)
}
