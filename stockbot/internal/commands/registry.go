package commands

import (
	"strings"

	"stockbot/internal/contracts"
)

// similar to the http handler type from stdlib
type Handler func(contracts.BotRequest) (contracts.BotResponseSubmit, bool)

type Registry struct {
	commandToHandler map[string]Handler
}

func NewRegistry() *Registry {
	return &Registry{commandToHandler: make(map[string]Handler)}
}

func (r *Registry) Register(command string, handler Handler) {
	key := strings.TrimPrefix(strings.ToLower(strings.TrimSpace(command)), "/")
	r.commandToHandler[key] = handler
}

// Dispatch routes a request to the appropriate handler by command.
func (r *Registry) Dispatch(req contracts.BotRequest) (contracts.BotResponseSubmit, bool) {
	key := strings.TrimPrefix(strings.ToLower(strings.TrimSpace(req.Command)), "/")
	if h, ok := r.commandToHandler[key]; ok {
		return h(req)
	}
	return contracts.BotResponseSubmit{}, false
}
