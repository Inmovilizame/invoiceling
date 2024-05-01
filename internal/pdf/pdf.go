package pdf

import (
	_ "embed"
	"fmt"
	"image"
	"os"
	"strconv"
	"strings"

	"github.com/Inmovilizame/invoiceling/pkg/model"

	"github.com/signintech/gopdf"
)

const (
	quantityColumnOffset = 360
	rateColumnOffset     = 405
	amountColumnOffset   = 480
)

const (
	subtotalLabel = "Subtotal"
	discountLabel = "Discount"
	taxLabel      = "Tax"
	totalLabel    = "Total"
)

var (
	margins = map[string]float64{"left": 40, "top": 40, "right": 40, "bottom": 40}
)

type Document struct {
	po       *gopdf.GoPdf
	lastYPos float64
}

func NewGoPdf(fonts map[string][]byte) (*Document, error) {
	pdfObject := gopdf.GoPdf{}
	pdfObject.Start(gopdf.Config{
		PageSize: *gopdf.PageSizeA4,
	})

	pdfObject.SetMargins(margins["left"], margins["top"], margins["right"], margins["bottom"])
	pdfObject.AddPage()

	for name, font := range fonts {
		err := pdfObject.AddTTFFontData(name, font)
		if err != nil {
			return nil, err
		}
	}

	return &Document{po: &pdfObject, lastYPos: pdfObject.GetY()}, nil
}

func (d *Document) SaveTo(dest string) error {
	err := d.po.WritePdf(dest)
	if err != nil {
		return err
	}

	return nil
}

func (d *Document) Logo(logo string) {
	if logo == "" {
		return
	}

	startX := d.po.GetX()
	startY := d.po.GetX()

	width, height := getImageDimension(logo)
	scaledWidth := 100.0
	scaledHeight := float64(height) * scaledWidth / float64(width)
	_ = d.po.Image(logo, startX, startY, &gopdf.Rect{W: scaledWidth, H: scaledHeight})
	d.po.Br(scaledHeight + 24)

	d.lastYPos = d.po.GetY()
}

func (d *Document) InvoiceInfo(id, date, due string) {
	d.po.SetX(420)
	d.po.SetY(40)

	_ = d.po.SetFont("Inter-Bold", "", 24)
	d.po.SetTextColor(0, 0, 0)
	_ = d.po.CellWithOption(
		&gopdf.Rect{W: 135, H: 36},
		"INVOICE",
		gopdf.CellOption{Align: gopdf.Center},
	)
	d.po.Br(36)

	_ = d.po.SetFont("Inter", "", 12)

	d.po.SetX(420)
	d.po.SetTextColor(100, 100, 100)
	_ = d.po.CellWithOption(&gopdf.Rect{W: 30, H: 20},
		"Invoice:",
		gopdf.CellOption{Align: gopdf.Right},
	)
	d.po.SetX(450)
	d.po.SetTextColor(0, 0, 0)
	_ = d.po.CellWithOption(
		&gopdf.Rect{W: 105, H: 20},
		id,
		gopdf.CellOption{Align: gopdf.Right},
	)
	d.po.Br(20)

	d.po.SetX(420)
	d.po.SetTextColor(100, 100, 100)
	_ = d.po.CellWithOption(&gopdf.Rect{W: 30, H: 20},
		"Date:",
		gopdf.CellOption{Align: gopdf.Right},
	)
	d.po.SetX(450)
	d.po.SetTextColor(0, 0, 0)
	_ = d.po.CellWithOption(
		&gopdf.Rect{W: 105, H: 20},
		date,
		gopdf.CellOption{Align: gopdf.Right},
	)
	d.po.Br(20)

	d.po.SetX(420)
	d.po.SetTextColor(100, 100, 100)
	_ = d.po.CellWithOption(&gopdf.Rect{W: 30, H: 20},
		"Due:",
		gopdf.CellOption{Align: gopdf.Right},
	)
	d.po.SetX(450)
	d.po.SetTextColor(0, 0, 0)
	_ = d.po.CellWithOption(
		&gopdf.Rect{W: 105, H: 20},
		due,
		gopdf.CellOption{Align: gopdf.Right},
	)
	d.po.Br(20)

	endY := d.po.GetY()
	if endY > d.lastYPos {
		d.lastYPos = endY
	}
}

