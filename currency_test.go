package monies_test

import (
	"github.com/Craftserve/monies"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCurrencyGetNonExistingCurrency(t *testing.T) {
	_, err := monies.CurrencyByCode("I*am*Not*a*CurrencyCode")
	assert.Error(t, err, monies.ErrCurrencyNotFound)
}

func TestCurrencyGetCurrencyByNumericCode(t *testing.T) {
	desired := monies.Currency{Decimal: ",", Thousand: ".", Code: monies.HUF, Fraction: 0, NumericCode: "348", Grapheme: "Ft", Template: "1 $"}
	currency, err := monies.CurrencyByNumericCode("348")

	assert.NoError(t, err)
	assert.Equal(t, desired, currency)
}

func TestCurrencyCurrencyByNumericCodeNonExisting(t *testing.T) {
	_, err := monies.CurrencyByNumericCode("0900990")
	assert.ErrorIs(t, err, monies.ErrCurrencyNotFound)
}
