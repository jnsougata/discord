package discord

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

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
	state                      *state
}

func (c *Channel) Send(draft Draft) (Message, error) {
	body, err := draft.marshal()
	path := fmt.Sprintf("/channels/%s/messages", c.Id)
	r := multipartReq("POST", path, body, c.state.Token, draft.Files...)
	bs, _ := io.ReadAll(r.fire().Body)
	var m Message
	_ = json.Unmarshal(bs, &m)
	m.state = c.state
	if draft.DeleteAfter > 0 {
		go func() {
			time.Sleep(time.Second * time.Duration(draft.DeleteAfter))
			m.Delete()
		}()
	}
	return m, err
}
