package service

import (
	"fmt"
	"image"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Inmovilizame/invoiceling/pkg/model"

	"github.com/signintech/gopdf"
)

const (
	DocWidth   = 595
	LineHeigth = 20
	Margin     = 40
)

const (
	HeaderLogoSize   = 100
	HeaderInfoStartX = 400
	HeaderInfoWidth  = 155
)

const (
	ToStart = 320
)

const (
	ItemDescWidth   = 300
	ItemQtyWidth    = 60
	ItemRateWidth   = 75
	ItemAmountWidth = 80
)

type Document struct {
	po        *gopdf.GoPdf
	lastYPos  float64
	debug     bool
	outputDir string
}

func NewInvoiceRender(fonts map[string][]byte, outputDir string, debug bool) (*Document, error) {
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
		po:        &pdfObject,
		lastYPos:  pdfObject.GetY(),
		outputDir: outputDir,
		debug:     debug,
	}, nil
}

func (d *Document) SaveTo(filename string) error {
	dest := filepath.Join(d.outputDir, filename)
	err := d.po.WritePdf(dest)

	if err != nil {
		return err
	}

	return nil
}

func schemaBlueColor() (red, green, blue uint8) {
	return 0, 0, 200 //nolint:gomnd // static value for color schema
}

func schemaGrayColor() (red, green, blue uint8) {
	return 192, 192, 192 //nolint:gomnd // static value for color schema
}

func (d *Document) Render(invoice *model.Invoice) error {
	err := d.header(invoice.Logo, invoice.ID, invoice.Date, invoice.Due)
	if err != nil {
		return err
	}

	d.po.SetY(d.lastYPos)
	d.po.SetStrokeColor(schemaBlueColor())
	d.po.Line(Margin, d.po.GetY(), gopdf.PageSizeA4.W-Margin, d.po.GetY())
	d.po.Br(LineHeigth)

	err = d.sendingInfo(invoice.From, invoice.To)
	if err != nil {
		return err
	}

	d.po.SetY(d.lastYPos)
	d.po.SetStrokeColor(schemaBlueColor())
	d.po.Line(Margin, d.po.GetY(), gopdf.PageSizeA4.W-Margin, d.po.GetY())
	d.po.Br(LineHeigth)

	err = d.items(invoice.Items, invoice.Currency, invoice.Tax, invoice.Retention)
	if err != nil {
		return err
	}

	d.po.SetY(d.lastYPos)
	d.po.Br(LineHeigth)

	err = d.paymentInfo(invoice.Payment.Holder, invoice.Payment.Iban, invoice.Payment.Swift)
	if err != nil {
		return err
	}

	d.po.SetY(d.lastYPos)
	d.po.Br(LineHeigth)

	err = d.notes(invoice.Note)
	if err != nil {
		return err
	}

	return nil
}

