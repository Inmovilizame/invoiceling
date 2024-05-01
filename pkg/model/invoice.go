package model

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"os"
	"path/filepath"
	"time"
)

type DateFormat string

const (
	DF_Text DateFormat = "Jan 02, 2006"
	DF_YMD  DateFormat = "2006-01-02"
	DF_DMY  DateFormat = "02/01/2006"
)

var roMask = os.FileMode(0400)

type Invoice struct {
	ID   string `json:"id" yaml:"id"`
	Logo string `json:"logo" yaml:"logo"`

	From *Freelancer `json:"from" yaml:"from"`
	To   *Client     `json:"to" yaml:"to"`

	Date string `json:"date" yaml:"date"`
	Due  string `json:"due" yaml:"due"`

	Items []Item `json:"items" yaml:"items"`

	Tax      float64 `json:"tax" yaml:"tax"`
	Discount float64 `json:"discount" yaml:"discount"`
	Currency string  `json:"currency" yaml:"currency"`

	Note string `json:"note" yaml:"note"`
}

func NewInvoice(me *Freelancer, to *Client, dueDays int, dateFormat DateFormat) *Invoice {
	invoiceDate := time.Now()
	dueDate := invoiceDate.AddDate(0, 0, dueDays)

	items := NewItems(5)

	return &Invoice{
		ID:       "F24-001",
		Logo:     "logo.png",
		From:     me,
		To:       to,
		Date:     invoiceDate.Format(string(dateFormat)),
		Due:      dueDate.Format(string(dateFormat)),
		Items:    items,
		Tax:      0,
		Discount: 0,
		Currency: "EUR",
	}
}

type Item struct {
	Description string  `json:"description" yaml:"description"`
	Quantity    int     `json:"quantity" yaml:"quantity"`
	Rate        float64 `json:"rate" yaml:"rate"`
}

func NewItems(amount int) []Item {
	items := make([]Item, amount)
	for i := 0; i < amount; i++ {
		items[i] = Item{
			Description: fmt.Sprintf("Product %d", i),
			Quantity:    rand.IntN(10),
			Rate:        rand.Float64() * 100.0,
		}
	}

	return items
}

func (i *Invoice) Save(basePath string) error {
	jsonBytes, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		return err
	}

	invoicePath := filepath.Join(basePath, i.ID+".json")

	err = os.WriteFile(invoicePath, jsonBytes, roMask)
	if err != nil {
		return err
	}

	return nil
}
