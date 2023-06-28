package monies

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

var (
	ErrCurrencyMismatch = errors.New("currencies mismatched")
	ErrNegativeSplit    = errors.New("split be must be positive")
	ErrNoRatios         = errors.New("no ratios provided")
	ErrInvalidText      = errors.New("invalid text")
)

// Money represents a monetary value
type Money struct {
	amount   int64
	currency Currency
}

// New creates and returns new instance of Money.
func New(amount int64, code CurrencyCode) (m Money, err error) {
	currency, err := CurrencyByCode(code)
	if err != nil {
		return m, ErrCurrencyNotFound
	}

	return Money{
		amount:   amount,
		currency: currency,
	}, nil
}

func (m Money) Currency() Currency {
	return m.currency
}

func (m Money) Amount() int64 {
	return m.amount
}

func (m Money) String() string {
	// Work with absolute amount value
	sa := strconv.FormatInt(absolute(m.amount), 10)
	c := m.currency

	if len(sa) <= m.currency.Fraction {
		sa = strings.Repeat("0", c.Fraction-len(sa)+1) + sa
	}

	if c.Thousand != "" {
		for i := len(sa) - c.Fraction - 3; i > 0; i -= 3 {
			sa = sa[:i] + c.Thousand + sa[i:]
		}
	}

	if c.Fraction > 0 {
		sa = sa[:len(sa)-c.Fraction] + c.Decimal + sa[len(sa)-c.Fraction:]
	}
	sa = strings.Replace(c.Template, "1", sa, 1)
	sa = strings.Replace(sa, "$", c.Grapheme, 1)

	// Add minus sign for negative amount.
	if m.amount < 0 {
		sa = "-" + sa
	}

	return sa
}

func (m Money) AsMajorUnits() float64 {
	if m.currency.Fraction == 0 {
		return float64(m.amount)
	}

	return float64(m.amount) / float64(math.Pow10(m.currency.Fraction))
}

func (m *Money) UnmarshalJSON(b []byte) error {
	type moneyJSON struct {
		Currency CurrencyCode `json:"currency"`
		Amount   int64        `json:"amount"`
	}

	var ref moneyJSON
	err := json.Unmarshal(b, &ref)
	if err != nil {
		return err
	}

	money, err := New(ref.Amount, ref.Currency)
	if err != nil {
		return err
	}

	*m = money
	return nil
}
func (m *Money) UnmarshalText(text []byte) error {
	if len(text) < 5 {
		return ErrInvalidText
	}

	currencyCode := strings.TrimSpace(string(text)[len(text)-4:])
	amountStr := string(text[:len(text)-4])

	currency, err := CurrencyByCode(CurrencyCode(currencyCode))
	if err != nil {
		return err
	}

	majorUnitsStr := strings.Split(amountStr, currency.Decimal)[0]
	minorUnitsStr := strings.Split(amountStr, currency.Decimal)[1]

	majorUnits, err := strconv.ParseInt(majorUnitsStr, 10, 64)
	if err != nil {
		return err
	}

	minorUnits, err := strconv.ParseInt(minorUnitsStr, 10, 64)
	if err != nil {
		return err
	}

	var amount = majorUnits*int64(math.Pow10(currency.Fraction)) + minorUnits

	money, _ := New(amount, CurrencyCode(currencyCode))

	*m = money

	return nil
}

func (m Money) MarshalText() ([]byte, error) {
	var majorUnits = m.amount / int64(math.Pow10(m.currency.Fraction))
	var majorUnitsStr = strconv.FormatInt(majorUnits, 10)

	var minorUnits = m.amount - (majorUnits * int64(math.Pow10(m.currency.Fraction)))
	var minorUnitsStr = strconv.FormatInt(minorUnits, 10)

	if len(minorUnitsStr) <= m.currency.Fraction {
		minorUnitsStr = strings.Repeat("0", m.currency.Fraction-len(minorUnitsStr)) + minorUnitsStr
	}

	return []byte(fmt.Sprintf(`%s%s%s %s`, majorUnitsStr, m.currency.Decimal, minorUnitsStr, m.currency.Code)), nil
}

