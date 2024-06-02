package service

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Inmovilizame/invoiceling/internal/repository"

	"github.com/Inmovilizame/invoiceling/pkg/model"
)

type DateFormat string

type InvoiceRepo interface {
	List(filter repository.Filter[*model.Invoice]) []*model.Invoice
	Create(invoice *model.Invoice) error
	Read(invoiceID string) *model.Invoice
	Update(invoiceID string, invoice *model.Invoice) *model.Invoice
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

type InvoiceService struct {
	iRepo   InvoiceRepo
	cRepo   ClientRepo
	cfgRepo CfgRepo
}

func NewInvoiceService(iRepo InvoiceRepo, cRepo ClientRepo, cfgRepo CfgRepo) *InvoiceService {
	return &InvoiceService{
		iRepo:   iRepo,
		cRepo:   cRepo,
		cfgRepo: cfgRepo,
	}
}

func (is *InvoiceService) List(filter repository.Filter[*model.Invoice]) []*model.Invoice {
	return is.iRepo.List(filter)
}

func (is *InvoiceService) Create(
	id int,
	clientID string,
	dueDays int,
	note string,
	vat,
	retention float64,
) (*model.Invoice, error) {
	idString := is.getFormatedID(id)
	notes := is.cfgRepo.GetNotes()
	currency := is.cfgRepo.GetCurrency()
	payment := is.cfgRepo.GetPaymentInfo()
	dueHours := dueDays * 24

	due, err := time.ParseDuration(fmt.Sprintf("%dh", dueHours))
	if err != nil {
		return nil, err
	}

	invoice := model.NewInvoice(idString, due, currency, payment)
	invoice.Logo = is.cfgRepo.GetLogo()
	invoice.From = is.cfgRepo.GetFreelancer()
	invoice.To = *is.cRepo.Read(clientID)
	invoice.Notes = model.Notes{Default: note}
	invoice.SetTaxes(vat, retention, notes)

	err = is.iRepo.Create(invoice)
	if err != nil {
		return nil, err
	}

	return invoice, nil
}

func (is *InvoiceService) Read(id string) *model.Invoice {
	return is.iRepo.Read(id)
}

func (is *InvoiceService) Update(invoiceID string, invoice *model.Invoice) *model.Invoice {
	return is.iRepo.Update(invoiceID, invoice)
}

func (is *InvoiceService) Delete(invoiceID string) error {
	return is.iRepo.Delete(invoiceID)
}

func (is *InvoiceService) getFormatedID(id int) string {
	// TODO: Make a better solution. Maybe a repo based GetLastID
	if id == 0 {
		id = 1
		invoices := is.iRepo.List(noFilter())

		if len(invoices) > 0 {
			invoice := invoices[len(invoices)-1]

			idParts := strings.Split(invoice.ID, "-")
			id, _ = strconv.Atoi(idParts[len(idParts)-1]) //nolint:errcheck,errcheck //hack based on known format
			id++
		}
	}

	return fmt.Sprintf(is.cfgRepo.GetIDFormat(), time.Now().Format("06"), id)
}

func noFilter() repository.Filter[*model.Invoice] {
	return func(_ *model.Invoice) bool {
		return true
	}
}
