package stock

import (
	"encoding/csv"
	"errors"
	"io"
	"net/http"
	"strings"
)

// Service fetches and parses quote data.
type Service struct {
	HTTP        *http.Client
	URLTemplate string
}

func NewService(httpClient *http.Client, urlTemplate string) *Service {
	return &Service{HTTP: httpClient, URLTemplate: urlTemplate}
}

// Fetch returns a parsed QuoteRow for the given symbol.
func (s *Service) Fetch(symbol string) (QuoteRow, error) {
	url := replaceSymbol(s.URLTemplate, symbol)
	httpClient := s.HTTP
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	resp, err := httpClient.Get(url)
	if err != nil {
		return QuoteRow{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return QuoteRow{}, errors.New("upstream status")
	}
	r := csv.NewReader(resp.Body)
	if _, err := r.Read(); err != nil { // header
		return QuoteRow{}, err
	}
	row, err := r.Read()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return QuoteRow{}, errors.New("no data")
		}
		return QuoteRow{}, err
	}
	if len(row) < 8 {
		return QuoteRow{}, errors.New("short row")
	}
	return QuoteRow{
		Symbol:  row[0],
		DateISO: row[1],
		TimeISO: row[2],
		Open:    row[3],
		High:    row[4],
		Low:     row[5],
		Close:   row[6],
		Volume:  row[7],
	}, nil
}

func replaceSymbol(tpl, symbol string) string {
	return strings.ReplaceAll(tpl, "{{symbol}}", strings.ToLower(symbol))
}
