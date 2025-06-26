package types

type Currency string
type Currencies []Currency

const (
	CurrencyUSD Currency = "USD"
	CurrencyEUR Currency = "EUR"
	CurrencyGBP Currency = "GBP"
	CurrencyAED Currency = "AED"
	CurrencyBHD Currency = "BHD"
	CurrencySAR Currency = "SAR"
)

func (c Currency) String() string {
	return string(c)
}

func (c Currencies) String() []string {
	strs := make([]string, len(c))
	for i, currency := range c {
		strs[i] = currency.String()
	}
	return strs
}

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
