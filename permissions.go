package discord

type Permission int

const (
	zero        = 1 << iota
	one         = 1 << iota
	two         = 1 << iota
	three       = 1 << iota
	four        = 1 << iota
	five        = 1 << iota
	six         = 1 << iota
	seven       = 1 << iota
	eight       = 1 << iota
	nine        = 1 << iota
	ten         = 1 << iota
	eleven      = 1 << iota
	twelve      = 1 << iota
	thirteen    = 1 << iota
	fourteen    = 1 << iota
	fifteen     = 1 << iota
	sixteen     = 1 << iota
	seventeen   = 1 << iota
	eighteen    = 1 << iota
	nineteen    = 1 << iota
	twenty      = 1 << iota
	twentyone   = 1 << iota
	twentytwo   = 1 << iota
	twentythree = 1 << iota
	twentyfour  = 1 << iota
	twentyfive  = 1 << iota
	twentysix   = 1 << iota
	twentyseven = 1 << iota
	twentyeight = 1 << iota
	twentynine  = 1 << iota
	thirty      = 1 << iota
	thirtyone   = 1 << iota
	thirtytwo   = 1 << iota
	thirtythree = 1 << iota
	thirtyfour  = 1 << iota
	thirtyfive  = 1 << iota
	thirtysix   = 1 << iota
	thirtyseven = 1 << iota
	thirtyeight = 1 << iota
	thirtynine  = 1 << iota
	forty       = 1 << iota
)

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

func (p *permissions) Check(total Permission, other Permission) bool {
	return int(total)&int(other) == int(other)
}

func Permissions() *permissions {
	return &permissions{
		CreateInstantInvite:     Permission(zero),
		KickMembers:             Permission(one),
		BanMembers:              Permission(two),
		Administrator:           Permission(three),
		ManageChannels:          Permission(four),
		ManageGuild:             Permission(five),
		AddReactions:            Permission(six),
		ViewAuditLog:            Permission(seven),
		PrioritySpeaker:         Permission(eight),
		Stream:                  Permission(nine),
		ViewChannel:             Permission(ten),
		SendMessages:            Permission(eleven),
		SendTTSMessages:         Permission(twelve),
		ManageMessages:          Permission(thirteen),
		EmbedLinks:              Permission(fourteen),
		AttachFiles:             Permission(fifteen),
		ReadMessageHistory:      Permission(sixteen),
		MentionEveryone:         Permission(seventeen),
		UseExternalEmojis:       Permission(eighteen),
		ViewGuildInsights:       Permission(nineteen),
		Connect:                 Permission(twenty),
		Speak:                   Permission(twentyone),
		MuteMembers:             Permission(twentytwo),
		DeafenMembers:           Permission(twentythree),
		MoveMembers:             Permission(twentyfour),
		UseVAD:                  Permission(twentyfive),
		ChangeNickname:          Permission(twentysix),
		ManageNicknames:         Permission(twentyseven),
		ManageRoles:             Permission(twentyeight),
		ManageWebhooks:          Permission(twentynine),
		ManageEmojisAndStickers: Permission(thirty),
		UseApplicationCommands:  Permission(thirtyone),
		RequestToSpeak:          Permission(thirtytwo),
		ManageEvents:            Permission(thirtythree),
		ManageThreads:           Permission(thirtyfour),
		CreatePublicThreads:     Permission(thirtyfive),
		CreatePrivateThreads:    Permission(thirtysix),
		UseExternalStickers:     Permission(thirtyseven),
		SendMessageInThreads:    Permission(thirtyeight),
		UseEmbeddedActivity:     Permission(thirtynine),
		ModerateMembers:         Permission(forty),
	}
}

func (p *permissions) check(all Permission, permissions ...Permission) bool {
	res := true
	for _, p := range permissions {
		res = res && (int(all)&int(p) == int(p))
	}
	return res
}
