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
	strs := make([]string, 0, len(c))
	for _, currency := range c {
		strs = append(strs, currency.String())
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