func (d *Document) WriteHeader(logo, title, id, date string) {
	startX := d.po.GetX()
	startY := d.po.GetX()

	if logo != "" {
		width, height := getImageDimension(logo)
		scaledWidth := 100.0
		scaledHeight := float64(height) * scaledWidth / float64(width)
		_ = d.po.Image(logo, startX, startY, &gopdf.Rect{W: scaledWidth, H: scaledHeight})
		d.po.Br(scaledHeight + 24)
	}

	//logoY := d.po.GetY()

	d.po.SetX(300)
	d.po.SetY(startY)

	_ = d.po.SetFont("Inter-Bold", "", 24)
	d.po.SetTextColor(0, 0, 0)
	_ = d.po.Cell(nil, title)
	d.po.Br(36)

	d.po.SetX(300)
	_ = d.po.SetFont("Inter", "", 12)
	d.po.SetTextColor(100, 100, 100)
	_ = d.po.Cell(nil, "#")
	_ = d.po.Cell(nil, id)
	d.po.SetTextColor(150, 150, 150)
	_ = d.po.Cell(nil, "  ·  ")
	d.po.SetTextColor(100, 100, 100)
	_ = d.po.Cell(nil, date)
	d.po.Br(48)
}

func (d *Document) WriteFromTo(from model.Freelancer) {
	d.po.SetTextColor(55, 55, 55)
	d.po.SetFont("Inter", "", 14)
	d.po.Cell(nil, from.Name)
	d.po.Br(18)

	d.po.SetFont("Inter", "", 10)

	if from.Company != "" {
		d.po.SetX(300)
		d.po.Cell(nil, from.Company)
		d.po.Br(15)
	}

	if from.Address1 != "" {
		d.po.SetX(300)
		d.po.Cell(nil, from.Address1)
		d.po.Br(15)
	}

	if from.Address2 != "" {
		d.po.SetX(300)
		d.po.Cell(nil, from.Address2)
		d.po.Br(15)
	}

	if from.Phone != "" {
		d.po.SetX(300)
		d.po.Cell(nil, from.Phone)
		d.po.Br(15)
	}

	if from.VatID != "" {
		d.po.SetX(300)
		d.po.Cell(nil, from.VatID)
		d.po.Br(15)
	}

	d.po.Br(21)
	d.po.SetStrokeColor(225, 225, 225)
	d.po.Line(300, d.po.GetY(), 555, d.po.GetY())
	d.po.Br(36)
}

func WriteTitle(pdf *gopdf.GoPdf, title, id, date string) {

}

func WriteDueDate(pdf *gopdf.GoPdf, due string) {
	_ = pdf.SetFont("Inter", "", 9)
	pdf.SetTextColor(75, 75, 75)
	pdf.SetX(rateColumnOffset)
	_ = pdf.Cell(nil, "Due Date")
	pdf.SetTextColor(0, 0, 0)
	_ = pdf.SetFontSize(11)
	pdf.SetX(amountColumnOffset - 15)
	_ = pdf.Cell(nil, due)
	pdf.Br(12)
}

func WriteBillTo(pdf *gopdf.GoPdf, to string) {
	pdf.SetTextColor(75, 75, 75)
	_ = pdf.SetFont("Inter", "", 9)
	_ = pdf.Cell(nil, "BILL TO")
	pdf.Br(18)
	pdf.SetTextColor(75, 75, 75)

	formattedTo := strings.ReplaceAll(to, `\n`, "\n")
	toLines := strings.Split(formattedTo, "\n")

	for i := 0; i < len(toLines); i++ {
		if i == 0 {
			_ = pdf.SetFont("Inter", "", 15)
			_ = pdf.Cell(nil, toLines[i])
			pdf.Br(20)
		} else {
			_ = pdf.SetFont("Inter", "", 10)
			_ = pdf.Cell(nil, toLines[i])
			pdf.Br(15)
		}
	}
	pdf.Br(64)
}

