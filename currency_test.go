package money_test

import (
	"github.com/Craftserve/monies"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAddCurrency(t *testing.T) {
	var code money.CurrencyCode = "TEST"
	decimals := 5
	money.AddCurrency(code, "T$", "1 $", ".", ",", decimals, "0")

	m, err := money.New(1, code)
	assert.NoError(t, err)

	assert.Equal(t, code, m.Currency().Code)
	assert.Equal(t, decimals, m.Currency().Fraction)
}

func TestCurrencyGetCurrency(t *testing.T) {
	var code money.CurrencyCode = "KLINGONDOLLAR"
	desired := money.Currency{Decimal: ".", Thousand: ",", Code: code, Fraction: 2, Grapheme: "$", Template: "$1"}
	money.AddCurrency(desired.Code, desired.Grapheme, desired.Template, desired.Decimal, desired.Thousand, desired.Fraction, desired.NumericCode)

	currency, err := money.GetCurrency(code)
	require.NoError(t, err)
	assert.Equal(t, desired, currency)

}

func TestCurrencyGetNonExistingCurrency(t *testing.T) {
	_, err := money.GetCurrency("I*am*Not*a*CurrencyCode")
	assert.Error(t, err, money.ErrCurrencyNotFound)
}

func TestCurrencyGetCurrencyByNumericCode(t *testing.T) {
	var code money.CurrencyCode = "EUROGĄBKI"
	desired := money.Currency{Decimal: ".", Thousand: ",", Code: code, Fraction: 2, Grapheme: "$", Template: "$1", NumericCode: "9999"}
	money.AddCurrency(desired.Code, desired.Grapheme, desired.Template, desired.Decimal, desired.Thousand, desired.Fraction, desired.NumericCode)

	currency, err := money.Currencies.CurrencyByNumericCode("9999")
	assert.NoError(t, err)
	assert.Equal(t, desired, currency)

}

func TestCurrencyCurrencyByNumericCodeNonExisting(t *testing.T) {
	_, err := money.Currencies.CurrencyByNumericCode("0900990")
	assert.ErrorIs(t, err, money.ErrCurrencyNotFound)
}
