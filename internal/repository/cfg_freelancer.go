package repository

import (
	"github.com/Inmovilizame/invoiceling/pkg/model"
	"github.com/spf13/viper"
)

type CfgFreelancer struct{}

func (c CfgFreelancer) Get() *model.Freelancer {
	return &model.Freelancer{
		Company:  viper.GetString("freelancer.company"),
		Name:     viper.GetString("freelancer.name"),
		Email:    viper.GetString("freelancer.email"),
		Phone:    viper.GetString("freelancer.phone"),
		VatID:    viper.GetString("freelancer.vat_id"),
		Address1: viper.GetString("freelancer.address1"),
		Address2: viper.GetString("freelancer.address2"),
	}
}
