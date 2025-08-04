# Internationalization (i18n) Support

Invoiceling now supports multiple languages for PDF generation. Language support is implemented at the rendering level only - your invoice data remains language-neutral while the PDF output can be generated in different languages.

## Supported Languages

- **English** (`en`) - Default language
- **Spanish** (`es`) - Full Spanish translation including tax terminology

## Usage

### Generating PDFs with Language Support

Use the `--language` or `-l` flag when generating PDFs:

```bash
# Generate PDF in English (default)
./invoiceling pdf --invoice F24-001

# Generate PDF in English (explicit)
./invoiceling pdf --invoice F24-001 --language en

# Generate PDF in Spanish
./invoiceling pdf --invoice F24-001 --language es

# Generate draft PDF in Spanish
./invoiceling pdf --invoice F24-001 --language es --draft
```

## What Gets Translated

The following elements are translated in the PDF:

### Header Section

- "INVOICE" / "FACTURA"
- "Invoice" / "Factura"
- "Date" / "Fecha"
- "Due" / "Vence"

### Address Section

- "From" / "De"
- "To" / "Cliente"

### Item Table

- "Description" / "Descripción"
- "Quantity" / "Cantidad"
- "Rate" / "Precio"
- "Amount" / "Subtotal"

### Payment Information

- "Payment Info" / "Información de Pago"
- "Holder:" / "Titular:"
- "IBAN:" (same in both languages)
- "Swift:" (same in both languages)

### Totals Section

- "Subtotal" (same in both languages)
- "VAT" / "IVA"
- "IRPF" (same in both languages - Spanish tax retention)
- "Total" (same in both languages)

### Status

- "DRAFT" / "BORRADOR"

## Examples

### English Invoice

```bash
./invoiceling pdf -i F24-001 -l en
```

This generates a PDF with English labels like "INVOICE", "From", "To", "Description", "VAT", etc.

### Spanish Invoice

```bash
./invoiceling pdf -i F24-001 -l es
```

This generates a PDF with Spanish labels like "FACTURA", "De", "Para", "Descripción", "IVA", etc.

### Spanish Draft Invoice

```bash
./invoiceling pdf -i F24-001 -l es -d
```

This generates a draft PDF in Spanish with "BORRADOR" watermark.

## Technical Details

- Language is specified only at PDF generation time
- Invoice data files (JSON) remain language-neutral
- Default language is English if not specified
- Invalid language codes fall back to English
- Language setting doesn't affect the invoice data, only the PDF rendering

## Error Handling

If you specify an unsupported language:

```bash
./invoiceling pdf -i F24-001 -l fr
```

The system will return an error: `unsupported language: fr`

## Extending Language Support

To add more languages, developers can extend the translations in `pkg/i18n/i18n.go` by adding new language constants and translation maps.

Currently supported language codes follow ISO 639-1 standard:

- `en` - English
- `es` - Spanish (Español)
