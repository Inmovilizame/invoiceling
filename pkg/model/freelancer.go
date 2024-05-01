package model

import "github.com/spf13/viper"

const notValidConfig = "Not Valid"

type Freelancer struct {
	Company  string `json:"company" yaml:"company"`
	Name     string `json:"name" yaml:"name"`
	Email    string `json:"email" yaml:"email"`
	Phone    string `json:"phone" yaml:"phone"`
	VatID    string `json:"vat_id" yaml:"vat_id"`
	Address1 string `json:"address1" yaml:"address1"`
	Address2 string `json:"address2" yaml:"address2"`
}

func LoadFreelancer() Freelancer {
	return Freelancer{
		Company:  viper.GetString("freelancer.company"),
		Name:     viper.GetString("freelancer.name"),
		Email:    viper.GetString("freelancer.email"),
		Phone:    viper.GetString("freelancer.phone"),
		VatID:    viper.GetString("freelancer.vat_id"),
		Address1: viper.GetString("freelancer.address1"),
		Address2: viper.GetString("freelancer.address2"),
	}
}
