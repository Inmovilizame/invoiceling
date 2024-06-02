package render

import (
	"fmt"
	"image"
	"os"
	"strconv"
	"time"

	"github.com/Inmovilizame/invoiceling/assets"
	"github.com/Inmovilizame/invoiceling/pkg/model"
	"github.com/signintech/gopdf"
)

type DateFormat string

const (
	DFText DateFormat = "Jan 02, 2006"
	DFYMD  DateFormat = "2006-01-02"
	DFDMY  DateFormat = "02/01/2006"
)

const (
	FontSizeNormal       = 10
	FontSizeSubtleNormal = 12
	FontSizeSubtleTotal  = 14
	FontSizeTitle        = 24
	FontSizeDraftMark    = 92
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

const (
	DraftText            = "DRAFT"
	DraftAlpha           = 0.65
	DraftVerticalShift   = 200
	DraftHorizontalShift = 350
)

type PdfBasic struct {
	debug    bool
	lastYPos float64
	gopdf.GoPdf
}

func NewPdfBasicRender() (*PdfBasic, error) {
	interFont, err := assets.FS.ReadFile("fonts/Inter.ttf")
	if err != nil {
		return nil, err
	}

	interBoldFont, err := assets.FS.ReadFile("fonts/Inter-Bold.ttf")
	if err != nil {
		return nil, err
	}

	pb := PdfBasic{}
	pb.Start(gopdf.Config{
		PageSize: *gopdf.PageSizeA4,
	})
	pb.SetMargins(Margin, Margin, Margin, Margin)
	pb.AddPage()

	err = pb.AddTTFFontData("Inter", interFont)
	if err != nil {
		return nil, err
	}

	err = pb.AddTTFFontData("Inter-Bold", interBoldFont)
	if err != nil {
		return nil, err
	}

	return &pb, nil
}

func (p *PdfBasic) Render(invoice *model.Invoice, draft bool) error {
	err := p.header(invoice.Logo, invoice.ID, invoice.Date, invoice.Due)
	if err != nil {
		return err
	}

	p.SetY(p.lastYPos)
	p.SetStrokeColor(colorBlue())
	p.Line(Margin, p.GetY(), gopdf.PageSizeA4.W-Margin, p.GetY())
	p.Br(LineHeight)

	err = p.sendingInfo(&invoice.From, &invoice.To)
	if err != nil {
		return err
	}

	p.SetY(p.lastYPos)
	p.SetStrokeColor(colorBlue())
	p.Line(Margin, p.GetY(), gopdf.PageSizeA4.W-Margin, p.GetY())
	p.Br(LineHeight)

	err = p.items(invoice.Items, invoice.Tax, invoice.Currency, invoice.Payment)
	if err != nil {
		return err
	}

	p.SetY(p.lastYPos)
	p.Br(LineHeight)

	err = p.notes(invoice.Notes)
	if err != nil {
		return err
	}

	if draft {
		err = p.draftOverlay()
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *PdfBasic) SaveTo(filename string) error {
	return p.WritePdf(filename)
}

func (p *PdfBasic) header(logo, id string, date time.Time, due time.Duration) error {
	err := p.headingTitle()
	if err != nil {
		return err
	}

	err = p.headingInfoLine("Invoice", id)
	if err != nil {
		return err
	}

	err = p.headingInfoLine("Date", date.Format(string(DFYMD)))
	if err != nil {
		return err
	}

	err = p.headingInfoLine("Due", date.Add(due).Format(string(DFYMD)))
	if err != nil {
		return err
	}

	if p.GetY() > p.lastYPos {
		p.lastYPos = p.GetY()
	}

	p.SetX(HeaderInfoStartX)
	p.SetY(Margin)

	if logo != "" {
		startX := p.GetX()
		startY := p.GetY()
		width, height := getImageScaledDimension(logo)

		err := p.Image(logo, startX, startY, &gopdf.Rect{W: width, H: height})
		if err != nil {
			return err
		}

		p.Br(height + LineHeight)

		if p.GetY() > p.lastYPos {
			p.lastYPos = p.GetY()
		}
	}

	if p.lastYPos < HeaderMinHeight {
		p.lastYPos = HeaderMinHeight
	}

	return nil
}

func (p *PdfBasic) sendingInfo(from *model.Freelancer, client *model.Client) error {
	startY := p.GetY()

	err := p.from(from)
	if err != nil {
		return nil
	}

	p.SetY(startY)

	err = p.to(client)
	if err != nil {
		return nil
	}

	return nil
}

//nolint:funlen //TODO fix func length
func (p *PdfBasic) items(items []*model.Item, tax model.TaxInfo, currency string, payment model.Payment) error {
	currSymbol := model.GetCurrencySymbol(currency)

	p.setSubtleNormalText()

	err := p.itemTableRow("Description", "Quantity", "Rate", "Amount")
	if err != nil {
		return err
	}

	p.setNormalText()

	subtotal := 0.

	for _, item := range items {
		subtotal += item.GetAmount()

		err := p.itemTableRow(
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

	p.Br(LineHeight)
	startY := p.GetY()
	p.setSubtleNormalText()

	err = p.CellWithOption(&gopdf.Rect{W: ItemQtyWidth}, "Payment Info", p.getCellOptions(gopdf.Left))
	if err != nil {
		return err
	}

	p.Br(LineHeight)

	p.SetStrokeColor(colorGray())
	p.SetFillColor(colorGray())

	err = p.Rectangle(
		Margin,
		p.GetY(),
		ItemDescWidth,
		p.GetY()+3*LineHeight,
		"DF",
		0., //nolint:gomnd //static value
		0,
	)
	if err != nil {
		return err
	}

	p.Br(5)            //nolint:gomnd //static value
	p.SetX(Margin + 5) //nolint:gomnd //static value
	p.setNormalText()

	err = p.CellWithOption(&gopdf.Rect{W: ItemQtyWidth}, "Holder: "+payment.Holder, p.getCellOptions(gopdf.Left))
	if err != nil {
		return err
	}

	p.Br(FromToLineHeight)
	p.SetX(Margin + 5) //nolint:gomnd //static value

	err = p.CellWithOption(&gopdf.Rect{W: ItemQtyWidth}, "IBAN: "+payment.Iban, p.getCellOptions(gopdf.Left))
	if err != nil {
		return err
	}

	p.Br(FromToLineHeight)
	p.SetX(Margin + 5) //nolint:gomnd //static value

	err = p.CellWithOption(&gopdf.Rect{W: ItemQtyWidth}, "Swift: "+payment.Swift, p.getCellOptions(gopdf.Left))
	if err != nil {
		return err
	}

	p.Br(LineHeight)

	if p.GetY() > p.lastYPos {
		p.lastYPos = p.GetY()
	}

	p.SetY(startY)
	p.SetStrokeColor(colorBlue())
	p.Br(LineHeight)
	p.Line(Margin+ItemDescWidth, p.GetY(), gopdf.PageSizeA4.W-Margin, p.GetY())
	p.Br(LineHeight / 2) //nolint:gomnd //static value

	err = p.itemTableRow(
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

	err = p.itemTableRow(
		"",
		vatLabel,
		strconv.FormatFloat(tax.Vat, 'f', 0, 64)+"%",
		strconv.FormatFloat(tax.GetVat(subtotal), 'f', 2, 64)+currSymbol)
	if err != nil {
		return err
	}

	if tax.Retention != 0 {
		err = p.itemTableRow(
			"",
			retentionLabel,
			"-"+strconv.FormatFloat(tax.Retention, 'f', 0, 64)+"%",
			strconv.FormatFloat(-tax.GetRetention(subtotal), 'f', 2, 64)+currSymbol)
		if err != nil {
			return err
		}
	}

	p.setSubtleTotalText()

	err = p.itemTableRow(
		"",
		"Total",
		"",
		strconv.FormatFloat(total, 'f', 2, 64)+currSymbol)
	if err != nil {
		return err
	}

	if p.GetY() > p.lastYPos {
		p.lastYPos = p.GetY()
	}

	return nil
}

func (p *PdfBasic) notes(notes model.Notes) error {
	notesSlice := notes.ToSlice()
	mark := ""

	p.setNormalText()
	p.SetY(float64(842 - Margin - len(notesSlice)*20))

	for _, line := range notesSlice {
		err := p.MultiCell(
			&gopdf.Rect{W: gopdf.PageSizeA4.W - 2*Margin, H: LineHeight},
			mark+line,
		)
		if err != nil {
			return err
		}

		mark += "*"

		p.Br(5) //nolint:gomnd //static value
	}

	return nil
}

func (p *PdfBasic) draftOverlay() error {
	p.setDraftText()

	for i := 0; i < 4; i++ {
		p.SetX(Margin)
		p.SetY(Margin + float64(i*DraftVerticalShift))

		if i%2 == 1 {
			p.SetX(gopdf.PageSizeA4.W - DraftHorizontalShift)
		}

		err := p.CellWithOption(nil, DraftText, gopdf.CellOption{
			Transparency: &gopdf.Transparency{Alpha: DraftAlpha, BlendModeType: gopdf.Overlay},
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *PdfBasic) itemTableRow(desc, qty, rate, total string) error {
	err := p.CellWithOption(&gopdf.Rect{W: ItemDescWidth}, desc, p.getCellOptions(gopdf.Left))
	if err != nil {
		return err
	}

	err = p.CellWithOption(&gopdf.Rect{W: ItemQtyWidth}, qty, p.getCellOptions(gopdf.Right))
	if err != nil {
		return err
	}

	err = p.CellWithOption(&gopdf.Rect{W: ItemRateWidth}, rate, p.getCellOptions(gopdf.Right))
	if err != nil {
		return err
	}

	err = p.CellWithOption(&gopdf.Rect{W: ItemAmountWidth}, total, p.getCellOptions(gopdf.Right))
	if err != nil {
		return err
	}

	p.Br(LineHeight)

	return nil
}

func (p *PdfBasic) headingTitle() error {
	p.setTitleText()

	err := p.CellWithOption(
		&gopdf.Rect{W: HeaderInfoWidth},
		"INVOICE",
		p.getCellOptions(gopdf.Center),
	)
	if err != nil {
		return err
	}

	p.Br(36) //nolint:gomnd //static value

	return nil
}

func (p *PdfBasic) headingInfoLine(key, value string) error {
	p.setSubtleNormalText()

	err := p.CellWithOption(&gopdf.Rect{W: HeaderInfoName},
		key,
		p.getCellOptions(gopdf.Left),
	)
	if err != nil {
		return err
	}

	err = p.CellWithOption(&gopdf.Rect{W: HeaderInfoSeparator},
		":",
		p.getCellOptions(gopdf.Left),
	)
	if err != nil {
		return err
	}

	p.setNormalText()

	err = p.CellWithOption(
		&gopdf.Rect{W: HeaderInfoValue},
		value,
		p.getCellOptions(gopdf.Right),
	)
	if err != nil {
		return err
	}

	p.Br(LineHeight)

	return nil
}

func (p *PdfBasic) from(from *model.Freelancer) error {
	p.setSubtleNormalText()

	err := p.Cell(&gopdf.Rect{W: FromWidth}, "From")
	if err != nil {
		fmt.Printf("invoiceService.from: error %v", err)
	}

	p.Br(LineHeight)
	p.setNormalText()

	err = p.Cell(&gopdf.Rect{W: FromWidth}, from.Name)
	if err != nil {
		fmt.Printf("invoiceService.from: error %v", err)
	}

	p.Br(FromToLineHeight)

	if from.Company != "" {
		err = p.Cell(&gopdf.Rect{W: FromWidth}, from.Company)
		if err != nil {
			fmt.Printf("invoiceService.from: error %v", err)
		}

		p.Br(FromToLineHeight)
	}

	if from.VatID != "" {
		err = p.Cell(&gopdf.Rect{W: FromWidth}, from.VatID)
		if err != nil {
			fmt.Printf("invoiceService.from: error %v", err)
		}

		p.Br(FromToLineHeight)
	}

	if from.Address1 != "" {
		err = p.Cell(&gopdf.Rect{W: FromWidth}, from.Address1)
		if err != nil {
			fmt.Printf("invoiceService.from: error %v", err)
		}

		p.Br(FromToLineHeight)
	}

	if from.Address2 != "" {
		err = p.Cell(&gopdf.Rect{W: FromWidth}, from.Address2)
		if err != nil {
			fmt.Printf("invoiceService.from: error %v", err)
		}

		p.Br(FromToLineHeight)
	}

	if from.Phone != "" {
		err = p.Cell(&gopdf.Rect{W: FromWidth}, from.Phone)
		if err != nil {
			fmt.Printf("invoiceService.from: error %v", err)
		}

		p.Br(FromToLineHeight)
	}

	p.Br(LineHeight)

	endY := p.GetY()
	if endY > p.lastYPos {
		p.lastYPos = endY
	}

	return err
}

func (p *PdfBasic) to(client *model.Client) error {
	p.SetX(ToStart)
	p.setSubtleNormalText()

	err := p.Cell(&gopdf.Rect{W: ToWidth}, "To")
	if err != nil {
		fmt.Printf("invoiceService.to: error %v", err)
	}

	p.Br(LineHeight)
	p.SetX(ToStart)
	p.setNormalText()

	err = p.Cell(&gopdf.Rect{W: ToWidth}, client.Name)
	if err != nil {
		fmt.Printf("invoiceService.to: error %v", err)
	}

	p.Br(FromToLineHeight)
	p.SetX(ToStart)

	err = p.Cell(&gopdf.Rect{W: ToWidth}, client.VatID)
	if err != nil {
		fmt.Printf("invoiceService.to: error %v", err)
	}

	p.Br(FromToLineHeight)

	if client.Address1 != "" {
		p.SetX(ToStart)

		err = p.Cell(&gopdf.Rect{W: ToWidth}, client.Address1)
		if err != nil {
			fmt.Printf("invoiceService.to: error %v", err)
		}

		p.Br(FromToLineHeight)
	}

	if client.Address2 != "" {
		p.SetX(ToStart)

		err = p.Cell(&gopdf.Rect{W: ToWidth}, client.Address2)
		if err != nil {
			fmt.Printf("invoiceService.to: error %v", err)
		}

		p.Br(FromToLineHeight)
	}

	endY := p.GetY()
	if endY > p.lastYPos {
		p.lastYPos = endY
	}

	return err
}

func (p *PdfBasic) getCellOptions(align int) gopdf.CellOption {
	co := gopdf.CellOption{Align: align}
	if p.debug {
		co.Border = gopdf.AllBorders
	}

	return co
}

func (p *PdfBasic) setNormalText() {
	p.SetTextColor(colorBlack())

	err := p.SetFont("Inter", "", FontSizeNormal)
	if err != nil {
		fmt.Println("Error Loading font: 'Inter'")
	}
}

func (p *PdfBasic) setSubtleNormalText() {
	p.SetTextColor(colorLavender())

	err := p.SetFont("Inter", "", FontSizeSubtleNormal)
	if err != nil {
		fmt.Println("Error Loading font: 'Inter'")
	}
}

func (p *PdfBasic) setSubtleTotalText() {
	p.SetTextColor(colorLavender())

	err := p.SetFont("Inter-Bold", "", FontSizeSubtleTotal)
	if err != nil {
		fmt.Println("Error Loading font: 'Inter-Bold'")
	}
}

func (p *PdfBasic) setTitleText() {
	p.SetTextColor(colorBlack())

	err := p.SetFont("Inter-Bold", "", FontSizeTitle)
	if err != nil {
		fmt.Println("Error Loading font: 'Inter-Bold'")
	}
}

func (p *PdfBasic) setDraftText() {
	p.SetTextColor(colorGray())

	err := p.SetFont("Inter-Bold", "", FontSizeDraftMark)
	if err != nil {
		fmt.Println("Error Loading font: 'Inter-Bold'")
	}
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
