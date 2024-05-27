package container

import (
	"github.com/Inmovilizame/invoiceling/assets"
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
		invoiceRepo,
		clientRepo,
		repository.CfgRepo{},
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
		repository.CfgRepo{},
	)
}

func NewDocumentService() (*service.Document, error) {
	repo := repository.CfgRepo{}

	interFont, err := assets.FS.ReadFile("fonts/Inter.ttf")
	if err != nil {
		return nil, err
	}

	interBoldFont, err := assets.FS.ReadFile("fonts/Inter-Bold.ttf")
	if err != nil {
		return nil, err
	}

	fontMap := map[string][]byte{
		"Inter":      interFont,
		"Inter-Bold": interBoldFont,
	}

	doc, err := service.NewInvoiceRender(fontMap, repo.GetPdfOutputDir(), false)
	if err != nil {
		return nil, err
	}

	return doc, nil
}
