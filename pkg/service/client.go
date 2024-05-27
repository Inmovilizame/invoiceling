package service

import (
	"errors"
	"regexp"
	"strings"

	"github.com/Inmovilizame/invoiceling/internal/repository"
	"github.com/Inmovilizame/invoiceling/pkg/model"
)

const (
	countryCodeLength = 2
)

var (
	ErrNotValidVatFormat = errors.New("vat: format not valid")
	ErrCountryNotFound   = errors.New("vat: country not found")
)

type ClientRepo interface {
	List(filter repository.Filter[*model.Client]) []*model.Client
	Create(client *model.Client) error
	Read(clientID string) *model.Client
	Update(clientID string, invoice *model.Client) *model.Client
	Delete(clientID string) error
}

type Client struct {
	repo ClientRepo
}

func NewClientService(repo ClientRepo) *Client {
	return &Client{
		repo: repo,
	}
}

func (cs *Client) List(filter repository.Filter[*model.Client]) []*model.Client {
	return cs.repo.List(filter)
}

func (cs *Client) Create(id, name, vatID, address1, address2 string) error {
	err := ValidateNumberFormat(vatID)
	if err != nil {
		return err
	}

	if id == "" {
		b := strings.Builder{}
		b.WriteString("client")
		b.WriteString("-")
		b.WriteString(vatID)
	}

	client := &model.Client{
		ID:       id,
		Name:     name,
		VatID:    vatID,
		Address1: address1,
		Address2: address2,
	}

	return cs.repo.Create(client)
}

func (cs *Client) Read(id string) *model.Client {
	return cs.repo.Read(id)
}

func (cs *Client) Update(invoiceID string, invoice *model.Client) *model.Client {
	return cs.repo.Update(invoiceID, invoice)
}

func (cs *Client) Delete(invoiceID string) error {
	return cs.repo.Delete(invoiceID)
}

// ValidateNumberFormat validates a VAT number by its format.
func ValidateNumberFormat(n string) error {
	n = strings.ToUpper(n)
	if len(n) <= countryCodeLength {
		return ErrNotValidVatFormat
	}

	patterns := map[string]string{
		"AT": `U[A-Z0-9]{8}`,
		"BE": `(0[0-9]{9}|[0-9]{10})`,
		"BG": `[0-9]{9,10}`,
		"CH": `(?:E(?:-| )[0-9]{3}(?:\.| )[0-9]{3}(?:\.| )[0-9]{3}( MWST)?|E[0-9]{9}(?:MWST)?)`,
		"CY": `[0-9]{8}[A-Z]`,
		"CZ": `[0-9]{8,10}`,
		"DE": `[0-9]{9}`,
		"DK": `[0-9]{8}`,
		"EE": `[0-9]{9}`,
		"EL": `[0-9]{9}`,
		"ES": `[A-Z][0-9]{7}[A-Z]|[0-9]{8}[A-Z]|[A-Z][0-9]{8}`,
		"FI": `[0-9]{8}`,
		"FR": `([A-Z]{2}|[0-9]{2})[0-9]{9}`,
		"GB": `[0-9]{9}|[0-9]{12}|(GD|HA)[0-9]{3}`,
		"HR": `[0-9]{11}`,
		"HU": `[0-9]{8}`,
		"IE": `[A-Z0-9]{7}[A-Z]|[A-Z0-9]{7}[A-W][A-I]`,
		"IT": `[0-9]{11}`,
		"LT": `([0-9]{9}|[0-9]{12})`,
		"LU": `[0-9]{8}`,
		"LV": `[0-9]{11}`,
		"MT": `[0-9]{8}`,
		"NL": `[0-9]{9}B[0-9]{2}`,
		"PL": `[0-9]{10}`,
		"PT": `[0-9]{9}`,
		"RO": `[0-9]{2,10}`,
		"SE": `[0-9]{12}`,
		"SI": `[0-9]{8}`,
		"SK": `[0-9]{10}`,
	}

	pattern, ok := patterns[n[0:2]]
	if !ok {
		return ErrCountryNotFound
	}

	matched, err := regexp.MatchString(pattern, n[2:])
	if !matched {
		return ErrNotValidVatFormat
	}

	return err
}
