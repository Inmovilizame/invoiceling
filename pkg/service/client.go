package service

import (
	"github.com/Inmovilizame/invoiceling/internal/repository"
	"github.com/Inmovilizame/invoiceling/pkg/model"
)

type ClientRepo interface {
	List() []*model.Client
	Create(client *model.Client) error
	Read(clientID string) *model.Client
	Update(clientID string, invoice *model.Client) *model.Client
	Delete(clientID string) error
}

type Client struct {
	repo ClientRepo
}

func NewClientFS(basePath string) *Client {
	return &Client{
		repo: repository.NewFsClient(basePath),
	}
}

func (cs *Client) List() []*model.Client {
	return cs.repo.List()
}

func (cs *Client) Create(invoice *model.Client) error {
	return cs.repo.Create(invoice)
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
