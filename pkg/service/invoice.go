package service

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Inmovilizame/invoiceling/pkg/model"
)

type DateFormat string

const (
	DF_Text DateFormat = "Jan 02, 2006"
	DF_YMD  DateFormat = "2006-01-02"
	DF_DMY  DateFormat = "02/01/2006"
)

type InvoiceRepo interface {
	List() []*model.Invoice
	Create(invoice *model.Invoice) error
	Read(invoiceID string) *model.Invoice
	Update(invoiceID string, invoice *model.Invoice) *model.Invoice
	Delete(invoiceID string) error
}

type Invoice struct {
	currency string
	idFormat string
	logo     string
	iRepo    InvoiceRepo
	cRepo    ClientRepo
}

func NewInvoice(currency, idFormat, logoPath string, iRepo InvoiceRepo, cRepo ClientRepo) *Invoice {
	return &Invoice{
		currency: currency,
		idFormat: idFormat,
		logo:     logoPath,
		iRepo:    iRepo,
		cRepo:    cRepo,
	}
}

func (is *Invoice) List() []*model.Invoice {
	return is.iRepo.List()
}

func (is *Invoice) Create(id int, me *model.Freelancer, to *model.Client, dueDays int, dateFormat DateFormat) (*model.Invoice, error) {
	now := time.Now()
	due := now.AddDate(0, 0, dueDays)
	invoiceID := is.getFormatedID(id)

	items := []model.Item{
		{
			Description: "Product Description",
			Quantity:    1,
			Rate:        1.0,
		},
	}

	invoice := &model.Invoice{
		ID:       invoiceID,
		Status:   "CREATED",
		Logo:     is.logo,
		From:     me,
		To:       to,
		Date:     now.Format(string(dateFormat)),
		Due:      due.Format(string(dateFormat)),
		Items:    items,
		Tax:      0,
		Discount: 0,
		Currency: is.currency,
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
		invoices := is.iRepo.List()
		invoice := invoices[len(invoices)-1]

		idParts := strings.Split(invoice.ID, "-")
		id, _ = strconv.Atoi(idParts[len(idParts)-1]) //nolint:errcheck,errcheck hack
		id++
	}

	return fmt.Sprintf(is.idFormat, time.Now().Format("06"), id)
}
