package repository

import (
	"encoding/json"
	"io"
	"os"

	"github.com/Inmovilizame/invoiceling/pkg/model"
)

type Filter[T any] func(T) bool

func readClientFromFile(invoicePath string) (*model.Client, error) {
	jsonFile, err := os.Open(invoicePath)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	jsonBytes, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	client := &model.Client{}

	err = json.Unmarshal(jsonBytes, client)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func readInvoiceFromFile(invoicePath string) (*model.Invoice, error) {
	jsonFile, err := os.Open(invoicePath)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	jsonBytes, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	invoice := &model.Invoice{}

	err = json.Unmarshal(jsonBytes, invoice)
	if err != nil {
		return nil, err
	}

	return invoice, nil
}

func checkFileExists(pathname string) bool {
	_, err := os.Stat(pathname)
	return !os.IsNotExist(err)
}
