package smith

type Msg struct {
	Type string `json:"type"`
}

type MessageMsg struct {
	Type      string `json:"type"`
	Channel   string `json:"channel"`
	User      string `json:"user"`
	Text      string `json:"text"`
	Timestamp string `json:"ts"`
	Team      string `json:"team"`
}

type ChatMsg struct {
	Id      int    `json:"id"`
	Type    string `json:"type"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
}
