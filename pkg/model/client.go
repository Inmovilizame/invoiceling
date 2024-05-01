package model

import (
	"encoding/json"
	"io"
	"os"
)

type Client struct {
	Name     string `json:"name" yaml:"name"`
	VatID    string `json:"vat_id" yaml:"vat_id"`
	Address1 string `json:"address1" yaml:"address1"`
	Address2 string `json:"address2" yaml:"address2"`
}

func LoadClient(clientSrc string) (client Client, err error) {
	jsonFile, err := os.Open(clientSrc)
	if err != nil {
		return client, err
	}
	defer jsonFile.Close()

	jsonBytes, err := io.ReadAll(jsonFile)
	if err != nil {
		return client, err
	}

	err = json.Unmarshal(jsonBytes, &client)
	if err != nil {
		return client, err
	}

	return client, nil
}
