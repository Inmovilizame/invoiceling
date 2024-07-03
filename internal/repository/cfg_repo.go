package repository

import (
	"github.com/Inmovilizame/invoiceling/pkg/model"
	"github.com/spf13/viper"
)

type CfgRepo struct{}

func (c CfgRepo) GetDebug() bool {
	return viper.GetBool("debug")
}

func (c CfgRepo) GetNotes() map[string]string {
	return map[string]string{
		"no_due":          viper.GetString("notes.no_due"),
		"retention_not_0": viper.GetString("notes.retention_not_0"),
		"vat_0":           viper.GetString("notes.vat_0"),
	}
}

func (c CfgRepo) GetPdfOutputDir() string {
	return viper.GetString("dirs.pdf")
}

func (c CfgRepo) GetCurrency() string {
	return viper.GetString("invoice.currency")
}

func (c CfgRepo) GetIDFormat() string {
	return viper.GetString("invoice.id_format")
}

func (c CfgRepo) GetLogo() string {
	return viper.GetString("invoice.logo")
}

func (c CfgRepo) GetFreelancer() model.Freelancer {
	return model.Freelancer{
		Company:  viper.GetString("freelancer.company"),
		Name:     viper.GetString("freelancer.name"),
		Email:    viper.GetString("freelancer.email"),
		Phone:    viper.GetString("freelancer.phone"),
		VatID:    viper.GetString("freelancer.vat_id"),
		Address1: viper.GetString("freelancer.address1"),
		Address2: viper.GetString("freelancer.address2"),
	}
}

func (c CfgRepo) GetPaymentInfo() model.Payment {
	return model.Payment{
		Holder: viper.GetString("payment.holder"),
		Iban:   viper.GetString("payment.iban"),
		Swift:  viper.GetString("payment.swift"),
	}
}
