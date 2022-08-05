package disgo

// Attachment represents a base Discord attachment
type Attachment struct {
	ID          string `json:"id"`
	Filename    string `json:"filename"`
	Description string `json:"description"`
	ContentType string `json:"content_type"`
	Size        int    `json:"size"`
	URL         string `json:"url"`
	ProxyURL    string `json:"proxy_url"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	Ephemeral   bool   `json:"ephemeral"`
}

// PartialAttachment represents a partial Discord attachment
type PartialAttachment struct {
	Id          string `json:"id"`
	Filename    string `json:"filename"`
	Description string `json:"description"`
}
