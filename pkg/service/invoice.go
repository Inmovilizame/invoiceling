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

const (
	DFText DateFormat = "Jan 02, 2006"
	DFYMD  DateFormat = "2006-01-02"
	DFDMY  DateFormat = "02/01/2006"
)

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
	dateFormat DateFormat,
	note string,
	vat,
	retention float64,
) (*model.Invoice, error) {
	now := time.Now()
	due := now.AddDate(0, 0, dueDays)

	invoice := &model.Invoice{
		ID:       is.getFormatedID(id),
		Status:   "CREATED",
		Logo:     is.cfgRepo.GetLogo(),
		From:     is.cfgRepo.GetFreelancer(),
		To:       *is.cRepo.Read(clientID),
		Date:     now.Format(string(dateFormat)),
		Due:      due.Format(string(dateFormat)),
		Items:    []model.Item{},
		Tax:      model.TaxInfo{},
		Discount: 0,
		Currency: is.cfgRepo.GetCurrency(),
		Payment:  is.cfgRepo.GetPaymentInfo(),
		Notes:    model.Notes{Default: note},
	}

	notes := is.cfgRepo.GetNotes()

	invoice.Tax.Vat = vat
	if vat == 0 {
		invoice.Notes.Vat0 = notes["vat_0"]
	}

	invoice.Tax.Retention = retention
	if retention != 0 {
		invoice.Notes.RetentionNot0 = notes["retention_not_0"]
	}

	err := is.iRepo.Create(invoice)
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