func (m Money) MarshalJSON() ([]byte, error) {
	buff := bytes.NewBufferString(fmt.Sprintf(`{"amount": %d, "currency": "%s"}`, m.Amount(), m.Currency().Code))
	return buff.Bytes(), nil
}

func (m Money) SameCurrency(om Money) bool {
	return m.currency == om.currency
}

func (m Money) assertSameCurrency(om Money) error {
	if !m.SameCurrency(om) {
		return ErrCurrencyMismatch
	}

	return nil
}

// Compare methods

func (m Money) Equals(om Money) (bool, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return false, err
	}

	return m.amount == om.amount, nil
}

func (m Money) Less(om Money) (bool, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return false, err
	}

	return m.amount < om.amount, nil
}

// Asserts

func (m Money) IsZero() bool {
	return m.amount == 0
}

func (m Money) IsPositive() bool {
	return m.amount > 0
}

func (m Money) IsNegative() bool {
	return m.amount < 0
}

// Operations

func (m Money) Absolute() Money {
	return Money{amount: absolute(m.amount), currency: m.currency}
}

func (m Money) Negative() Money {
	if m.amount < 0 {
		return m
	}
	return Money{amount: -1 * m.amount, currency: m.currency}
}

func (m Money) Add(om Money) (Money, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return om, err
	}

	return Money{amount: m.amount + om.amount, currency: m.currency}, nil
}

func (m Money) Subtract(om Money) (Money, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return om, err
	}

	return Money{amount: m.amount - om.amount, currency: m.currency}, nil
}

func (m Money) Multiply(mul int64) Money {
	return Money{amount: m.amount * mul, currency: m.currency}
}

func (m Money) Round() Money {
	return Money{amount: round(m.amount, m.currency.Fraction), currency: m.currency}
}

// Helpers

// Split tries to evenly distribute the value of the Money struct among the parties.
// If there are not enough pennies to fully distribute, the remainder will be distributed round-robin amongst the parties.
func (m Money) Split(n int) ([]Money, error) {
	if n <= 0 {
		return nil, ErrNegativeSplit
	}

	a := m.amount / int64(n)
	ms := make([]Money, n)

	for i := 0; i < n; i++ {
		ms[i] = Money{amount: a, currency: m.currency}
	}

	r := modulus(m.amount, int64(n))
	l := absolute(r)
	// Add leftovers to the first parties.

	v := int64(1)
	if m.amount < 0 {
		v = -1
	}
	for p := 0; l != 0; p++ {
		ms[p].amount = ms[p].amount + v
		l--
	}

	return ms, nil
}

// Allocate returns slice of Money structs with split Self value in given Ratios.
// It lets split money by given Ratios without losing pennies and as Split operations distributes
// leftover pennies amongst the parties with round-robin principle.
func (m Money) Allocate(rs ...int) ([]Money, error) {
	if len(rs) == 0 {
		return nil, ErrNoRatios
	}

	// Calculate sum of Ratios.
	var sum int
	for _, r := range rs {
		sum += r
	}

	var total int64
	ms := make([]Money, 0, len(rs))
	for _, r := range rs {
		party := Money{
			amount:   allocate(m.amount, r, sum),
			currency: m.currency,
		}

		ms = append(ms, party)
		total += party.amount
	}

	// Calculate leftover value and divide to first parties.
	lo := m.amount - total
	sub := int64(1)
	if lo < 0 {
		sub = -sub
	}

	for p := 0; lo != 0; p++ {
		ms[p].amount = ms[p].amount + sub
		lo -= sub
	}

	return ms, nil
}

// It panics if error occurs.
func MustNew(amount int64, code CurrencyCode) Money {
	m, err := New(amount, code)
	if err != nil {
		panic(err)
	}

	return m
}
