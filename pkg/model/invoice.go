package model

const (
	DefaultDueSpan = 30
)

type Invoice struct {
	ID     string `json:"id" yaml:"id"`
	Status string `json:"status" yaml:"status"`
	Logo   string `json:"logo" yaml:"logo"`

	From Freelancer `json:"from" yaml:"from"`
	To   Client     `json:"to" yaml:"to"`

	Date string `json:"date" yaml:"date"`
	Due  string `json:"due" yaml:"due"`

	Items []Item `json:"items" yaml:"items"`

	Tax      TaxInfo `json:"tax" yaml:"tax"`
	Discount float64 `json:"discount" yaml:"discount"`
	Currency string  `json:"currency" yaml:"currency"`

	Payment Payment `json:"payment" yaml:"payment"`

	Notes Notes `json:"notes" yaml:"notes"`
}

type TaxInfo struct {
	Vat       float64 `json:"vat" yaml:"vat"`
	Retention float64 `json:"retention" yaml:"retention"`
}

func (t TaxInfo) GetVat(amount float64) float64 {
	return amount * t.Vat / 100 //nolint:gomnd //calculating percentage
}

func (t TaxInfo) GetRetention(amount float64) float64 {
	return amount * t.Retention / 100 //nolint:gomnd //calculating percentage
}

type Item struct {
	Description string  `json:"description" yaml:"description"`
	Quantity    int     `json:"quantity" yaml:"quantity"`
	Rate        float64 `json:"rate" yaml:"rate"`
}

func (i Item) GetAmount() float64 {
	return float64(i.Quantity) * i.Rate
}

type Payment struct {
	Holder string `json:"holder" yaml:"holder"`
	Iban   string `json:"iban" yaml:"iban"`
	Swift  string `json:"swift" yaml:"swift"`
}

type Notes struct {
	Default       string `json:"default" yaml:"default"`
	Vat0          string `json:"vat0" yaml:"vat0"`
	RetentionNot0 string `json:"retentionNot0" yaml:"retentionNot0"`
}

func (n Notes) ToSlice() []string {
	notes := []string{}

	notes = append(notes, n.Default)

	if n.Vat0 != "" {
		notes = append(notes, n.Vat0)
	}

	if n.RetentionNot0 != "" {
		notes = append(notes, n.RetentionNot0)
	}

	return notes
}
