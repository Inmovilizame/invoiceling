package i18n

import "fmt"

// Language represents a supported language
type Language string

const (
	English Language = "en"
	Spanish Language = "es"
)

// Translator provides translation functionality
type Translator interface {
	T(key string) string
	SetLanguage(lang Language)
	GetLanguage() Language
}

// translator implements the Translator interface
type translator struct {
	language     Language
	translations map[Language]map[string]string
}

// NewTranslator creates a new translator with default language English
func NewTranslator() Translator {
	t := &translator{
		language:     English,
		translations: make(map[Language]map[string]string),
	}
	t.loadTranslations()

	return t
}

// T translates a key to the current language
// TODO: better defaults handling, why handle language at translation time?
func (t *translator) T(key string) string {
	if translations, exists := t.translations[t.language]; exists {
		if translation, exists := translations[key]; exists {
			return translation
		}
	}

	// Fallback to English if translation not found
	if t.language != English {
		if translations, exists := t.translations[English]; exists {
			if translation, exists := translations[key]; exists {
				return translation
			}
		}
	}

	// Return key if no translation found
	return key
}

// SetLanguage sets the current language
func (t *translator) SetLanguage(lang Language) {
	t.language = lang
}

// GetLanguage returns the current language
func (t *translator) GetLanguage() Language {
	return t.language
}

// loadTranslations loads all translation data
func (t *translator) loadTranslations() {
	t.translations[English] = map[string]string{
		// PDF Header
		"invoice":      "Invoice",
		"invoice_caps": "INVOICE",
		"date":         "Date",
		"due":          "Due",

		// PDF Sections
		"from": "From",
		"to":   "To",

		// Table Headers
		"description": "Description",
		"quantity":    "Quantity",
		"rate":        "Rate",
		"amount":      "Amount",

		// Payment Section
		"payment_info": "Payment Info",
		"holder":       "Holder",
		"iban":         "IBAN",
		"swift":        "Swift",

		// Totals
		"subtotal": "Subtotal",
		"vat":      "VAT",
		"irpf":     "IRPF",
		"total":    "Total",

		// Status
		"draft": "DRAFT",

		// Payment Info Labels
		"holder_label": "Holder: ",
		"iban_label":   "IBAN: ",
		"swift_label":  "Swift: ",
	}

	t.translations[Spanish] = map[string]string{
		// PDF Header
		"invoice":      "Factura",
		"invoice_caps": "FACTURA",
		"date":         "Fecha",
		"due":          "Vence",

		// PDF Sections
		"from": "De",
		"to":   "Para",

		// Table Headers
		"description": "Descripción",
		"quantity":    "Cantidad",
		"rate":        "Precio",
		"amount":      "Importe",

		// Payment Section
		"payment_info": "Información de Pago",
		"holder":       "Titular",
		"iban":         "IBAN",
		"swift":        "Swift",

		// Totals
		"subtotal": "Subtotal",
		"vat":      "IVA",
		"irpf":     "IRPF",
		"total":    "Total",

		// Status
		"draft": "BORRADOR",

		// Payment Info Labels
		"holder_label": "Titular: ",
		"iban_label":   "IBAN: ",
		"swift_label":  "Swift: ",
	}
}

// GetSupportedLanguages returns a list of supported languages
func GetSupportedLanguages() []Language {
	return []Language{English, Spanish}
}

// IsLanguageSupported checks if a language is supported
// TODO: rely on slices.Contains?
func IsLanguageSupported(lang Language) bool {
	for _, supportedLang := range GetSupportedLanguages() {
		if lang == supportedLang {
			return true
		}
	}

	return false
}

// ParseLanguage parses a string to Language type with validation
func ParseLanguage(lang string) (Language, error) {
	parsedLang := Language(lang)
	if !IsLanguageSupported(parsedLang) {
		return English, fmt.Errorf("unsupported language: %s", lang)
	}

	return parsedLang, nil
}
