package model

import (
	"time"
)

const (
	DefaultDueSpan = 30
)

type Item struct {
	Description string  `json:"description" yaml:"description"`
	Quantity    int     `json:"quantity" yaml:"quantity"`
	Vat         float64 `json:"vat" yaml:"vat"`
	Rate        float64 `json:"rate" yaml:"rate"`
}

func (i *Item) GetAmount() float64 {
	return float64(i.Quantity) * i.Rate
}

func (i *Item) GetVat() float64 {
	return i.GetAmount() * i.Vat / 100 //nolint:mnd //static percentage calculation
}

type Notes struct {
	Default       string `json:"default" yaml:"default"`
	RetentionNot0 string `json:"retentionNot0" yaml:"retentionNot0"`
	Vat0          string `json:"vat0" yaml:"vat0"`
}

func (n Notes) ToSlice() []string {
	var notes []string

	notes = append(notes, n.Default)

	if n.Vat0 != "" {
		notes = append(notes, n.Vat0)
	}

	if n.RetentionNot0 != "" {
		notes = append(notes, n.RetentionNot0)
	}

	return notes
}

type Payment struct {
	Holder string `json:"holder" yaml:"holder"`
	Iban   string `json:"iban" yaml:"iban"`
	Swift  string `json:"swift" yaml:"swift"`
}

type TaxInfo struct {
	Vat       float64 `json:"vat" yaml:"vat"`
	Retention float64 `json:"retention" yaml:"retention"`
}

func (t TaxInfo) VatAmount(amount float64) float64 {
	return amount * t.Vat / 100 //nolint:mnd //calculating percentage
}

func (t TaxInfo) RetentionAmount(amount float64) float64 {
	return amount * t.Retention / 100 //nolint:mnd //calculating percentage
}

type Invoice struct {
	ID     string `json:"id" yaml:"id"`
	Status string `json:"status" yaml:"status"`
	Logo   string `json:"logo" yaml:"logo"`

	From Freelancer `json:"from" yaml:"from"`
	To   Client     `json:"to" yaml:"to"`

	Date time.Time     `json:"date" yaml:"date"`
	Due  time.Duration `json:"due" yaml:"due"`

	Items []*Item `json:"items" yaml:"items"`

	Tax      TaxInfo `json:"tax" yaml:"tax"`
	Discount float64 `json:"discount" yaml:"discount"`
	Currency string  `json:"currency" yaml:"currency"`

	Payment Payment `json:"payment" yaml:"payment"`

	Notes Notes `json:"notes" yaml:"notes"`
}

func NewInvoice(id string, due time.Duration, currency, note, noDueNote string) *Invoice {
	notes := Notes{Default: note}
	if due == 0 {
		notes.Default += " " + noDueNote
	}

	return &Invoice{
		ID:       id,
		Status:   "CREATED",
		Date:     time.Now(),
		Due:      due,
		Items:    []*Item{},
		Tax:      TaxInfo{},
		Discount: 0,
		Currency: currency,
		Notes:    notes,
	}
}

func (i *Invoice) SetTaxes(vat, retention float64, configNotes map[string]string) {
	i.Tax = TaxInfo{
		Vat:       vat,
		Retention: retention,
	}

	i.Tax.Vat = vat
	if vat == 0 {
		i.Notes.Vat0 = configNotes["vat_0"]
	}

	i.Tax.Retention = retention
	if retention != 0 {
		i.Notes.RetentionNot0 = configNotes["retention_not_0"]
	}
}

func (i *Invoice) AddItem(item Item) {
	if item.Vat == 0 {
		item.Vat = i.Tax.Vat
	}

	i.Items = append(i.Items, &item)
}
