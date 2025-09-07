package help

import (
	"testing"

	"stockbot/internal/contracts"
)

func TestHandle_Help(t *testing.T) {
	out, ok := Handle(contracts.BotRequest{Command: "help", RoomID: "r1"})
	if !ok {
		t.Fatal("expected handled")
	}
	if out.Text == "" {
		t.Fatal("expected help text")
	}
}
