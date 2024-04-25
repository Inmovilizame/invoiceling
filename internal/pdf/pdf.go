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

type Document struct {
	PDF *gopdf.GoPdf
}

func NewGoPdf(fonts map[string][]byte) Document {
	pdfObject := gopdf.GoPdf{}
	pdfObject.Start(gopdf.Config{
		PageSize: *gopdf.PageSizeA4,
	})

	pdfObject.SetMargins(40, 40, 40, 40)
	pdfObject.AddPage()

	for name, font := range fonts {
		_ = pdfObject.AddTTFFontData(name, font)
	}

	return Document{&pdfObject}
}

func (d Document) WriteLogo(logo string, from string) {
	if logo != "" {
		width, height := getImageDimension(logo)
		scaledWidth := 100.0
		scaledHeight := float64(height) * scaledWidth / float64(width)
		_ = d.PDF.Image(logo, d.PDF.GetX(), d.PDF.GetY(), &gopdf.Rect{W: scaledWidth, H: scaledHeight})
		d.PDF.Br(scaledHeight + 24)
	}
	d.PDF.SetTextColor(55, 55, 55)

	formattedFrom := strings.ReplaceAll(from, `\n`, "\n")
	fromLines := strings.Split(formattedFrom, "\n")

	for i := 0; i < len(fromLines); i++ {
		if i == 0 {
			_ = d.PDF.SetFont("Inter", "", 12)
			_ = d.PDF.Cell(nil, fromLines[i])
			d.PDF.Br(18)
		} else {
			_ = d.PDF.SetFont("Inter", "", 10)
			_ = d.PDF.Cell(nil, fromLines[i])
			d.PDF.Br(15)
		}
	}
	d.PDF.Br(21)
	d.PDF.SetStrokeColor(225, 225, 225)
	d.PDF.Line(d.PDF.GetX(), d.PDF.GetY(), 260, d.PDF.GetY())
	d.PDF.Br(36)
}

func WriteTitle(pdf *gopdf.GoPdf, title, id, date string) {
	_ = pdf.SetFont("Helvetica-Bold", "", 24)
	pdf.SetTextColor(0, 0, 0)
	_ = pdf.Cell(nil, title)
	pdf.Br(36)
	_ = pdf.SetFont("Inter", "", 12)
	pdf.SetTextColor(100, 100, 100)
	_ = pdf.Cell(nil, "#")
	_ = pdf.Cell(nil, id)
	pdf.SetTextColor(150, 150, 150)
	_ = pdf.Cell(nil, "  ·  ")
	pdf.SetTextColor(100, 100, 100)
	_ = pdf.Cell(nil, date)
	pdf.Br(48)
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
