package discord

type (
	commandTypes struct {
		Slash   int
		User    int
		Message int
	}

	channelTypes struct {
		Text          int
		DM            int
		Voice         int
		GroupDM       int
		Category      int
		News          int
		NewsThread    int
		PublicThread  int
		PrivateThread int
		StageVoice    int
		Directory     int
		Forum         int
	}

	buttonStyles struct {
		Blue  int
		Green int
		Red   int
		Grey  int
		Link  int
	}

	optionKind struct {
		String      int
		Integer     int
		Boolean     int
		User        int
		Channel     int
		Role        int
		Mentionable int
		Number      int
		Attachment  int
	}
)

var (
	OptionTypes = optionKind{
		String:      3,
		Integer:     4,
		Boolean:     5,
		User:        6,
		Channel:     7,
		Role:        8,
		Mentionable: 9,
		Number:      10,
		Attachment:  11,
	}

	ChannelTypes = channelTypes{
		Text:          0,
		DM:            1,
		Voice:         2,
		GroupDM:       3,
		Category:      4,
		News:          5,
		NewsThread:    10,
		PublicThread:  11,
		PrivateThread: 12,
		StageVoice:    13,
		Directory:     14,
		Forum:         15,
	}

	CommandTypes = commandTypes{
		Slash:   1,
		User:    2,
		Message: 3,
	}

	ButtonStyles = buttonStyles{
		Blue:  1,
		Grey:  2,
		Green: 3,
		Red:   4,
		Link:  5,
	}
)
