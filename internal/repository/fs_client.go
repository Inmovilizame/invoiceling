package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Inmovilizame/invoiceling/pkg/model"
)

type FsClient struct {
	basePath string
}

func NewFsClient(baseDir string) *FsClient {
	basePath, err := filepath.Abs(baseDir)
	if err != nil {
		fmt.Printf("Error while getting absolute path for invoice dir. %v", err)
		return nil
	}

	return &FsClient{
		basePath: basePath,
	}
}

func (fc *FsClient) List(filter Filter[*model.Client]) []*model.Client {
	clients := make([]*model.Client, 0)

	files, err := os.ReadDir(fc.basePath)
	if err != nil {
		fmt.Printf("Error while opening invoice dir. %v", err)
		return clients
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileExt := filepath.Ext(file.Name())
		if fileExt != ".json" {
			continue
		}

		clientPath := filepath.Join(fc.basePath, file.Name())

		client, err := readClientFromFile(clientPath)
		if err != nil {
			fmt.Printf("Error while loading invoice %s. %v", file.Name(), err)
			continue
		}

		if filter(client) {
			clients = append(clients, client)
		}
	}

	return clients
}

func (fc *FsClient) Create(client *model.Client) error {
	jsonBytes, err := json.MarshalIndent(client, "", "  ")
	if err != nil {
		return err
	}

	clientPath := filepath.Join(fc.basePath, client.ID+".json")
	if checkFileExists(clientPath) {
		return errors.New("client already exists")
	}

	err = os.WriteFile(clientPath, jsonBytes, rwMask)
	if err != nil {
		return err
	}

	return nil
}

func (fc *FsClient) Read(clientID string) *model.Client {
	clientPath := filepath.Join(fc.basePath, clientID+".json")

	client, err := readClientFromFile(clientPath)
	if err != nil {
		fmt.Printf("Error while loading invoice %s. %v", clientID, err)
		return nil
	}

	return client
}

func (fc *FsClient) Update(clientID string, client *model.Client) *model.Client {
	jsonBytes, err := json.MarshalIndent(client, "", "  ")
	if err != nil {
		fmt.Printf("Error while marshaling client %s. %v", clientID, err)
		return nil
	}

	clientPath := filepath.Join(fc.basePath, clientID+".json")

	err = os.WriteFile(clientPath, jsonBytes, rwMask)
	if err != nil {
		fmt.Printf("Error while updating client %s. %v", clientID, err)
		return nil
	}

	return client
}

func (fc *FsClient) Delete(clientID string) error {
	clientPath := filepath.Join(fc.basePath, clientID+".json")

	err := os.WriteFile(clientPath, []byte{}, roMask)
	if err != nil {
		fmt.Printf("Error while updating invoice %s. %v", clientID, err)
		return err
	}

	return nil
}
