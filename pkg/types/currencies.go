package types

type Currency string

const (
	CurrencyUSD Currency = "USD"
	CurrencyEUR Currency = "EUR"
	CurrencyGBP Currency = "GBP"
	CurrencyAED Currency = "AED"
	CurrencyBHD Currency = "BHD"
	CurrencySAR Currency = "SAR"
)

func GetCurrencies() []Currency {
	return []Currency{
		CurrencyUSD,
		CurrencyEUR,
		CurrencyGBP,
		CurrencyAED,
		CurrencyBHD,
		CurrencySAR,
	}
}
