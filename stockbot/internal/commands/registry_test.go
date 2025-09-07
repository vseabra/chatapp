package commands

import (
	"testing"

	"stockbot/internal/contracts"
)

func TestRegistry_Dispatch(t *testing.T) {
	reg := NewRegistry()
	reg.Register("echo", func(req contracts.BotRequest) (contracts.BotResponseSubmit, bool) {
		return contracts.BotResponseSubmit{RoomID: req.RoomID, Text: req.Args}, true
	})
	out, ok := reg.Dispatch(contracts.BotRequest{Command: "echo", Args: "x", RoomID: "r1"})
	if !ok || out.Text != "x" {
		t.Fatal("expected dispatch to echo handler")
	}
}
