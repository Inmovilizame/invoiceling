package model

type Invoice struct {
	ID     string `json:"id" yaml:"id"`
	Status string `json:"status" yaml:"status"`
	Logo   string `json:"logo" yaml:"logo"`

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

type Item struct {
	Description string  `json:"description" yaml:"description"`
	Quantity    int     `json:"quantity" yaml:"quantity"`
	Rate        float64 `json:"rate" yaml:"rate"`
}