func WriteHeaderRow(pdf *gopdf.GoPdf) {
	_ = pdf.SetFont("Inter", "", 9)
	pdf.SetTextColor(55, 55, 55)
	_ = pdf.Cell(nil, "ITEM")
	pdf.SetX(quantityColumnOffset)
	_ = pdf.Cell(nil, "QTY")
	pdf.SetX(rateColumnOffset)
	_ = pdf.Cell(nil, "RATE")
	pdf.SetX(amountColumnOffset)
	_ = pdf.Cell(nil, "AMOUNT")
	pdf.Br(24)
}

func WriteNotes(pdf *gopdf.GoPdf, notes string) {
	pdf.SetY(600)

	_ = pdf.SetFont("Inter", "", 9)
	pdf.SetTextColor(55, 55, 55)
	_ = pdf.Cell(nil, "NOTES")
	pdf.Br(18)
	_ = pdf.SetFont("Inter", "", 9)
	pdf.SetTextColor(0, 0, 0)

	formattedNotes := strings.ReplaceAll(notes, `\n`, "\n")
	notesLines := strings.Split(formattedNotes, "\n")

	for i := 0; i < len(notesLines); i++ {
		_ = pdf.Cell(nil, notesLines[i])
		pdf.Br(15)
	}

	pdf.Br(48)
}
func WriteFooter(pdf *gopdf.GoPdf, id string) {
	pdf.SetY(800)

	_ = pdf.SetFont("Inter", "", 10)
	pdf.SetTextColor(55, 55, 55)
	_ = pdf.Cell(nil, id)
	pdf.SetStrokeColor(225, 225, 225)
	pdf.Line(pdf.GetX()+10, pdf.GetY()+6, 550, pdf.GetY()+6)
	pdf.Br(48)
}

func WriteRow(pdf *gopdf.GoPdf, item string, quantity int, rate float64, currency string) {
	_ = pdf.SetFont("Inter", "", 11)
	pdf.SetTextColor(0, 0, 0)

	total := float64(quantity) * rate
	amount := strconv.FormatFloat(total, 'f', 2, 64)

	_ = pdf.Cell(nil, item)
	pdf.SetX(quantityColumnOffset)
	_ = pdf.Cell(nil, strconv.Itoa(quantity))
	pdf.SetX(rateColumnOffset)
	_ = pdf.Cell(nil, model.GetCurrencySymbol(currency)+strconv.FormatFloat(rate, 'f', 2, 64))
	pdf.SetX(amountColumnOffset)
	_ = pdf.Cell(nil, model.GetCurrencySymbol(currency)+amount)
	pdf.Br(24)
}

func WriteTotals(pdf *gopdf.GoPdf, subtotal float64, tax float64, discount float64, currency string) {
	pdf.SetY(600)

	writeTotal(pdf, subtotalLabel, subtotal, currency)
	if tax > 0 {
		writeTotal(pdf, taxLabel, tax, currency)
	}
	if discount > 0 {
		writeTotal(pdf, discountLabel, discount, currency)
	}
	writeTotal(pdf, totalLabel, subtotal+tax-discount, currency)
}

func writeTotal(pdf *gopdf.GoPdf, label string, total float64, currency string) {
	_ = pdf.SetFont("Inter", "", 9)
	pdf.SetTextColor(75, 75, 75)
	pdf.SetX(rateColumnOffset)
	_ = pdf.Cell(nil, label)
	pdf.SetTextColor(0, 0, 0)
	_ = pdf.SetFontSize(12)
	pdf.SetX(amountColumnOffset - 15)
	if label == totalLabel {
		_ = pdf.SetFont("Helvetica-Bold", "", 11.5)
	}
	_ = pdf.Cell(nil, model.GetCurrencySymbol(currency)+strconv.FormatFloat(total, 'f', 2, 64))
	pdf.Br(24)
}

func getImageDimension(imagePath string) (int, int) {
	file, err := os.Open(imagePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	defer file.Close()

	img, _, err := image.DecodeConfig(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", imagePath, err)
	}

	return img.Width, img.Height
}
