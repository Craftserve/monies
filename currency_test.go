package monies_test

import (
	"github.com/Craftserve/monies"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAddCurrency(t *testing.T) {
	var code monies.CurrencyCode = "TEST"
	decimals := 5
	monies.AddCurrency(code, "T$", "1 $", ".", ",", decimals, "0")

	m, err := monies.New(1, code)
	assert.NoError(t, err)

	assert.Equal(t, code, m.Currency().Code)
	assert.Equal(t, decimals, m.Currency().Fraction)
}

func TestCurrencyGetCurrency(t *testing.T) {
	var code monies.CurrencyCode = "KLINGONDOLLAR"
	desired := monies.Currency{Decimal: ".", Thousand: ",", Code: code, Fraction: 2, Grapheme: "$", Template: "$1"}
	monies.AddCurrency(desired.Code, desired.Grapheme, desired.Template, desired.Decimal, desired.Thousand, desired.Fraction, desired.NumericCode)

	currency, err := monies.GetCurrency(code)
	require.NoError(t, err)
	assert.Equal(t, desired, currency)

}

func TestCurrencyGetNonExistingCurrency(t *testing.T) {
	_, err := monies.GetCurrency("I*am*Not*a*CurrencyCode")
	assert.Error(t, err, monies.ErrCurrencyNotFound)
}

func TestCurrencyGetCurrencyByNumericCode(t *testing.T) {
	var code monies.CurrencyCode = "EUROGĄBKI"
	desired := monies.Currency{Decimal: ".", Thousand: ",", Code: code, Fraction: 2, Grapheme: "$", Template: "$1", NumericCode: "9999"}
	monies.AddCurrency(desired.Code, desired.Grapheme, desired.Template, desired.Decimal, desired.Thousand, desired.Fraction, desired.NumericCode)

	currency, err := monies.Currencies.CurrencyByNumericCode("9999")
	assert.NoError(t, err)
	assert.Equal(t, desired, currency)

}

func TestCurrencyCurrencyByNumericCodeNonExisting(t *testing.T) {
	_, err := monies.Currencies.CurrencyByNumericCode("0900990")
	assert.ErrorIs(t, err, monies.ErrCurrencyNotFound)
}
