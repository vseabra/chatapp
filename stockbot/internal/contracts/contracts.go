package contracts

// BotRequest models the incoming payload for bot.requested
type BotRequest struct {
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
