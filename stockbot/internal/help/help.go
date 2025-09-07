package help

import (
	"strings"

	"stockbot/internal/contracts"
)

func Handle(req contracts.BotRequest) (contracts.BotResponseSubmit, bool) {
	if strings.ToLower(strings.TrimSpace(req.Command)) != "help" {
		return contracts.BotResponseSubmit{}, false
	}
	text := strings.Join([]string{
		"Available commands:",
		"- /echo <text>",
		"- /stock <symbol> or /stock=<symbol> (e.g., aapl.us)",
		"- /help",
	}, "\n")
	return contracts.BotResponseSubmit{RoomID: req.RoomID, Text: text}, true
}
