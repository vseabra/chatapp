package echo

import (
	"testing"

	"stockbot/internal/contracts"
)

func TestHandle_Echo(t *testing.T) {
	out, ok := Handle(contracts.BotRequest{Command: "echo", Args: "hello", RoomID: "r1"})
	if !ok {
		t.Fatal("expected handled")
	}
	if out.Text != "hello" {
		t.Fatalf("expected echo 'hello', got %q", out.Text)
	}
}