func (d *Document) header(logo, id, date, due string) error {
	if logo != "" {
		startX := d.po.GetX()
		startY := d.po.GetY()
		width, height := getImageScaledDimension(logo)

		err := d.po.Image(logo, startX, startY, &gopdf.Rect{W: width, H: height})
		if err != nil {
			return err
		}

		d.po.Br(height + LineHeigth)

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

func (d *Document) items(items []model.Item, currency string, tax, retention float64) error {
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
	retentions := subtotal * retention / 100.
	total := subtotal + taxes - retentions

	d.po.Br(LineHeigth / 2)
	d.po.Line(Margin+ItemDescWidth, d.po.GetY(), DocWidth-Margin, d.po.GetY())
	d.po.Br(LineHeigth / 2)

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

	if retention != 0. {
		err = d.itemTableRow(
			"",
			"IRPF",
			"-"+strconv.FormatFloat(retention, 'f', 0, 64)+"%",
			strconv.FormatFloat(retentions, 'f', 2, 64)+currSymbol)
		if err != nil {
			return err
		}
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
	d.po.SetX(Margin + ItemDescWidth)
	d.setSubtleNormalText()

	err := d.po.CellWithOption(&gopdf.Rect{W: ItemQtyWidth}, "Payment Info", d.getCellOptions(gopdf.Left))
	if err != nil {
		return err
	}

	d.po.Br(LineHeigth)

	d.po.SetStrokeColor(schemaGrayColor())
	d.po.SetFillColor(schemaGrayColor())

	err = d.po.Rectangle(
		Margin+ItemDescWidth-5,
		d.po.GetY(),
		Margin+ItemDescWidth+ItemQtyWidth+ItemRateWidth+ItemAmountWidth,
		d.po.GetY()+60,
		"DF",
		0.,
		0,
	)
	if err != nil {
		return err
	}

	d.po.Br(5)
	d.po.SetX(Margin + ItemDescWidth)
	d.setNormalText()

	err = d.po.CellWithOption(&gopdf.Rect{W: ItemQtyWidth}, "Holder: "+holder, d.getCellOptions(gopdf.Left))
	if err != nil {
		return err
	}

	d.po.Br(18)
	d.po.SetX(Margin + ItemDescWidth)

	err = d.po.CellWithOption(&gopdf.Rect{W: ItemQtyWidth}, "IBAN: "+iban, d.getCellOptions(gopdf.Left))
	if err != nil {
		return err
	}

	d.po.Br(18)
	d.po.SetX(Margin + ItemDescWidth)

	err = d.po.CellWithOption(&gopdf.Rect{W: ItemQtyWidth}, "Swift: "+swift, d.getCellOptions(gopdf.Left))
	if err != nil {
		return err
	}

	d.po.Br(LineHeigth)

	if d.po.GetY() > d.lastYPos {
		d.lastYPos = d.po.GetY()
	}

	return nil
}

func (d *Document) notes(notes []string) error {
	d.setNormalText()

	d.po.SetY(float64(842 - Margin - len(notes)*20))

	for _, line := range notes {
		err := d.po.MultiCell(
			&gopdf.Rect{W: DocWidth - 2*Margin, H: LineHeigth},
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
	err := d.po.CellWithOption(&gopdf.Rect{W: ItemDescWidth}, desc, d.getCellOptions(gopdf.Left))
	if err != nil {
		return err
	}

	err = d.po.CellWithOption(&gopdf.Rect{W: ItemQtyWidth}, qty, d.getCellOptions(gopdf.Right))
	if err != nil {
		return err
	}

	err = d.po.CellWithOption(&gopdf.Rect{W: ItemRateWidth}, rate, d.getCellOptions(gopdf.Right))
	if err != nil {
		return err
	}

	err = d.po.CellWithOption(&gopdf.Rect{W: ItemAmountWidth}, total, d.getCellOptions(gopdf.Right))
	if err != nil {
		return err
	}

	d.po.Br(LineHeigth)

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

	d.po.Br(LineHeigth)

	return nil
}

func (d *Document) from(from *model.Freelancer) error {
	d.setSubtleNormalText()
	_ = d.po.Cell(&gopdf.Rect{W: 250}, "From")
	d.po.Br(LineHeigth)

	d.setNormalText()
	_ = d.po.Cell(&gopdf.Rect{W: 250}, from.Name)
	d.po.Br(18)

	if from.Company != "" {
		_ = d.po.Cell(&gopdf.Rect{W: 250}, from.Company)
		d.po.Br(18)
	}

	if from.VatID != "" {
		_ = d.po.Cell(&gopdf.Rect{W: 250}, from.VatID)
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

	d.po.Br(LineHeigth)

	endY := d.po.GetY()
	if endY > d.lastYPos {
		d.lastYPos = endY
	}

	return nil
}

func (d *Document) to(client *model.Client) error {
	d.po.SetX(ToStart)

	d.setSubtleNormalText()
	_ = d.po.Cell(&gopdf.Rect{W: 155}, "To")
	d.po.Br(LineHeigth)

	d.po.SetX(ToStart)
	d.setNormalText()
	_ = d.po.Cell(&gopdf.Rect{W: 250}, client.Name)
	d.po.Br(18)

	d.po.SetX(ToStart)
	_ = d.po.Cell(&gopdf.Rect{W: 250}, client.VatID)
	d.po.Br(18)

	if client.Address1 != "" {
		d.po.SetX(ToStart)
		_ = d.po.Cell(&gopdf.Rect{W: 250}, client.Address1)
		d.po.Br(18)
	}

	if client.Address2 != "" {
		d.po.SetX(ToStart)
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
