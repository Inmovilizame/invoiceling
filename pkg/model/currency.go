package model

var currencySymbols = map[string]string{
	"USD": "$",
	"EUR": "€",
	"GBP": "£",
	"JPY": "¥",
	"CNY": "¥",
	"INR": "₹",
	"RUB": "₽",
	"KRW": "₩",
	"BRL": "R$",
	"SGD": "SGD$",
}

func GetCurrencySymbol(code string) string {
	return currencySymbols[code]
}
