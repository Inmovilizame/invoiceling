package model

type Freelancer struct {
	Company  string `json:"company" yaml:"company"`
	Name     string `json:"name" yaml:"name"`
	Email    string `json:"email" yaml:"email"`
	Phone    string `json:"phone" yaml:"phone"`
	VatID    string `json:"vat_id" yaml:"vat_id"`
	Address1 string `json:"address1" yaml:"address1"`
	Address2 string `json:"address2" yaml:"address2"`
}
