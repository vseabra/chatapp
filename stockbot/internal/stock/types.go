package stock

// QuoteRow represents one parsed CSV row.
type QuoteRow struct {
	Symbol  string
	DateISO string
	TimeISO string
	Open    string
	High    string
	Low     string
	Close   string
	Volume  string
}
