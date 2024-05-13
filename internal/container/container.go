package container

import (
	"github.com/Inmovilizame/invoiceling/internal/repository"
	"github.com/Inmovilizame/invoiceling/pkg/service"
	"github.com/spf13/viper"
)

func NewInvoiceService() *service.Invoice {
	invoiceRepo := repository.NewFsInvoice(
		viper.GetString("dirs.invoice"),
	)

	clientRepo := repository.NewFsClient(
		viper.GetString("dirs.client"),
	)

	return service.NewInvoice(
		viper.GetString("invoice.currency"),
		viper.GetString("invoice.id_format"),
		viper.GetString("invoice.logo"),
		invoiceRepo,
		clientRepo,
	)
}

func NewClientService() *service.Client {
	clientRepo := repository.NewFsClient(
		viper.GetString("dirs.client"),
	)

	return service.NewClient(clientRepo)

}

func NewFreelancerService() *service.Freelancer {
	return service.NewFreelancer(
		repository.CfgFreelancer{},
	)
}
