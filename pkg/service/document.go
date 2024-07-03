package service

import (
	"path/filepath"

	"github.com/Inmovilizame/invoiceling/pkg/model"
)

type Document struct {
	debug     bool
	draft     bool
	outputDir string
	renderer  RendererInterface
}

func NewDocumentService(debug, draft bool, outputDir string) (*Document, error) {
	return &Document{
		debug:     debug,
		draft:     draft,
		outputDir: outputDir,
	}, nil
}

func (d *Document) SetRenderer(r RendererInterface) {
	d.renderer = r
}

func (d *Document) Render(invoice *model.Invoice) error {
	filename := invoice.ID + ".pdf"
	if d.draft {
		filename = invoice.ID + "_DRAFT.pdf"
	}

	err := d.renderer.Render(invoice, d.draft)
	if err != nil {
		return err
	}

	dest := filepath.Join(d.outputDir, filename)

	err = d.renderer.SaveTo(dest)
	if err != nil {
		return err
	}

	return nil
}
