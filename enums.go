package discord

type commandKind int
type commandKinds struct {
	Slash   commandKind
	User    commandKind
	Message commandKind
}

type ChannelKind int
type channelKinds struct {
	Text          ChannelKind
	DM            ChannelKind
	Voice         ChannelKind
	GroupDM       ChannelKind
	Category      ChannelKind
	News          ChannelKind
	NewsThread    ChannelKind
	PublicThread  ChannelKind
	PrivateThread ChannelKind
	StageVoice    ChannelKind
	Directory     ChannelKind
	Forum         ChannelKind
}

type buttonStyle int
type buttonStyles struct {
	Blue  buttonStyle
	Green buttonStyle
	Red   buttonStyle
	Grey  buttonStyle
	Link  buttonStyle
}

type optionKind int

const (
	stringOption      optionKind = 3
	integerOption     optionKind = 4
	booleanOption     optionKind = 5
	userOption        optionKind = 6
	channelOption     optionKind = 7
	roleOption        optionKind = 8
	mentionableOption optionKind = 9
	numberOption      optionKind = 10
	attachmentOption  optionKind = 11
)

// exported variables

var CommandKinds = commandKinds{
	Slash:   commandKind(1),
	User:    commandKind(2),
	Message: commandKind(3),
}

var ChannelKinds = channelKinds{
	Text:          ChannelKind(0),
	DM:            ChannelKind(1),
	Voice:         ChannelKind(2),
	GroupDM:       ChannelKind(3),
	Category:      ChannelKind(4),
	News:          ChannelKind(5),
	NewsThread:    ChannelKind(10),
	PublicThread:  ChannelKind(11),
	PrivateThread: ChannelKind(12),
	StageVoice:    ChannelKind(13),
	Directory:     ChannelKind(14),
	Forum:         ChannelKind(15),
}

var ButtonStyles = buttonStyles{
	Blue:  buttonStyle(1),
	Grey:  buttonStyle(2),
	Green: buttonStyle(3),
	Red:   buttonStyle(4),
	Link:  buttonStyle(5),
}
