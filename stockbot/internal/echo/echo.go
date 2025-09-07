package echo

import "strings"

// BotRequested models the incoming payload for bot.requested
type BotRequested struct {
	Command        string `json:"command"`
	Args           string `json:"args"`
	RoomID         string `json:"roomId"`
	RequestUserID  string `json:"requestUserId"`
	RequestedAtISO string `json:"requestedAt"`
}

// BotResponseSubmit models the outgoing payload for bot.response.submit
type BotResponseSubmit struct {
	RoomID string `json:"roomId"`
	Text   string `json:"text"`
}

// Handle returns a response if the command is echo; otherwise empty, false.
func Handle(req BotRequested) (BotResponseSubmit, bool) {
	if strings.ToLower(strings.TrimSpace(req.Command)) != "echo" {
		return BotResponseSubmit{}, false
	}
	text := strings.TrimSpace(req.Args)
	return BotResponseSubmit{
		RoomID: req.RoomID,
		Text:   text,
	}, true
}
