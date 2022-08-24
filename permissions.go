package discord

type Permission int

type permissions struct {
	CreateInstantInvite     Permission
	KickMembers             Permission
	BanMembers              Permission
	Administrator           Permission
	ManageChannels          Permission
	ManageGuild             Permission
	AddReactions            Permission
	ViewAuditLog            Permission
	PrioritySpeaker         Permission
	Stream                  Permission
	ViewChannel             Permission
	SendMessages            Permission
	SendTTSMessages         Permission
	ManageMessages          Permission
	EmbedLinks              Permission
	AttachFiles             Permission
	ReadMessageHistory      Permission
	MentionEveryone         Permission
	UseExternalEmojis       Permission
	ViewGuildInsights       Permission
	Connect                 Permission
	Speak                   Permission
	MuteMembers             Permission
	DeafenMembers           Permission
	MoveMembers             Permission
	UseVAD                  Permission
	ChangeNickname          Permission
	ManageNicknames         Permission
	ManageRoles             Permission
	ManageWebhooks          Permission
	ManageEmojisAndStickers Permission
	UseApplicationCommands  Permission
	RequestToSpeak          Permission
	ManageEvents            Permission
	ManageThreads           Permission
	CreatePublicThreads     Permission
	CreatePrivateThreads    Permission
	UseExternalStickers     Permission
	SendMessageInThreads    Permission
	UseEmbeddedActivity     Permission
	ModerateMembers         Permission
	All                     Permission
}

func (p *permissions) Check(reference Permission, others ...Permission) bool {
	check := true
	for _, o := range others {
		check = check && int(reference)&int(o) == int(o)
	}
	return check
}

func perms() *permissions {
	whole := 0
	for n := 0; n < 41; n++ {
		whole |= 1 << n
	}
	return &permissions{
		All:                     Permission(whole),
		CreateInstantInvite:     Permission(1 << 0),
		KickMembers:             Permission(1 << 1),
		BanMembers:              Permission(1 << 2),
		Administrator:           Permission(1 << 3),
		ManageChannels:          Permission(1 << 4),
		ManageGuild:             Permission(1 << 5),
		AddReactions:            Permission(1 << 6),
		ViewAuditLog:            Permission(1 << 7),
		PrioritySpeaker:         Permission(1 << 8),
		Stream:                  Permission(1 << 9),
		ViewChannel:             Permission(1 << 10),
		SendMessages:            Permission(1 << 11),
		SendTTSMessages:         Permission(1 << 12),
		ManageMessages:          Permission(1 << 13),
		EmbedLinks:              Permission(1 << 14),
		AttachFiles:             Permission(1 << 15),
		ReadMessageHistory:      Permission(1 << 16),
		MentionEveryone:         Permission(1 << 17),
		UseExternalEmojis:       Permission(1 << 18),
		ViewGuildInsights:       Permission(1 << 19),
		Connect:                 Permission(1 << 20),
		Speak:                   Permission(1 << 21),
		MuteMembers:             Permission(1 << 22),
		DeafenMembers:           Permission(1 << 23),
		MoveMembers:             Permission(1 << 24),
		UseVAD:                  Permission(1 << 25),
		ChangeNickname:          Permission(1 << 26),
		ManageNicknames:         Permission(1 << 27),
		ManageRoles:             Permission(1 << 28),
		ManageWebhooks:          Permission(1 << 29),
		ManageEmojisAndStickers: Permission(1 << 30),
		UseApplicationCommands:  Permission(1 << 31),
		RequestToSpeak:          Permission(1 << 32),
		ManageEvents:            Permission(1 << 33),
		ManageThreads:           Permission(1 << 34),
		CreatePublicThreads:     Permission(1 << 35),
		CreatePrivateThreads:    Permission(1 << 36),
		UseExternalStickers:     Permission(1 << 37),
		SendMessageInThreads:    Permission(1 << 38),
		UseEmbeddedActivity:     Permission(1 << 39),
		ModerateMembers:         Permission(1 << 40),
	}
}

var Permissions = perms()
