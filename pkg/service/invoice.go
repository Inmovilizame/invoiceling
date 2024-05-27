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
	DF_Text DateFormat = "Jan 02, 2006"
	DF_YMD  DateFormat = "2006-01-02"
	DF_DMY  DateFormat = "02/01/2006"
)

type InvoiceRepo interface {
	List(filter repository.Filter[*model.Invoice]) []*model.Invoice
	Create(invoice *model.Invoice) error
	Read(invoiceID string) *model.Invoice
	Update(invoiceID string, invoice *model.Invoice) *model.Invoice
	Delete(invoiceID string) error
}

type CfgRepo interface {
	GetPdfOutputDir() string
	GetCurrency() string
	GetIDFormat() string
	GetLogo() string
	GetFreelancer() *model.Freelancer
	GetPaymentInfo() *model.Payment
}

type Invoice struct {
	iRepo   InvoiceRepo
	cRepo   ClientRepo
	cfgRepo CfgRepo
}

func NewInvoice(iRepo InvoiceRepo, cRepo ClientRepo, cfgRepo CfgRepo) *Invoice {
	return &Invoice{
		iRepo:   iRepo,
		cRepo:   cRepo,
		cfgRepo: cfgRepo,
	}
}

func (is *Invoice) List(filter repository.Filter[*model.Invoice]) []*model.Invoice {
	return is.iRepo.List(filter)
}

func (is *Invoice) Create(id int, me *model.Freelancer, to *model.Client, dueDays int, dateFormat DateFormat, note string) (*model.Invoice, error) {
	now := time.Now()
	due := now.AddDate(0, 0, dueDays)
	invoiceID := is.getFormatedID(id)

	invoice := &model.Invoice{
		ID:        invoiceID,
		Status:    "CREATED",
		Logo:      is.cfgRepo.GetLogo(),
		From:      me,
		To:        to,
		Date:      now.Format(string(dateFormat)),
		Due:       due.Format(string(dateFormat)),
		Items:     []model.Item{},
		Tax:       0,
		Discount:  0,
		Retention: 0,
		Currency:  is.cfgRepo.GetCurrency(),
		Payment:   is.cfgRepo.GetPaymentInfo(),
		Note: []string{
			"Thank you for your business. Please add the invoice number to your payment description.",
		},
	}

	if note != "" {
		invoice.Note = append(invoice.Note, note)
	}

	err := is.iRepo.Create(invoice)
	if err != nil {
		return nil, err
	}

	return invoice, nil
}

func (is *Invoice) Read(id string) *model.Invoice {
	return is.iRepo.Read(id)
}

func (is *Invoice) Update(invoiceID string, invoice *model.Invoice) *model.Invoice {
	return is.iRepo.Update(invoiceID, invoice)
}

func (is *Invoice) Delete(invoiceID string) error {
	return is.iRepo.Delete(invoiceID)
}

func (is *Invoice) getFormatedID(id int) string {
	// TODO: Make a better solution. Maybe a repo based GetLastID
	if id == 0 {
		id = 1
		invoices := is.iRepo.List(noFilter())

		if len(invoices) > 0 {
			invoice := invoices[len(invoices)-1]

			idParts := strings.Split(invoice.ID, "-")
			id, _ = strconv.Atoi(idParts[len(idParts)-1]) //nolint:errcheck,errcheck hack
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
