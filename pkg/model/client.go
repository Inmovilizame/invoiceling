package model

type Client struct {
	ID       string `json:"id" yaml:"id"`
	Name     string `json:"name" yaml:"name"`
	VatID    string `json:"vat_id" yaml:"vat_id"`
	Address1 string `json:"address1" yaml:"address1"`
	Address2 string `json:"address2" yaml:"address2"`
}
