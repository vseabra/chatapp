package echo

import (
	"strings"

	"stockbot/internal/contracts"
)

// Handle returns a response for the echo command.
func Handle(req contracts.BotRequest) (contracts.BotResponseSubmit, bool) {
	if strings.ToLower(strings.TrimSpace(req.Command)) != "echo" {
		return contracts.BotResponseSubmit{}, false
	}
	text := strings.TrimSpace(req.Args)
	return contracts.BotResponseSubmit{
		RoomID: req.RoomID,
		Text:   text,
	}, true
}
