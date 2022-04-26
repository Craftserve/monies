package monies_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Craftserve/monies"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func MustNew(amount int64, code monies.CurrencyCode) monies.Money {
	m, err := monies.New(amount, code)
	if err != nil {
		panic(err)
	}

	return m
}

func TestNew(t *testing.T) {
	type testCase struct {
		Name  string
		Input struct {
			CurrencyCode monies.CurrencyCode
			Amount       int64
		}
		ExpectedErr error
	}

	var testCases = []testCase{
		{
			Name: "SUCCESS",
			Input: struct {
				CurrencyCode monies.CurrencyCode
				Amount       int64
			}{
				CurrencyCode: monies.EUR,
				Amount:       1000,
			},
			ExpectedErr: nil,
		},
		{
			Name: "CURRENCY_NOT_FOUND",
			Input: struct {
				CurrencyCode monies.CurrencyCode
				Amount       int64
			}{
				CurrencyCode: "UNDEFINED_CURRENCY",
				Amount:       1000,
			},
			ExpectedErr: monies.ErrCurrencyNotFound,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.Name, func(t *testing.T) {
			t.Parallel()
			m, err := monies.New(tC.Input.Amount, tC.Input.CurrencyCode)
			assert.ErrorIs(t, tC.ExpectedErr, err)

			if tC.ExpectedErr == nil {
				assert.Equal(t, tC.Input.Amount, m.Amount())
				assert.Equal(t, tC.Input.CurrencyCode, m.Currency().Code)
			}
		})
	}
}

func TestSameCurrency(t *testing.T) {
	m, err := monies.New(0, monies.EUR)
	assert.NoError(t, err)

	other, err := monies.New(0, monies.USD)
	assert.NoError(t, err)

	assert.Equal(t, false, m.SameCurrency(other))

	sameCurrency, err := monies.New(0, monies.EUR)
	assert.NoError(t, err)
	assert.Equal(t, true, m.SameCurrency(sameCurrency))
}

func TestEquals(t *testing.T) {
	testCases := []struct {
		Name        string
		Money       monies.Money
		OtherMoney  monies.Money
		Expected    bool
		ExpectedErr error
	}{
		{"SUCCESS", MustNew(100, monies.EUR), MustNew(100, monies.EUR), true, nil},
		{"FAIL_OTHER_CURRENCIES", MustNew(100, monies.USD), MustNew(100, monies.EUR), false, monies.ErrCurrencyMismatch},
		{"FAIL_OTHER_AMOUNTS", MustNew(50, monies.EUR), MustNew(100, monies.EUR), false, nil},
	}

	for _, tC := range testCases {
		t.Run(tC.Name, func(t *testing.T) {
			equals, err := tC.Money.Equals(tC.OtherMoney)
			if tC.ExpectedErr == nil {
				assert.Equal(t, tC.Expected, equals)
			}

			assert.ErrorIs(t, tC.ExpectedErr, err)
		})
	}
}

func TestGreaterThan(t *testing.T) {
	testCases := []struct {
		Name        string
		Money       monies.Money
		OtherMoney  monies.Money
		Expected    bool
		ExpectedErr error
	}{
		{"SUCCESS", MustNew(1000, monies.EUR), MustNew(100, monies.EUR), true, nil},
		{"FAIL_OTHER_CURRENCIES", MustNew(100, monies.USD), MustNew(100, monies.EUR), false, monies.ErrCurrencyMismatch},
		{"FAIL_OTHER_AMOUNTS", MustNew(50, monies.EUR), MustNew(100, monies.EUR), false, nil},
	}

	for _, tC := range testCases {
		t.Run(tC.Name, func(t *testing.T) {
			equals, err := tC.Money.GreaterThan(tC.OtherMoney)
			if tC.ExpectedErr == nil {
				assert.Equal(t, tC.Expected, equals)
			}

			assert.ErrorIs(t, tC.ExpectedErr, err)
		})
	}
}
func TestGreaterThanOrEqual(t *testing.T) {
	testCases := []struct {
		Name        string
		Money       monies.Money
		OtherMoney  monies.Money
		Expected    bool
		ExpectedErr error
	}{
		{"SUCCESS_GREATER", MustNew(1000, monies.EUR), MustNew(100, monies.EUR), true, nil},
		{"SUCCESS_EQUAL", MustNew(100, monies.EUR), MustNew(100, monies.EUR), true, nil},
		{"FAIL_OTHER_CURRENCIES", MustNew(100, monies.USD), MustNew(100, monies.EUR), false, monies.ErrCurrencyMismatch},
		{"FAIL_OTHER_AMOUNTS", MustNew(50, monies.EUR), MustNew(100, monies.EUR), false, nil},
	}

	for _, tC := range testCases {
		t.Run(tC.Name, func(t *testing.T) {
			equals, err := tC.Money.GreaterThanOrEqual(tC.OtherMoney)
			if tC.ExpectedErr == nil {
				assert.Equal(t, tC.Expected, equals)
			}

			assert.ErrorIs(t, tC.ExpectedErr, err)
		})
	}
}

