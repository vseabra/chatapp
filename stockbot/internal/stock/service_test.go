package stock

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"
)

// roundTripFunc allows mocking http.Client transport
type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func newTestClient(fn roundTripFunc) *http.Client {
	return &http.Client{Transport: fn}
}

func TestService_Fetch_Success(t *testing.T) {
	body := "Symbol,Date,Time,Open,High,Low,Close,Volume\nAAPL.US,2025-09-05,22:00:17,239.995,241.32,238.4901,239.69,54870397\n"
	cli := newTestClient(func(r *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewBufferString(body))}, nil
	})
	svc := NewService(cli, "http://test/?s={{symbol}}")
	q, err := svc.Fetch("aapl.us")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if q.Close != "239.69" {
		t.Fatalf("expected close 239.69, got %s", q.Close)
	}
}

func TestService_Fetch_UpstreamError(t *testing.T) {
	cli := newTestClient(func(r *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 500, Body: io.NopCloser(bytes.NewBuffer(nil))}, nil
	})
	svc := NewService(cli, "http://test/?s={{symbol}}")
	_, err := svc.Fetch("aapl.us")
	if err == nil {
		t.Fatal("expected error for upstream status")
	}
}

func TestService_Fetch_NoData(t *testing.T) {
	body := "Symbol,Date,Time,Open,High,Low,Close,Volume\n"
	cli := newTestClient(func(r *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewBufferString(body))}, nil
	})
	svc := NewService(cli, "http://test/?s={{symbol}}")
	_, err := svc.Fetch("aapl.us")
	if err == nil {
		t.Fatal("expected error for no data")
	}
}

func TestService_Fetch_ShortRow(t *testing.T) {
	body := "Symbol,Date,Time,Open,High,Low,Close,Volume\nAAPL.US,2025-09-05\n"
	cli := newTestClient(func(r *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewBufferString(body))}, nil
	})
	svc := NewService(cli, "http://test/?s={{symbol}}")
	_, err := svc.Fetch("aapl.us")
	if err == nil {
		t.Fatal("expected error for short row")
	}
}

func TestService_Fetch_NetworkError(t *testing.T) {
	cli := newTestClient(func(r *http.Request) (*http.Response, error) {
		return nil, errors.New("dial")
	})
	svc := NewService(cli, "http://test/?s={{symbol}}")
	_, err := svc.Fetch("aapl.us")
	if err == nil {
		t.Fatal("expected network error")
	}
}
