package stock

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"stockbot/internal/contracts"
)

// Handler adapts Service to the bot command interface.
type Handler struct {
	Service *Service
}

func NewHandler(urlTemplate string, httpClient *http.Client) *Handler {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 8 * time.Second}
	}
	return &Handler{Service: NewService(httpClient, urlTemplate)}
}

func (h *Handler) Handle(req contracts.BotRequest) (contracts.BotResponseSubmit, bool) {
	cmdBase := strings.ToLower(strings.TrimSpace(req.Command))
	if i := strings.Index(cmdBase, "="); i >= 0 {
		cmdBase = strings.TrimSpace(cmdBase[:i])
	}
	if cmdBase != "stock" {
		return contracts.BotResponseSubmit{}, false
	}
	symbol := parseSymbol(req.Command, req.Args)
	if symbol == "" {
		return contracts.BotResponseSubmit{RoomID: req.RoomID, Text: "Usage: /stock <symbol> or /stock=<symbol> (e.g., /stock aapl.us)"}, true
	}
	quote, err := h.Service.Fetch(symbol)
	if err != nil {
		return contracts.BotResponseSubmit{RoomID: req.RoomID, Text: fmt.Sprintf("Could not get quote for %s", symbol)}, true
	}
	closeVal := strings.TrimSpace(quote.Close)
	if closeVal == "" || closeVal == "0" || strings.EqualFold(closeVal, "N/D") {
		return contracts.BotResponseSubmit{RoomID: req.RoomID, Text: fmt.Sprintf("No quote for symbol %s", symbol)}, true
	}
	text := fmt.Sprintf("%s quote is %s per share", strings.ToUpper(symbol), closeVal)
	return contracts.BotResponseSubmit{RoomID: req.RoomID, Text: text}, true
}

func parseSymbol(command, args string) string {
	s := strings.TrimSpace(args)
	if s == "" {
		if i := strings.Index(command, "="); i >= 0 && i+1 < len(command) {
			return strings.TrimSpace(command[i+1:])
		}
		if i := strings.Index(args, "="); i >= 0 && i+1 < len(args) {
			return strings.TrimSpace(args[i+1:])
		}
	}
	return s
}
