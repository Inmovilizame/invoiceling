package container

import (
	"github.com/Inmovilizame/invoiceling/assets"
	"github.com/Inmovilizame/invoiceling/internal/repository"
	"github.com/Inmovilizame/invoiceling/pkg/service"
	"github.com/spf13/viper"
)

func NewInvoiceService() *service.InvoiceService {
	invoiceRepo := repository.NewFsInvoice(
		viper.GetString("dirs.invoice"),
	)

	clientRepo := repository.NewFsClient(
		viper.GetString("dirs.client"),
	)

	return service.NewInvoiceService(
		invoiceRepo,
		clientRepo,
		repository.CfgRepo{},
	)
}

func NewClientService() *service.Client {
	clientRepo := repository.NewFsClient(
		viper.GetString("dirs.client"),
	)

	return service.NewClientService(clientRepo)
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
