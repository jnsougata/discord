package disgo

import "encoding/json"

// Channel represents a Discord channel of any type
type Channel struct {
	Id                         string        `json:"id"`
	Type                       int           `json:"type"`
	GuildId                    string        `json:"guild_id"`
	Position                   int           `json:"position"`
	Overwrites                 []interface{} `json:"permission_overwrites"`
	Name                       string        `json:"name"`
	Topic                      string        `json:"topic"`
	NSFW                       bool          `json:"nsfw"`
	LastMessageId              string        `json:"last_message_id"`
	Bitrate                    int           `json:"bitrate"`
	UserLimit                  int           `json:"user_limit"`
	RateLimitPerUser           int           `json:"rate_limit_per_user"`
	Recipients                 []interface{} `json:"recipients"`
	Icon                       string        `json:"icon"`
	OwnerId                    string        `json:"owner_id"`
	ApplicationId              string        `json:"application_id"`
	ParentId                   string        `json:"parent_id"`
	LastPinTime                int           `json:"last_pin_timestamp"`
	RTCRegion                  string        `json:"rtc_region"`
	VideoQualityMode           int           `json:"video_quality_mode"`
	MessageCount               int           `json:"message_count"`
	ThreadMetaData             interface{}   `json:"thread_metadata"`
	Member                     interface{}   `json:"member"`
	DefaultAutoArchiveDuration int           `json:"default_auto_archive_days"`
	Permissions                string        `json:"permissions"`
	Flags                      int           `json:"flags"`
	TotalMessages              int           `json:"total_messages"`
}

func UnmarshalChannel(payload interface{}) *Channel {
	ch := &Channel{}
	data, _ := json.Marshal(payload)
	_ = json.Unmarshal(data, ch)
	return ch
}
