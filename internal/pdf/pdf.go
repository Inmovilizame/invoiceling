package pdf

import (
	_ "embed"
	"fmt"
	"github.com/Inmovilizame/invoiceling/pkg/model"
	"image"
	"os"
	"strconv"
	"strings"

	"github.com/signintech/gopdf"
)

const (
	Margin     = 40
	LineHeigth = 20
)

const (
	HeaderLogoSize   = 100
	HeaderInfoStartX = 400
	HeaderInfoWidth  = 155
)

const (
	itemDescWidth   = 300
	itemQtyWidth    = 60
	itemRateWidth   = 75
	itemAmountWidth = 80
)

type Document struct {
	po       *gopdf.GoPdf
	lastYPos float64
	debug    bool
}

func NewGoPdf(fonts map[string][]byte) (*Document, error) {
	pdfObject := gopdf.GoPdf{}
	pdfObject.Start(gopdf.Config{
		PageSize: *gopdf.PageSizeA4,
	})

	pdfObject.SetMargins(Margin, Margin, Margin, Margin)
	pdfObject.AddPage()

	for name, font := range fonts {
		err := pdfObject.AddTTFFontData(name, font)
		if err != nil {
			return nil, err
		}
	}

	return &Document{
		po:       &pdfObject,
		lastYPos: pdfObject.GetY(),
		//debug:    true,
	}, nil
}

func (d *Document) SaveTo(dest string) error {
	err := d.po.WritePdf(dest)
	if err != nil {
		return err
	}

	return nil
}

func (d *Document) Render(invoice *model.Invoice) error {
	err := d.header(invoice.Logo, invoice.ID, invoice.Due, invoice.Due)
	if err != nil {
		return err
	}

	d.po.SetY(d.lastYPos)
	d.po.SetStrokeColor(0, 0, 200)
	d.po.Line(Margin, d.po.GetY(), gopdf.PageSizeA4.W-Margin, d.po.GetY())
	d.po.Br(20)

	err = d.sendingInfo(invoice.From, invoice.To)
	if err != nil {
		return err
	}

	d.po.SetY(d.lastYPos)
	d.po.SetStrokeColor(0, 0, 200)
	d.po.Line(Margin, d.po.GetY(), gopdf.PageSizeA4.W-Margin, d.po.GetY())
	d.po.Br(20)

	err = d.items(invoice.Items, invoice.Currency, invoice.Tax)
	if err != nil {
		return err
	}

	d.po.SetY(d.lastYPos)
	d.po.Br(20)

	err = d.paymentInfo(invoice.Payment.Holder, invoice.Payment.Iban, invoice.Payment.Swift)
	if err != nil {
		return err
	}

	d.po.SetY(d.lastYPos)
	d.po.Br(20)

	err = d.notes(invoice.Note)
	if err != nil {
		return err
	}

	return nil
}

func (d *Document) header(logo, id, date, due string) error {
	if logo != "" {
		startX := d.po.GetX()
		startY := d.po.GetX()
		width, height := getImageScaledDimension(logo)

		err := d.po.Image(logo, startX, startY, &gopdf.Rect{W: width, H: height})
		if err != nil {
			return err
		}

		d.po.Br(height + 20)

		if d.po.GetY() > d.lastYPos {
			d.lastYPos = d.po.GetY()
		}
	}

	d.po.SetY(Margin)

	err := d.headingTitle()
	if err != nil {
		return err
	}

	err = d.headingInfoLine("Invoice", id)
	if err != nil {
		return err
	}

	err = d.headingInfoLine("Date", date)
	if err != nil {
		return err
	}

	err = d.headingInfoLine("Due", due)
	if err != nil {
		return err
	}

	if d.po.GetY() > d.lastYPos {
		d.lastYPos = d.po.GetY()
	}

	return nil
}

func (d *Document) sendingInfo(from *model.Freelancer, client *model.Client) error {
	startY := d.po.GetY()

	err := d.from(from)
	if err != nil {
		return nil
	}

	d.po.SetY(startY)

	err = d.to(client)
	if err != nil {
		return nil
	}

	return nil
}

