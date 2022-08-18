package discord

type intent int

type intents struct {
	Guilds                intent
	GuildMembers          intent
	GuildBans             intent
	GuildEmojis           intent
	GuildIntegrations     intent
	GuildWebhooks         intent
	GuildInvites          intent
	GuildVoiceStates      intent
	GuildPresences        intent
	GuildMessages         intent
	GuildMessageReactions intent
	GuildMessageTyping    intent
	DM                    intent
	DMReactions           intent
	DMTyping              intent
	MessageContent        intent
	GuildScheduleEvents   intent
	AutoModConfiguration  intent
	AutoModExecution      intent
}

var Intents = intents{
	Guilds:                intent(1 << 0),
	GuildMembers:          intent(1 << 1),
	GuildBans:             intent(1 << 2),
	GuildEmojis:           intent(1 << 3),
	GuildIntegrations:     intent(1 << 4),
	GuildWebhooks:         intent(1 << 5),
	GuildInvites:          intent(1 << 6),
	GuildVoiceStates:      intent(1 << 7),
	GuildPresences:        intent(1 << 8),
	GuildMessages:         intent(1 << 9),
	GuildMessageReactions: intent(1 << 10),
	GuildMessageTyping:    intent(1 << 11),
	DM:                    intent(1 << 12),
	DMReactions:           intent(1 << 13),
	DMTyping:              intent(1 << 14),
	MessageContent:        intent(1 << 15),
	GuildScheduleEvents:   intent(1 << 16),
	AutoModConfiguration:  intent(1 << 17),
	AutoModExecution:      intent(1 << 18),
}

func (i *intents) Make(intents ...intent) intent {
	ini := 0
	for _, intent := range intents {
		ini |= 1 << int(intent)
	}
	return intent(ini)
}

func (i *intents) All() intent {
	ini := 0
	for i := 0; i < 19; i++ {
		ini |= 1 << i
	}
	return intent(ini)
}