func TestLessThan(t *testing.T) {
	testCases := []struct {
		Name        string
		Money       monies.Money
		OtherMoney  monies.Money
		Expected    bool
		ExpectedErr error
	}{
		{"SUCCESS", MustNew(50, monies.EUR), MustNew(100, monies.EUR), true, nil},
		{"FAIL_OTHER_CURRENCIES", MustNew(100, monies.USD), MustNew(100, monies.EUR), false, monies.ErrCurrencyMismatch},
		{"FAIL_OTHER_AMOUNTS", MustNew(150, monies.EUR), MustNew(100, monies.EUR), false, nil},
	}

	for _, tC := range testCases {
		t.Run(tC.Name, func(t *testing.T) {
			equals, err := tC.Money.LessThan(tC.OtherMoney)
			if tC.ExpectedErr == nil {
				assert.Equal(t, tC.Expected, equals)
			}

			assert.ErrorIs(t, tC.ExpectedErr, err)
		})
	}
}

func TestLessThanOrEqual(t *testing.T) {
	testCases := []struct {
		Name        string
		Money       monies.Money
		OtherMoney  monies.Money
		Expected    bool
		ExpectedErr error
	}{
		{"SUCCESS_LESS", MustNew(50, monies.EUR), MustNew(100, monies.EUR), true, nil},
		{"SUCCESS_EQUAL", MustNew(50, monies.EUR), MustNew(50, monies.EUR), true, nil},
		{"FAIL_OTHER_CURRENCIES", MustNew(100, monies.USD), MustNew(100, monies.EUR), false, monies.ErrCurrencyMismatch},
		{"FAIL_OTHER_AMOUNTS", MustNew(150, monies.EUR), MustNew(100, monies.EUR), false, nil},
	}

	for _, tC := range testCases {
		t.Run(tC.Name, func(t *testing.T) {
			equals, err := tC.Money.LessThanOrEqual(tC.OtherMoney)
			if tC.ExpectedErr == nil {
				assert.Equal(t, tC.Expected, equals)
			}

			assert.ErrorIs(t, tC.ExpectedErr, err)
		})
	}
}
func TestIsZero(t *testing.T) {
	testCases := []struct {
		Name     string
		Money    monies.Money
		Expected bool
	}{
		{"SUCCESS", MustNew(0, monies.EUR), true},
		{"NOT_ZERO", MustNew(1, monies.USD), false},
	}

	for _, tC := range testCases {
		t.Run(tC.Name, func(t *testing.T) {
			equals := tC.Money.IsZero()
			assert.Equal(t, tC.Expected, equals)
		})
	}
}

func TestIsNegative(t *testing.T) {
	testCases := []struct {
		Name     string
		Money    monies.Money
		Expected bool
	}{
		{"SUCCESS", MustNew(-1, monies.EUR), true},
		{"POSITIVE", MustNew(1, monies.EUR), false},
		{"ZERO", MustNew(0, monies.EUR), false},
	}

	for _, tC := range testCases {
		t.Run(tC.Name, func(t *testing.T) {
			equals := tC.Money.IsNegative()
			assert.Equal(t, tC.Expected, equals)
		})
	}
}

func TestIsPositive(t *testing.T) {
	testCases := []struct {
		Name     string
		Money    monies.Money
		Expected bool
	}{
		{"SUCCESS", MustNew(1, monies.EUR), true},
		{"NEGATIVE", MustNew(-1, monies.EUR), false},
		{"ZERO", MustNew(0, monies.EUR), false},
	}

	for _, tC := range testCases {
		t.Run(tC.Name, func(t *testing.T) {
			equals := tC.Money.IsPositive()
			assert.Equal(t, tC.Expected, equals)
		})
	}
}

