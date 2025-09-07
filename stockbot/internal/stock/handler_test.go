package stock

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"stockbot/internal/contracts"
)

func TestHandler_Handle_WithArg(t *testing.T) {
	body := "Symbol,Date,Time,Open,High,Low,Close,Volume\nAAPL.US,2025-09-05,22:00:17,239.995,241.32,238.4901,239.69,54870397\n"
	cli := newTestClient(func(r *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewBufferString(body))}, nil
	})
	h := NewHandler("http://test/?s={{symbol}}", cli)
	out, ok := h.Handle(contracts.BotRequest{Command: "stock", Args: "aapl.us", RoomID: "r1"})
	if !ok {
		t.Fatal("expected handled")
	}
	if out.Text == "" {
		t.Fatal("expected non-empty text")
	}
}

func TestHandler_Handle_EqualsSyntax(t *testing.T) {
	body := "Symbol,Date,Time,Open,High,Low,Close,Volume\nAAPL.US,2025-09-05,22:00:17,239.995,241.32,238.4901,239.69,54870397\n"
	cli := newTestClient(func(r *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewBufferString(body))}, nil
	})
	h := NewHandler("http://test/?s={{symbol}}", cli)
	out, ok := h.Handle(contracts.BotRequest{Command: "stock=aapl.us", Args: "", RoomID: "r1"})
	if !ok {
		t.Fatal("expected handled")
	}
	if out.Text == "" {
		t.Fatal("expected non-empty text")
	}
}

func TestHandler_Handle_ND(t *testing.T) {
	body := "Symbol,Date,Time,Open,High,Low,Close,Volume\nAAPL.US,2025-09-05,22:00:17,239.995,241.32,238.4901,N/D,54870397\n"
	cli := newTestClient(func(r *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewBufferString(body))}, nil
	})
	h := NewHandler("http://test/?s={{symbol}}", cli)
	out, ok := h.Handle(contracts.BotRequest{Command: "stock", Args: "aapl.us", RoomID: "r1"})
	if !ok {
		t.Fatal("expected handled")
	}
	if out.Text == "" || out.Text == "AAPL.US quote is  per share" {
		t.Fatal("expected friendly no-quote message")
	}
}
