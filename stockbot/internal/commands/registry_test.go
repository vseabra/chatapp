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
	reg.Register("stock", func(req contracts.BotRequest) (contracts.BotResponseSubmit, bool) {
		return contracts.BotResponseSubmit{RoomID: req.RoomID, Text: "stock quote"}, true
	})

	// Test basic command dispatch
	out, ok := reg.Dispatch(contracts.BotRequest{Command: "echo", Args: "x", RoomID: "r1"})
	if !ok || out.Text != "x" {
		t.Fatal("expected dispatch to echo handler")
	}

	// Test command with equals sign (like /stock=aapl.usxd)
	out, ok = reg.Dispatch(contracts.BotRequest{Command: "stock=aapl.usxd", Args: "", RoomID: "r1"})
	if !ok || out.Text != "stock quote" {
		t.Fatal("expected dispatch to stock handler for command with equals sign")
	}

	// Test command with equals sign and spaces
	out, ok = reg.Dispatch(contracts.BotRequest{Command: " stock = aapl.usxd ", Args: "", RoomID: "r1"})
	if !ok || out.Text != "stock quote" {
		t.Fatal("expected dispatch to stock handler for command with equals sign and spaces")
	}
}