func (d *Document) items(items []model.Item, currency string, tax float64) error {
	currSymbol := model.GetCurrencySymbol(currency)

	d.setSubtleNormalText()

	err := d.itemTableRow("Description", "Quantity", "Rate", "Amount")
	if err != nil {
		return err
	}

	d.setNormalText()

	subtotal := 0.

	for _, item := range items {
		itemAmount := float64(item.Quantity) * item.Rate
		subtotal += itemAmount

		err := d.itemTableRow(
			item.Description,
			strconv.Itoa(item.Quantity),
			strconv.FormatFloat(item.Rate, 'f', 2, 64)+currSymbol,
			strconv.FormatFloat(itemAmount, 'f', 2, 64)+currSymbol,
		)
		if err != nil {
			return err
		}
	}
	taxes := subtotal * tax / 100.
	total := subtotal + taxes

	d.po.Br(10)
	d.po.Line(Margin+itemDescWidth, d.po.GetY(), 595-Margin, d.po.GetY())
	d.po.Br(10)

	err = d.itemTableRow(
		"",
		"Subtotal",
		"",
		strconv.FormatFloat(subtotal, 'f', 2, 64)+currSymbol)
	if err != nil {
		return err
	}

	err = d.itemTableRow(
		"",
		"VAT",
		strconv.FormatFloat(tax, 'f', 0, 64)+"%",
		strconv.FormatFloat(taxes, 'f', 2, 64)+currSymbol)
	if err != nil {
		return err
	}

	d.setSubtleTotalText()

	err = d.itemTableRow(
		"",
		"Total",
		"",
		strconv.FormatFloat(total, 'f', 2, 64)+currSymbol)
	if err != nil {
		return err
	}

	if d.po.GetY() > d.lastYPos {
		d.lastYPos = d.po.GetY()
	}

	return nil
}

func (d *Document) paymentInfo(holder, iban, swift string) error {
	d.po.SetX(Margin + itemDescWidth)
	d.setSubtleNormalText()

	err := d.po.CellWithOption(&gopdf.Rect{W: itemQtyWidth}, "Payment Info", d.getCellOptions(gopdf.Left))
	if err != nil {
		return err
	}

	d.po.Br(20)

	d.po.SetStrokeColor(192, 192, 192)
	d.po.SetFillColor(192, 192, 192)

	err = d.po.Rectangle(
		Margin+itemDescWidth-5,
		d.po.GetY(),
		Margin+itemDescWidth+itemQtyWidth+itemRateWidth+itemAmountWidth,
		d.po.GetY()+60,
		"DF",
		0.,
		0,
	)
	if err != nil {
		return err
	}

	d.po.Br(5)
	d.po.SetX(Margin + itemDescWidth)
	d.setNormalText()

	err = d.po.CellWithOption(&gopdf.Rect{W: itemQtyWidth}, "Holder: "+holder, d.getCellOptions(gopdf.Left))
	if err != nil {
		return err
	}

	d.po.Br(18)
	d.po.SetX(Margin + itemDescWidth)

	err = d.po.CellWithOption(&gopdf.Rect{W: itemQtyWidth}, "IBAN: "+iban, d.getCellOptions(gopdf.Left))
	if err != nil {
		return err
	}

	d.po.Br(18)
	d.po.SetX(Margin + itemDescWidth)

	err = d.po.CellWithOption(&gopdf.Rect{W: itemQtyWidth}, "Swift: "+swift, d.getCellOptions(gopdf.Left))
	if err != nil {
		return err
	}

	d.po.Br(20)

	if d.po.GetY() > d.lastYPos {
		d.lastYPos = d.po.GetY()
	}

	return nil
}

func (d *Document) notes(txt string) error {
	d.setNormalText()

	lines := strings.Split(txt, "\n")

	d.po.SetY(float64(842 - Margin - len(lines)*20))

	for _, line := range lines {
		err := d.po.MultiCell(
			&gopdf.Rect{W: 595 - 2*Margin, H: 20},
			line,
		)
		if err != nil {
			return err
		}

		d.po.Br(5)
	}

	return nil
}

func (d *Document) itemTableRow(desc, qty, rate, total string) error {
	err := d.po.CellWithOption(&gopdf.Rect{W: itemDescWidth}, desc, d.getCellOptions(gopdf.Left))
	if err != nil {
		return err
	}

	err = d.po.CellWithOption(&gopdf.Rect{W: itemQtyWidth}, qty, d.getCellOptions(gopdf.Right))
	if err != nil {
		return err
	}

	err = d.po.CellWithOption(&gopdf.Rect{W: itemRateWidth}, rate, d.getCellOptions(gopdf.Right))
	if err != nil {
		return err
	}

	err = d.po.CellWithOption(&gopdf.Rect{W: itemAmountWidth}, total, d.getCellOptions(gopdf.Right))
	if err != nil {
		return err
	}

	d.po.Br(20)

	return nil
}

func (d *Document) headingTitle() error {
	d.po.SetX(HeaderInfoStartX) // have 115 width points

	err := d.setTitleText()
	if err != nil {
		return err
	}

	err = d.po.CellWithOption(
		&gopdf.Rect{W: HeaderInfoWidth},
		"INVOICE",
		d.getCellOptions(gopdf.Center),
	)
	if err != nil {
		return err
	}

	d.po.Br(36)

	return nil
}

