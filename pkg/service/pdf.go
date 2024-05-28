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
	FontNormalSize       = 10
	FontSubtleNormalSize = 12
	FontSubtleTotalSize  = 14
	FontTitleSize        = 24
	LineHeight           = 20
	Margin               = 40
)

const (
	HeaderInfoStartX    = 400
	HeaderInfoWidth     = 155
	HeaderInfoName      = 45
	HeaderInfoSeparator = 10
	HeaderInfoValue     = 100
	HeaderLogoSize      = 100
	HeaderMinHeight     = 160
)

const (
	FromToLineHeight = 18
	FromWidth        = 250
	ToStart          = 320
	ToWidth          = 275
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

func (d *Document) Render(invoice *model.Invoice) error {
	err := d.header(invoice.Logo, invoice.ID, invoice.Date, invoice.Due)
	if err != nil {
		return err
	}

	d.po.SetY(d.lastYPos)
	d.po.SetStrokeColor(colorBlue())
	d.po.Line(Margin, d.po.GetY(), gopdf.PageSizeA4.W-Margin, d.po.GetY())
	d.po.Br(LineHeight)

	err = d.sendingInfo(&invoice.From, &invoice.To)
	if err != nil {
		return err
	}

	d.po.SetY(d.lastYPos)
	d.po.SetStrokeColor(colorBlue())
	d.po.Line(Margin, d.po.GetY(), gopdf.PageSizeA4.W-Margin, d.po.GetY())
	d.po.Br(LineHeight)

	err = d.items(
		invoice.Items,
		invoice.Tax,
		invoice.Currency,
		invoice.Payment)
	if err != nil {
		return err
	}

	d.po.SetY(d.lastYPos)
	d.po.Br(LineHeight)

	err = d.notes(invoice.Notes)
	if err != nil {
		return err
	}

	return nil
}

func (d *Document) header(logo, id, date, due string) error {
	err := d.headingTitle()
	if err != nil {
		return err
	}

	err = d.headingInfoLine("InvoiceService", id)
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

	d.po.SetX(HeaderInfoStartX)
	d.po.SetY(Margin)

	if logo != "" {
		startX := d.po.GetX()
		startY := d.po.GetY()
		width, height := getImageScaledDimension(logo)

		err := d.po.Image(logo, startX, startY, &gopdf.Rect{W: width, H: height})
		if err != nil {
			return err
		}

		d.po.Br(height + LineHeight)

		if d.po.GetY() > d.lastYPos {
			d.lastYPos = d.po.GetY()
		}
	}

	if d.lastYPos < HeaderMinHeight {
		d.lastYPos = HeaderMinHeight
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

//nolint:funlen //TODO fix func length
func (d *Document) items(items []model.Item, tax model.TaxInfo, currency string, payment model.Payment) error {
	currSymbol := model.GetCurrencySymbol(currency)

	d.setSubtleNormalText()

	err := d.itemTableRow("Description", "Quantity", "Rate", "Amount")
	if err != nil {
		return err
	}

	d.setNormalText()

	subtotal := 0.

	for _, item := range items {
		subtotal += item.GetAmount()

		err := d.itemTableRow(
			item.Description,
			strconv.Itoa(item.Quantity),
			strconv.FormatFloat(item.Rate, 'f', 2, 64)+currSymbol,
			strconv.FormatFloat(item.GetAmount(), 'f', 2, 64)+currSymbol,
		)
		if err != nil {
			return err
		}
	}

	total := subtotal + tax.GetVat(subtotal) - tax.GetRetention(subtotal)

	d.po.Br(LineHeight)
	startY := d.po.GetY()
	d.setSubtleNormalText()

	err = d.po.CellWithOption(&gopdf.Rect{W: ItemQtyWidth}, "Payment Info", d.getCellOptions(gopdf.Left))
	if err != nil {
		return err
	}

	d.po.Br(LineHeight)

	d.po.SetStrokeColor(colorGray())
	d.po.SetFillColor(colorGray())

	err = d.po.Rectangle(
		Margin,
		d.po.GetY(),
		ItemDescWidth,
		d.po.GetY()+3*LineHeight,
		"DF",
		0., //nolint:gomnd //static value
		0,
	)
	if err != nil {
		return err
	}

	d.po.Br(5)            //nolint:gomnd //static value
	d.po.SetX(Margin + 5) //nolint:gomnd //static value
	d.setNormalText()

	err = d.po.CellWithOption(&gopdf.Rect{W: ItemQtyWidth}, "Holder: "+payment.Holder, d.getCellOptions(gopdf.Left))
	if err != nil {
		return err
	}

	d.po.Br(FromToLineHeight)
	d.po.SetX(Margin + 5) //nolint:gomnd //static value

	err = d.po.CellWithOption(&gopdf.Rect{W: ItemQtyWidth}, "IBAN: "+payment.Iban, d.getCellOptions(gopdf.Left))
	if err != nil {
		return err
	}

	d.po.Br(FromToLineHeight)
	d.po.SetX(Margin + 5) //nolint:gomnd //static value

	err = d.po.CellWithOption(&gopdf.Rect{W: ItemQtyWidth}, "Swift: "+payment.Swift, d.getCellOptions(gopdf.Left))
	if err != nil {
		return err
	}

	d.po.Br(LineHeight)

	if d.po.GetY() > d.lastYPos {
		d.lastYPos = d.po.GetY()
	}

	d.po.SetY(startY)
	d.po.SetStrokeColor(colorBlue())
	d.po.Br(LineHeight)
	d.po.Line(Margin+ItemDescWidth, d.po.GetY(), gopdf.PageSizeA4.W-Margin, d.po.GetY())
	d.po.Br(LineHeight / 2) //nolint:gomnd //static value

	err = d.itemTableRow(
		"",
		"Subtotal",
		"",
		strconv.FormatFloat(subtotal, 'f', 2, 64)+currSymbol)
	if err != nil {
		return err
	}

	mark := "*"
	vatLabel := "VAT"
	retentionLabel := "IRPF"

	if tax.Vat == 0 {
		vatLabel += mark
		mark += "*"
	}

	if tax.Retention != 0 {
		retentionLabel += mark
	}

	err = d.itemTableRow(
		"",
		vatLabel,
		strconv.FormatFloat(tax.Vat, 'f', 0, 64)+"%",
		strconv.FormatFloat(tax.GetVat(subtotal), 'f', 2, 64)+currSymbol)
	if err != nil {
		return err
	}

	if tax.Retention != 0 {
		err = d.itemTableRow(
			"",
			retentionLabel,
			"-"+strconv.FormatFloat(tax.Retention, 'f', 0, 64)+"%",
			strconv.FormatFloat(-tax.GetRetention(subtotal), 'f', 2, 64)+currSymbol)
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

func (d *Document) notes(notes model.Notes) error {
	notesSlice := notes.ToSlice()
	mark := ""

	d.setNormalText()
	d.po.SetY(float64(842 - Margin - len(notesSlice)*20))

	for _, line := range notesSlice {
		err := d.po.MultiCell(
			&gopdf.Rect{W: gopdf.PageSizeA4.W - 2*Margin, H: LineHeight},
			mark+line,
		)
		if err != nil {
			return err
		}

		mark += "*"

		d.po.Br(5) //nolint:gomnd //static value
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

	d.po.Br(LineHeight)

	return nil
}

func (d *Document) headingTitle() error {
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

	d.po.Br(36) //nolint:gomnd //static value

	return nil
}

func (d *Document) headingInfoLine(key, value string) error {
	d.setSubtleNormalText()

	err := d.po.CellWithOption(&gopdf.Rect{W: HeaderInfoName},
		key,
		d.getCellOptions(gopdf.Left),
	)
	if err != nil {
		return err
	}

	err = d.po.CellWithOption(&gopdf.Rect{W: HeaderInfoSeparator},
		":",
		d.getCellOptions(gopdf.Left),
	)
	if err != nil {
		return err
	}

	d.setNormalText()

	err = d.po.CellWithOption(
		&gopdf.Rect{W: HeaderInfoValue},
		value,
		d.getCellOptions(gopdf.Right),
	)
	if err != nil {
		return err
	}

	d.po.Br(LineHeight)

	return nil
}

func (d *Document) from(from *model.Freelancer) error {
	d.setSubtleNormalText()

	err := d.po.Cell(&gopdf.Rect{W: FromWidth}, "From")
	fmt.Printf("invoiceSerice.from: error %v", err)
	d.po.Br(LineHeight)

	d.setNormalText()
	err = d.po.Cell(&gopdf.Rect{W: FromWidth}, from.Name)
	fmt.Printf("invoiceSerice.from: error %v", err)
	d.po.Br(FromToLineHeight)

	if from.Company != "" {
		err = d.po.Cell(&gopdf.Rect{W: FromWidth}, from.Company)
		fmt.Printf("invoiceSerice.from: error %v", err)
		d.po.Br(FromToLineHeight)
	}

	if from.VatID != "" {
		err = d.po.Cell(&gopdf.Rect{W: FromWidth}, from.VatID)
		fmt.Printf("invoiceSerice.from: error %v", err)
		d.po.Br(FromToLineHeight)
	}

	if from.Address1 != "" {
		err = d.po.Cell(&gopdf.Rect{W: FromWidth}, from.Address1)
		fmt.Printf("invoiceSerice.from: error %v", err)
		d.po.Br(FromToLineHeight)
	}

	if from.Address2 != "" {
		err = d.po.Cell(&gopdf.Rect{W: FromWidth}, from.Address2)
		fmt.Printf("invoiceSerice.from: error %v", err)
		d.po.Br(FromToLineHeight)
	}

	if from.Phone != "" {
		err = d.po.Cell(&gopdf.Rect{W: FromWidth}, from.Phone)
		fmt.Printf("invoiceSerice.from: error %v", err)
		d.po.Br(FromToLineHeight)
	}

	d.po.Br(LineHeight)

	endY := d.po.GetY()
	if endY > d.lastYPos {
		d.lastYPos = endY
	}

	return err
}

func (d *Document) to(client *model.Client) error {
	d.po.SetX(ToStart)

	d.setSubtleNormalText()
	err := d.po.Cell(&gopdf.Rect{W: ToWidth}, "To")
	fmt.Printf("invoiceSerice.to: error %v", err)
	d.po.Br(LineHeight)

	d.po.SetX(ToStart)
	d.setNormalText()
	err = d.po.Cell(&gopdf.Rect{W: ToWidth}, client.Name)
	fmt.Printf("invoiceSerice.to: error %v", err)
	d.po.Br(FromToLineHeight)

	d.po.SetX(ToStart)
	err = d.po.Cell(&gopdf.Rect{W: ToWidth}, client.VatID)
	fmt.Printf("invoiceSerice.to: error %v", err)
	d.po.Br(FromToLineHeight)

	if client.Address1 != "" {
		d.po.SetX(ToStart)
		err = d.po.Cell(&gopdf.Rect{W: ToWidth}, client.Address1)
		fmt.Printf("invoiceSerice.to: error %v", err)
		d.po.Br(FromToLineHeight)
	}

	if client.Address2 != "" {
		d.po.SetX(ToStart)
		err = d.po.Cell(&gopdf.Rect{W: ToWidth}, client.Address2)
		fmt.Printf("invoiceSerice.to: error %v", err)
		d.po.Br(FromToLineHeight)
	}

	endY := d.po.GetY()
	if endY > d.lastYPos {
		d.lastYPos = endY
	}

	return err
}

func (d *Document) getCellOptions(align int) gopdf.CellOption {
	co := gopdf.CellOption{Align: align}
	if d.debug {
		co.Border = gopdf.AllBorders
	}

	return co
}

func (d *Document) setNormalText() {
	d.po.SetTextColor(colorBlack())

	err := d.po.SetFont("Inter", "", FontNormalSize)
	if err != nil {
		fmt.Println("Error Loading font: 'Inter'")
	}
}

func (d *Document) setSubtleNormalText() {
	d.po.SetTextColor(colorLavender())

	err := d.po.SetFont("Inter", "", FontSubtleNormalSize)
	if err != nil {
		fmt.Println("Error Loading font: 'Inter'")
	}
}

func (d *Document) setSubtleTotalText() {
	d.po.SetTextColor(colorLavender())

	err := d.po.SetFont("Inter-Bold", "", FontSubtleTotalSize)
	if err != nil {
		fmt.Println("Error Loading font: 'Inter-Bold'")
	}
}

func (d *Document) setTitleText() error {
	d.po.SetTextColor(colorBlack())

	err := d.po.SetFont("Inter-Bold", "", FontTitleSize)
	if err != nil {
		return err
	}

	return nil
}

func colorBlue() (red, green, blue uint8) {
	return 0, 0, 200 //nolint:gomnd // static value for color schema
}

func colorGray() (red, green, blue uint8) {
	return 192, 192, 192 //nolint:gomnd // static value for color schema
}

func colorBlack() (red, green, blue uint8) {
	return 24, 24, 24 //nolint:gomnd // static value for color schema
}

func colorLavender() (red, green, blue uint8) {
	return 128, 128, 192 //nolint:gomnd // static value for color schema
}

func getImageScaledDimension(imagePath string) (scaledWidth, scaledHeight float64) {
	file, err := os.Open(imagePath)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	defer file.Close()

	img, _, err := image.DecodeConfig(file)
	if err != nil {
		fmt.Printf("%s: %v\n", imagePath, err)
	}

	scaledHeight = HeaderLogoSize
	scaledWidth = float64(img.Width) * scaledHeight / float64(img.Height)

	return
}
