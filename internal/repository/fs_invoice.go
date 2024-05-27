package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Inmovilizame/invoiceling/pkg/model"
)

type FsInvoice struct {
	basePath string
}

func NewFsInvoice(baseDir string) *FsInvoice {
	basePath, err := filepath.Abs(baseDir)
	if err != nil {
		fmt.Printf("Error while getting absolute path for invoice dir. %v", err)
		return nil
	}

	return &FsInvoice{
		basePath: basePath,
	}
}

func (fi *FsInvoice) List(filter Filter[*model.Invoice]) []*model.Invoice {
	invoices := make([]*model.Invoice, 0)

	files, err := os.ReadDir(fi.basePath)
	if err != nil {
		fmt.Printf("Error while opening invoice dir. %v", err)
		return invoices
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileExt := filepath.Ext(file.Name())
		if fileExt != ".json" {
			continue
		}

		invoicePath := filepath.Join(fi.basePath, file.Name())

		invoice, err := readInvoiceFromFile(invoicePath)
		if err != nil {
			fmt.Printf("Error while loading invoice %s. %v", file.Name(), err)
			continue
		}

		if filter(invoice) {
			invoices = append(invoices, invoice)
		}
	}

	return invoices
}

func (fi *FsInvoice) Create(invoice *model.Invoice) error {
	jsonBytes, err := json.MarshalIndent(invoice, "", "  ")
	if err != nil {
		return err
	}

	invoicePath := filepath.Join(fi.basePath, invoice.ID+".json")
	if checkFileExists(invoicePath) {
		return errors.New("invoice already exists")
	}

	err = os.WriteFile(invoicePath, jsonBytes, rw_mask)
	if err != nil {
		return err
	}

	return nil
}

func (fi *FsInvoice) Read(invoiceID string) *model.Invoice {
	invoicePath := filepath.Join(fi.basePath, invoiceID+".json")

	invoice, err := readInvoiceFromFile(invoicePath)
	if err != nil {
		fmt.Printf("Error while loading invoice %s. %v", invoiceID, err)
		return nil
	}

	return invoice
}

func (fi *FsInvoice) Update(invoiceID string, invoice *model.Invoice) *model.Invoice {
	jsonBytes, err := json.MarshalIndent(invoice, "", "  ")
	if err != nil {
		fmt.Printf("Error while marshaling invoice %s. %v", invoiceID, err)
		return nil
	}

	invoicePath := filepath.Join(fi.basePath, invoiceID+".json")

	err = os.WriteFile(invoicePath, jsonBytes, rw_mask)
	if err != nil {
		fmt.Printf("Error while updating invoice %s. %v", invoiceID, err)
		return nil
	}

	return invoice
}

func (fi *FsInvoice) Delete(invoiceID string) error {
	invoicePath := filepath.Join(fi.basePath, invoiceID+".json")

	err := os.WriteFile(invoicePath, []byte{}, rw_mask)
	if err != nil {
		fmt.Printf("Error while deleting invoice %s. %v", invoiceID, err)
		return err
	}

	return nil
}