func TestAbsolute(t *testing.T) {
	testCases := []struct {
		Name     string
		Money    monies.Money
		Expected int64
	}{
		{"POSITIVE", MustNew(1, monies.EUR), 1},
		{"NEGATIVE", MustNew(-1, monies.EUR), 1},
		{"ZERO", MustNew(0, monies.EUR), 0},
	}

	for _, tC := range testCases {
		t.Run(tC.Name, func(t *testing.T) {

			result := tC.Money.Absolute()
			assert.Equal(t, tC.Expected, result.Amount())
		})
	}
}

func TestNegative(t *testing.T) {
	testCases := []struct {
		Name     string
		Money    monies.Money
		Expected int64
	}{
		{"POSITIVE", MustNew(1, monies.EUR), -1},
		{"NEGATIVE", MustNew(-1, monies.EUR), -1},
		{"ZERO", MustNew(0, monies.EUR), 0},
	}

	for _, tC := range testCases {
		t.Run(tC.Name, func(t *testing.T) {
			result := tC.Money.Negative()
			assert.Equal(t, tC.Expected, result.Amount())
		})
	}
}

func TestAdd(t *testing.T) {
	testCases := []struct {
		Name        string
		Money       monies.Money
		OtherMoney  monies.Money
		Expected    monies.Money
		ExpectedErr error
	}{
		{"POSITIVE_POSITIVE", MustNew(50, monies.EUR), MustNew(100, monies.EUR), MustNew(150, monies.EUR), nil},
		{"POSITIVE_NEGATIVE", MustNew(100, monies.EUR), MustNew(-50, monies.EUR), MustNew(50, monies.EUR), nil},
		{"CURRENCY_MISMATCH", MustNew(100, monies.EUR), MustNew(-50, monies.USD), MustNew(50, monies.USD), monies.ErrCurrencyMismatch},
	}

	for _, tC := range testCases {
		t.Run(tC.Name, func(t *testing.T) {

			result, err := tC.Money.Add(tC.OtherMoney)
			assert.ErrorIs(t, tC.ExpectedErr, err)

			if tC.ExpectedErr == nil {
				assert.NoError(t, err)
				assert.Equal(t, tC.Expected.Amount(), result.Amount())
			}
		})
	}
}

func TestSubstract(t *testing.T) {
	testCases := []struct {
		Name        string
		Money       monies.Money
		OtherMoney  monies.Money
		Expected    monies.Money
		ExpectedErr error
	}{
		{"POSITIVE_POSITIVE", MustNew(100, monies.EUR), MustNew(100, monies.EUR), MustNew(0, monies.EUR), nil},
		{"POSITIVE_NEGATIVE", MustNew(100, monies.EUR), MustNew(-50, monies.EUR), MustNew(150, monies.EUR), nil},
		{"CURRENCY_MISMATCH", MustNew(100, monies.EUR), MustNew(-50, monies.USD), MustNew(50, monies.USD), monies.ErrCurrencyMismatch},
	}

	for _, tC := range testCases {
		t.Run(tC.Name, func(t *testing.T) {
			result, err := tC.Money.Subtract(tC.OtherMoney)
			assert.ErrorIs(t, tC.ExpectedErr, err)

			if tC.ExpectedErr == nil {
				assert.Equal(t, tC.Expected.Amount(), result.Amount())
			}
		})
	}
}
func TestMultiply(t *testing.T) {
	testCases := []struct {
		Name       string
		Money      monies.Money
		Multiplier int64
		Expected   monies.Money
	}{
		{"BY_ONE", MustNew(100, monies.EUR), 1, MustNew(100, monies.EUR)},
		{"BY_ZERO", MustNew(100, monies.EUR), 0, MustNew(0, monies.EUR)},
		{"SUCCESS", MustNew(100, monies.EUR), 2, MustNew(200, monies.EUR)},
	}

	for _, tC := range testCases {
		t.Run(tC.Name, func(t *testing.T) {

			result := tC.Money.Multiply(tC.Multiplier)
			assert.Equal(t, tC.Expected, result)
		})
	}
}

func TestRound(t *testing.T) {
	monies.AddCurrency("TEST_EXPONENTIAL", "*", "$1", ".", ",", 3, "0")

	testCases := []struct {
		Name     string
		Money    monies.Money
		Expected int64
	}{
		{"125_100", MustNew(125, monies.EUR), 100},
		{"175_200", MustNew(175, monies.EUR), 200},
		{"349_300", MustNew(349, monies.EUR), 300},
		{"351_400", MustNew(351, monies.EUR), 400},
		{"0_0", MustNew(0, monies.EUR), 0},
		{"-1_0", MustNew(-1, monies.EUR), 0},
		{"-75_-100", MustNew(-75, monies.EUR), -100},
		{"TEST_EXPONENTIAL", MustNew(12555, monies.CurrencyCode("TEST_EXPONENTIAL")), 13000},
	}

	for _, tC := range testCases {
		t.Run(tC.Name, func(t *testing.T) {

			result := tC.Money.Round()
			assert.Equal(t, tC.Expected, result.Amount())
		})
	}
}

