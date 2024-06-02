package container

import (
	"github.com/Inmovilizame/invoiceling/internal/repository"
	"github.com/Inmovilizame/invoiceling/pkg/render"
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

func NewDocumentService(renderType string, draft bool) (*service.Document, error) {
	repo := repository.CfgRepo{}

	doc, err := service.NewDocumentService(repo.GetDebug(), draft, repo.GetPdfOutputDir())
	if err != nil {
		return nil, err
	}

	switch renderType {
	case "Basic":
		r, err := render.NewPdfBasicRender()
		if err != nil {
			return nil, err
		}

		doc.SetRenderer(r)
	}

	return doc, nil
}
