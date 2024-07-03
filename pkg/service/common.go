package service

import (
	"github.com/Inmovilizame/invoiceling/internal/repository"
	"github.com/Inmovilizame/invoiceling/pkg/model"
)

type ClientRepo interface {
	List(filter repository.Filter[*model.Client]) []*model.Client
	Create(client *model.Client) error
	Read(clientID string) *model.Client
	Update(invoice *model.Client) *model.Client
	Delete(clientID string) error
}

type InvoiceRepo interface {
	List(filter repository.Filter[*model.Invoice]) []*model.Invoice
	Create(invoice *model.Invoice) error
	Read(invoiceID string) *model.Invoice
	Update(invoice *model.Invoice) *model.Invoice
	Delete(invoiceID string) error
}

type CfgRepo interface {
	GetNotes() map[string]string
	GetPdfOutputDir() string
	GetCurrency() string
	GetIDFormat() string
	GetLogo() string
	GetFreelancer() model.Freelancer
	GetPaymentInfo() model.Payment
}

type RendererInterface interface {
	Render(invoice *model.Invoice, draft bool) error
	SaveTo(filename string) error
}
