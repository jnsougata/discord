package bot

import (
	"github.com/jnsougata/disgo/core/client"
	"github.com/jnsougata/disgo/core/command"
	"github.com/jnsougata/disgo/core/guild"
	"github.com/jnsougata/disgo/core/interaction"
	"github.com/jnsougata/disgo/core/message"
	"github.com/jnsougata/disgo/core/user"
)

const (
	onReady = "READY"
	// onResumed                      = "RESUMED"
	// onReconnect                 	  = "RECONNECT" --> handle internally
	// onInvalidSession 		      = "INVALID_SESSION" --> handle internally
	// onAppCmdPermsUpdate            = "APPLICATION_COMMAND_PERMISSIONS_UPDATE"
	// onAutoModRuleCreate            = "AUTO_MODERATION_RULE_CREATE"
	// onAutoModRuleDelete            = "AUTO_MODERATION_RULE_DELETE"
	// onAutoModRuleUpdate            = "AUTO_MODERATION_RULE_UPDATE"
	// onAutoModActionExec            = "AUTO_MODERATION_ACTION_EXECUTION"
	// onChannelCreate                = "CHANNEL_CREATE"
	// onChannelUpdate                = "CHANNEL_UPDATE"
	// onChannelDelete                = "CHANNEL_DELETE"
	// onChannelPinsUpdate            = "CHANNEL_PINS_UPDATE"
	// onThreadCreate                 = "THREAD_CREATE"
	// onThreadUpdate                 = "THREAD_UPDATE"
	// onThreadDelete                 = "THREAD_DELETE"
	// onThreadListSync               = "THREAD_LIST_SYNC"
	// onThreadMembersUpdate          = "THREAD_MEMBERS_UPDATE"
	onGuildCreate = "GUILD_CREATE"
	// onGuildUpdate                  = "GUILD_UPDATE"
	onGuildDelete = "GUILD_DELETE"
	// onGuildBanAdd                  = "GUILD_BAN_ADD"
	// onGuildBanRemove               = "GUILD_BAN_REMOVE"
	// onGuildEmojisUpdate            = "GUILD_EMOJIS_UPDATE"
	// onGuildStickersUpdate          = "GUILD_STICKERS_UPDATE"
	// onGuildIntegrationsUpdate      = "GUILD_INTEGRATIONS_UPDATE"
	// onGuildMemberAdd               = "GUILD_MEMBER_ADD"
	// onGuildMemberRemove            = "GUILD_MEMBER_REMOVE"
	// onGuildMemberUpdate            = "GUILD_MEMBER_UPDATE"
	// onGuildMembersChunk            = "GUILD_MEMBERS_CHUNK"
	// onGuildRoleCreate              = "GUILD_ROLE_CREATE"
	// onGuildRoleUpdate              = "GUILD_ROLE_UPDATE"
	// onGuildRoleDelete              = "GUILD_ROLE_DELETE"
	// onGuildScheduleEventCreate     = "GUILD_SCHEDULE_EVENT_CREATE"
	// onGuildScheduleEventUpdate     = "GUILD_SCHEDULE_EVENT_UPDATE"
	// onGuildScheduleEventDelete     = "GUILD_SCHEDULE_EVENT_DELETE"
	// onGuildScheduleEventUserAdd    = "GUILD_SCHEDULE_EVENT_USER_ADD"
	// onGuildScheduleEventUserRemove = "GUILD_SCHEDULE_EVENT_USER_REMOVE"
	// onIntegrationCreate            = "INTEGRATION_CREATE"
	// onIntegrationUpdate            = "INTEGRATION_UPDATE"
	// onIntegrationDelete            = "INTEGRATION_DELETE"
	onInteractionCreate = "INTERACTION_CREATE"
	// onInviteCreate                 = "INVITE_CREATE"
	// onInviteDelete                 = "INVITE_DELETE"
	onMessageCreate = "MESSAGE_CREATE"
	// onMessageUpdate                = "MESSAGE_UPDATE"
	// onMessageDelete                = "MESSAGE_DELETE"
	// onMessageDeleteBulk            = "MESSAGE_DELETE_BULK"
	// onMessageReactionAdd           = "MESSAGE_REACTION_ADD"
	// onMessageReactionRemove        = "MESSAGE_REACTION_REMOVE"
	// onMessageReactionRemoveAll     = "MESSAGE_REACTION_REMOVE_ALL"
	// onMessageReactionRemoveEmoji   = "MESSAGE_REACTION_REMOVE_EMOJI"
	// onPresenceUpdate               = "PRESENCE_UPDATE"
	// onStageInstanceCreate          = "STAGE_INSTANCE_CREATE"
	// onStageInstanceUpdate          = "STAGE_INSTANCE_UPDATE"
	// onStageInstanceDelete          = "STAGE_INSTANCE_DELETE"
	// onTypingStart                  = "TYPING_START"
	// onUserUpdate                   = "USER_UPDATE"
	// onVoiceStateUpdate             = "VOICE_STATE_UPDATE"
	// onVoiceServerUpdate            = "VOICE_SERVER_UPDATE"
	// onWebhooksUpdate               = "WEBHOOKS_UPDATE"
)

type Bot struct {
	intent int
	core   *client.Client
}

func New(intent int, memoize bool) *Bot {
	return &Bot{
		intent: intent,
		core: &client.Client{
			Intent:  intent,
			Memoize: memoize,
		},
	}
}

func (bot *Bot) Run(token string) {
	bot.core.Run(token)
}

func (bot *Bot) AddCommand(
	command command.SlashCommand, handler func(bot user.User, interaction interaction.Interaction, options ...interaction.Option)) {
	bot.core.Queue(command, handler)
}

func (bot *Bot) OnMessage(handler func(bot user.User, message message.Message)) {
	bot.core.AddHandler(onMessageCreate, handler)
}

func (bot *Bot) OnReady(handler func(bot user.User)) {
	bot.core.AddHandler(onReady, handler)
}

func (bot *Bot) OnInteraction(handler func(bot user.User, interaction interaction.Interaction)) {
	bot.core.AddHandler(onInteractionCreate, handler)
}

func (bot *Bot) OnGuildJoin(handler func(bot user.User, guild guild.Guild)) {
	bot.core.AddHandler(onGuildCreate, handler)
}

func (bot *Bot) OnGuildLeave(handler func(bot user.User, guild guild.Guild)) {
	bot.core.AddHandler(onGuildDelete, handler)
}