func TestMoneySplit(t *testing.T) {
	testCases := []struct {
		Money       monies.Money
		Split       int
		Expected    []int64
		ExpectedErr error
	}{
		{MustNew(100, monies.EUR), 3, []int64{34, 33, 33}, nil},
		{MustNew(100, monies.EUR), 4, []int64{25, 25, 25, 25}, nil},
		{MustNew(5, monies.EUR), 3, []int64{2, 2, 1}, nil},
		{MustNew(-101, monies.EUR), 4, []int64{-26, -25, -25, -25}, nil},
		{MustNew(-101, monies.EUR), 4, []int64{-26, -25, -25, -25}, nil},
		{MustNew(-2, monies.EUR), 3, []int64{-1, -1, 0}, nil},
		{MustNew(-2, monies.EUR), -1, []int64{-1, -1, 0}, monies.ErrNegativeSplit},
	}

	for index, tC := range testCases {
		t.Run(fmt.Sprintf(`#%d`, index), func(t *testing.T) {
			result, err := tC.Money.Split(tC.Split)
			assert.ErrorIs(t, tC.ExpectedErr, err)

			if tC.ExpectedErr == nil {
				var rs []int64
				for _, party := range result {
					rs = append(rs, party.Amount())
				}

				assert.Equal(t, tC.Expected, rs)
			}
		})
	}
}

func TestMoneyAllocate(t *testing.T) {
	testCases := []struct {
		m           monies.Money
		Ratios      []int
		Expected    []int64
		ExpectedErr error
	}{
		{MustNew(100, monies.EUR), []int{50, 50}, []int64{50, 50}, nil},
		{MustNew(100, monies.EUR), []int{30, 30, 30}, []int64{34, 33, 33}, nil},
		{MustNew(200, monies.EUR), []int{25, 25, 50}, []int64{50, 50, 100}, nil},
		{MustNew(5, monies.EUR), []int{50, 25, 25}, []int64{3, 1, 1}, nil},
		{MustNew(-101, monies.EUR), []int{50, 50}, []int64{-51, -50}, nil},
		{MustNew(-101, monies.EUR), []int{}, []int64{-26, -25}, monies.ErrNoRatios},
	}

	for index, tC := range testCases {
		t.Run(fmt.Sprintf("#%d", index), func(t *testing.T) {
			var rs []int64
			split, err := tC.m.Allocate(tC.Ratios...)

			assert.ErrorIs(t, tC.ExpectedErr, err)

			if tC.ExpectedErr == nil {
				for _, party := range split {
					rs = append(rs, party.Amount())
				}

				assert.Equal(t, tC.Expected, rs)
			}
		})
	}
}

func TestMoneyDisplay(t *testing.T) {
	testCases := []struct {
		m        monies.Money
		expected string
	}{
		{MustNew(100, monies.GBP), "£1.00"},
		{MustNew(100, monies.AED), "1.00 .\u062f.\u0625"},
		{MustNew(-100, monies.GBP), "-£1.00"},
		{MustNew(10, monies.GBP), "£0.10"},
		{MustNew(100000, monies.GBP), "£1,000.00"},
	}

	for _, tC := range testCases {
		display := tC.m.Display()
		assert.Equal(t, tC.expected, display)
	}
}

func TestAsMajorUnits(t *testing.T) {
	monies.AddCurrency("TEST", "T$", "1 $", ".", ",", 0, "0")

	testCases := []struct {
		m        monies.Money
		expected float64
	}{
		{MustNew(100, monies.GBP), 1.00},
		{MustNew(-100, monies.GBP), -1.00},
		{MustNew(0, monies.GBP), 0},
		{MustNew(1, "TEST"), 1},
	}

	for _, tC := range testCases {
		r := tC.m.AsMajorUnits()
		assert.Equal(t, tC.expected, r)
	}
}

func TestCurrency(t *testing.T) {
	pound, err := monies.New(100, monies.GBP)
	require.NoError(t, err)

	assert.Equal(t, monies.GBP, pound.Currency().Code)
}

