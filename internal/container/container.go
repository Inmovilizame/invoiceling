package container

import (
	"github.com/Inmovilizame/invoiceling/internal/repository"
	"github.com/Inmovilizame/invoiceling/pkg/i18n"
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

func NewDocumentService(renderType string, draft bool, language i18n.Language) (*service.Document, error) {
	repo := repository.CfgRepo{}

	doc, err := service.NewDocumentService(repo.GetDebug(), draft, repo.GetPdfOutputDir())
	if err != nil {
		return nil, err
	}

	translator := i18n.NewTranslator()
	translator.SetLanguage(language)

	switch renderType {
	case "Basic":
		r, err := render.NewPdfBasicRender(translator)
		if err != nil {
			return nil, err
		}

		doc.SetRenderer(r)
	default:
		r, err := render.NewPdfBasicRender(translator)
		if err != nil {
			return nil, err
		}

		doc.SetRenderer(r)
	}

	return doc, nil
}
