package service

import "github.com/Inmovilizame/invoiceling/pkg/model"

type InvoiceService interface {
	List() []*model.Invoice
	Create(invoice *model.Invoice) error
	Read(invoiceID string) *model.Invoice
	Update(invoiceID string, invoice *model.Invoice) *model.Invoice
	Delete(invoiceID string) error
}
