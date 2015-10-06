package smith

type Event struct {
	Type string `json:"type"`
}

type MessageEvent struct {
	Type      string `json:"type"`
	Subtype   string `json:"subtype,omitempty"`
	Channel   string `json:"channel"`
	UserId    string `json:"user"`
	UserName  string `json:"username"`
	Text      string `json:"text"`
	Timestamp string `json:"ts"`
	Edited    struct {
		User      string `json:"user"`
		Timestamp string `json:"timestamp"`
	} `json:"edited"`
	User *User `json:"-"`
}

type ChannelCreatedEvent struct {
	Type    string   `json:"type"`
	Channel struct{} `json:"channel"`
}