func TestMoney_Amount(t *testing.T) {
	pound, err := monies.New(100, monies.GBP)
	require.NoError(t, err)

	assert.Equal(t, int64(100), pound.Amount())
}

func TestDefaultMarshal(t *testing.T) {
	given, err := monies.New(12345, monies.IQD)
	assert.NoError(t, err)
	expected := `{"amount":12345,"currency":"IQD"}`

	b, err := json.Marshal(given)

	if err != nil {
		t.Error(err)
	}

	if string(b) != expected {
		t.Errorf("Expected %s got %s", expected, string(b))
	}

	given = monies.Money{}
	expected = `{"amount":0,"currency":""}`

	b, err = json.Marshal(given)

	if err != nil {
		t.Error(err)
	}

	if string(b) != expected {
		t.Errorf("Expected %s got %s", expected, string(b))
	}
}

func TestCustomMarshal(t *testing.T) {
	given, err := monies.New(12345, monies.IQD)
	assert.NoError(t, err)

	expected := `{"amount":12345,"currency_code":"IQD","currency_fraction":3}`
	monies.MarshalJSON = func(m monies.Money) ([]byte, error) {
		buff := bytes.NewBufferString(fmt.Sprintf(`{"amount": %d, "currency_code": "%s", "currency_fraction": %d}`, m.Amount(), m.Currency().Code, m.Currency().Fraction))
		return buff.Bytes(), nil
	}

	b, err := json.Marshal(given)

	if err != nil {
		t.Error(err)
	}

	if string(b) != expected {
		t.Errorf("Expected %s got %s", expected, string(b))
	}
}

func TestDefaultUnmarshal(t *testing.T) {
	type testCase struct {
		Name          string
		Input         []byte
		UnexpectedErr bool
		ExpectedErr   error
		Expected      monies.Money
	}

	var testCases = []testCase{
		{
			Name:          "SUCCESS",
			Input:         []byte(`{"amount": 100, "currency":"USD"}`),
			UnexpectedErr: false,
			ExpectedErr:   nil,
			Expected:      MustNew(100, monies.USD),
		},
		{
			Name:          "UNDEFINED_CURRENCY",
			Input:         []byte(`{"amount": 100, "currency":"UNDEFINED_CURRENCY"}`),
			UnexpectedErr: false,
			ExpectedErr:   monies.ErrCurrencyNotFound,
		},
		{
			Name:          "INVALID_AMOUNT",
			Input:         []byte(`{"amount": "foo", "currency":"UNDEFINED_CURRENCY"}`),
			UnexpectedErr: true,
		},
		{
			Name:          "WRONG_BYTES",
			Input:         []byte(`1112`),
			UnexpectedErr: true,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.Name, func(t *testing.T) {
			var m monies.Money
			err := json.Unmarshal(tC.Input, &m)
			if tC.UnexpectedErr {
				assert.Error(t, err)
			}

			if tC.ExpectedErr != nil {
				assert.ErrorIs(t, tC.ExpectedErr, err)
			}

			if tC.ExpectedErr == nil && !tC.UnexpectedErr {
				assert.Equal(t, tC.Expected, m)
			}
		})
	}
}

func TestCustomUnmarshal(t *testing.T) {
	given := `{"amount": 10012, "currency_code":"USD", "currency_fraction":2}`
	expected := "$100.12"
	monies.UnmarshalJSON = func(m *monies.Money, b []byte) error {
		data := make(map[string]interface{})
		err := json.Unmarshal(b, &data)
		if err != nil {
			return err
		}

		ref, err := monies.New(int64(data["amount"].(float64)), monies.CurrencyCode(data["currency_code"].(string)))
		require.NoError(t, err)

		*m = ref
		return nil
	}

	var m monies.Money
	err := json.Unmarshal([]byte(given), &m)
	if err != nil {
		t.Error(err)
	}

	if m.Display() != expected {
		t.Errorf("Expected %s got %s", expected, m.Display())
	}
}

func TestTextMarshaller(t *testing.T) {
	type testCase struct {
		Name     string
		Input    monies.Money
		Expected string
	}

	var testCases = []testCase{
		{
			Name:     "SUCCESS",
			Input:    MustNew(100, monies.USD),
			Expected: "$1.00",
		},
	}

	for _, tC := range testCases {
		t.Run(tC.Name, func(t *testing.T) {
			result, err := tC.Input.MarshalText()
			assert.NoError(t, err)

			assert.Equal(t, tC.Expected, string(result))
		})
	}

}