func (d *Document) headingInfoLine(key, value string) error {
	d.po.SetX(HeaderInfoStartX)
	d.setSubtleNormalText()

	err := d.po.CellWithOption(&gopdf.Rect{W: 45},
		key,
		d.getCellOptions(gopdf.Left),
	)
	if err != nil {
		return err
	}

	err = d.po.CellWithOption(&gopdf.Rect{W: 10},
		":",
		d.getCellOptions(gopdf.Left),
	)
	if err != nil {
		return err
	}

	d.setNormalText()

	err = d.po.CellWithOption(
		&gopdf.Rect{W: 100},
		value,
		d.getCellOptions(gopdf.Right),
	)
	if err != nil {
		return err
	}

	d.po.Br(20)

	return nil
}

func (d *Document) from(from *model.Freelancer) error {
	d.setSubtleNormalText()
	_ = d.po.Cell(&gopdf.Rect{W: 250}, "From")
	d.po.Br(20)

	d.setNormalText()
	_ = d.po.Cell(&gopdf.Rect{W: 250}, from.Name)
	d.po.Br(18)

	if from.Company != "" {
		_ = d.po.Cell(&gopdf.Rect{W: 250}, from.Company)
		d.po.Br(18)
	}

	if from.Address1 != "" {
		_ = d.po.Cell(&gopdf.Rect{W: 250}, from.Address1)
		d.po.Br(18)
	}

	if from.Address2 != "" {
		_ = d.po.Cell(&gopdf.Rect{W: 250}, from.Address2)
		d.po.Br(18)
	}

	if from.Phone != "" {
		_ = d.po.Cell(&gopdf.Rect{W: 250}, from.Phone)
		d.po.Br(18)
	}

	if from.VatID != "" {
		_ = d.po.Cell(&gopdf.Rect{W: 250}, from.VatID)
		d.po.Br(18)
	}

	d.po.Br(20)

	endY := d.po.GetY()
	if endY > d.lastYPos {
		d.lastYPos = endY
	}

	return nil
}

func (d *Document) to(client *model.Client) error {
	d.po.SetX(HeaderInfoStartX)

	d.setSubtleNormalText()
	_ = d.po.Cell(&gopdf.Rect{W: 155}, "To")
	d.po.Br(20)

	d.po.SetX(HeaderInfoStartX)
	d.setNormalText()
	_ = d.po.Cell(&gopdf.Rect{W: 250}, client.Name)
	d.po.Br(18)

	d.po.SetX(HeaderInfoStartX)
	_ = d.po.Cell(&gopdf.Rect{W: 250}, client.VatID)
	d.po.Br(18)

	if client.Address1 != "" {
		d.po.SetX(HeaderInfoStartX)
		_ = d.po.Cell(&gopdf.Rect{W: 250}, client.Address1)
		d.po.Br(18)
	}

	if client.Address2 != "" {
		d.po.SetX(HeaderInfoStartX)
		_ = d.po.Cell(&gopdf.Rect{W: 250}, client.Address2)
		d.po.Br(18)
	}

	endY := d.po.GetY()
	if endY > d.lastYPos {
		d.lastYPos = endY
	}

	return nil
}

func (d *Document) getCellOptions(align int) gopdf.CellOption {
	co := gopdf.CellOption{Align: align}
	if d.debug {
		co.Border = gopdf.AllBorders
	}

	return co
}

func (d *Document) setNormalText() {
	d.po.SetTextColor(24, 24, 24)

	err := d.po.SetFont("Inter", "", 10)
	if err != nil {
		fmt.Println("Error Loading font: 'Inter'")
	}
}

func (d *Document) setSubtleNormalText() {
	d.po.SetTextColor(128, 128, 192)

	err := d.po.SetFont("Inter", "", 12)
	if err != nil {
		fmt.Println("Error Loading font: 'Inter'")
	}

}

func (d *Document) setSubtleTotalText() {
	d.po.SetTextColor(128, 128, 192)

	err := d.po.SetFont("Inter-Bold", "", 14)
	if err != nil {
		fmt.Println("Error Loading font: 'Inter-Bold'")
	}
}

func (d *Document) setTitleText() error {
	d.po.SetTextColor(24, 24, 24)

	err := d.po.SetFont("Inter-Bold", "", 24)
	if err != nil {
		return err
	}

	return nil
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

func getImageScaledDimension(imagePath string) (scaledWidth, scaledHeight float64) {
	file, err := os.Open(imagePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	defer file.Close()

	img, _, err := image.DecodeConfig(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", imagePath, err)
	}

	scaledHeight = HeaderLogoSize
	scaledWidth = float64(img.Width) * scaledHeight / float64(img.Height)

	return
}
