package model

const (
	DefaultDueSpan = 30
)

type Invoice struct {
	ID     string `json:"id" yaml:"id"`
	Status string `json:"status" yaml:"status"`
	Logo   string `json:"logo" yaml:"logo"`

	From *Freelancer `json:"from" yaml:"from"`
	To   *Client     `json:"to" yaml:"to"`

	Date string `json:"date" yaml:"date"`
	Due  string `json:"due" yaml:"due"`

	Items []Item `json:"items" yaml:"items"`

	Tax       float64 `json:"tax" yaml:"tax"`
	Retention float64 `json:"retention" yaml:"retention"`
	Discount  float64 `json:"discount" yaml:"discount"`
	Currency  string  `json:"currency" yaml:"currency"`

	Payment *Payment `json:"payment" yaml:"payment"`

	Note []string `json:"note" yaml:"note"`
}

type Item struct {
	Description string  `json:"description" yaml:"description"`
	Quantity    int     `json:"quantity" yaml:"quantity"`
	Rate        float64 `json:"rate" yaml:"rate"`
}

type Payment struct {
	Holder string `json:"holder" yaml:"holder"`
	Iban   string `json:"iban" yaml:"iban"`
	Swift  string `json:"swift" yaml:"swift"`
}
