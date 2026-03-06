package money

import (
	"errors"

	"github.com/cockroachdb/apd/v3"
)

var mathCtx = apd.Context{
	Precision: 28,

	MaxExponent: 15,
	MinExponent: -15,

	// RoundHalfEven redondeo recomendado por el estándar IEEE 754
	Rounding: apd.RoundHalfEven,

	Traps: apd.Overflow | apd.DivisionByZero | apd.InvalidOperation,
}

var (
	ErrInvalidAmount  = errors.New("formato de dinero invalido")
	ErrMathOverflow   = errors.New("el calculo excedio los limites permitidos")
	ErrDivisionByZero = errors.New("error matematico: division por cero")
)

type Amount struct {
	value *apd.Decimal
}

type Rate struct {
	value *apd.Decimal
}

func NewAmmount(val string) (Amount, error) {
	d, _, err := apd.NewFromString(val)
	if err != nil {
		return Amount{}, ErrInvalidAmount
	}
	return Amount{value: d}, nil
}

func NewRate(val string) (Rate, error) {
	d, _, err := apd.NewFromString(val)
	if err != nil {
		return Rate{}, ErrInvalidAmount
	}
	return Rate{value: d}, nil
}

func (m Rate) Mul(rate string) (Rate, error) {
	r, _, err := apd.NewFromString(rate)
	if err != nil {
		return Rate{}, ErrInvalidAmount
	}

	result := apd.New(0, 0)

	condition, err := mathCtx.Mul(result, m.value, r)
	if err != nil {
		return Rate{}, err
	}
	if condition.Overflow() {
		return Rate{}, ErrMathOverflow
	}

	return Rate{value: result}, nil
}

func (m Amount) Mul(r Rate) (Amount, error) {
	result := apd.New(0, 0)

	condition, err := mathCtx.Mul(result, m.value, r.value)
	if err != nil {
		return Amount{}, err
	}
	if condition.Overflow() {
		return Amount{}, ErrMathOverflow
	}

	return Amount{value: result}, nil
}

func (m Amount) Div(other Amount) (Amount, error) {
	if other.IsZero() {
		return Amount{}, ErrDivisionByZero
	}

	result := apd.New(0, 0)

	condition, err := mathCtx.Quo(result, m.value, other.value)
	if err != nil {
		return Amount{}, err
	}

	if condition.Overflow() {
		return Amount{}, ErrMathOverflow
	}

	return Amount{value: result}, nil
}

func (m Amount) Add(other Amount) (Amount, error) {
	result := apd.New(0, 0)
	condition, err := mathCtx.Add(result, m.value, other.value)
	if err != nil {
		return Amount{}, err
	}
	if condition.Overflow() {
		return Amount{}, ErrMathOverflow
	}

	return Amount{value: result}, nil

}

func (m Amount) Sub(other Amount) (Amount, error) {
	result := apd.New(0, 0)
	condition, err := mathCtx.Sub(result, m.value, other.value)
	if err != nil {
		return Amount{}, err
	}
	if condition.Overflow() {
		return Amount{}, ErrMathOverflow
	}

	return Amount{value: result}, nil

}

func (m Amount) IsNegative() bool {
	if m.value == nil {
		return false
	}
	return m.value.Sign() < 0
}

func (m Amount) IsZero() bool {
	if m.value == nil {
		return true
	}
	return m.value.IsZero()
}

func (m Rate) IsNegative() bool {
	if m.value == nil {
		return false
	}
	return m.value.Sign() < 0
}

// IsGreaterThan compara el Rate original con otro Rate.
// Retorna true si el Rate original es MAYOR que el otro Rate.
func (r Rate) IsGreaterThan(other Rate) bool {
	if r.value == nil || other.value == nil {
		return false
	}

	result := r.value.Cmp(other.value)

	// Cmp devuelve:
	//  1 si el original es mayor
	//  0 si son exactamente iguales
	// -1 si el original es menor
	if result == 1 {
		return true
	} else {
		return false
	}
}

func (m Amount) AmmountToString() string {
	if m.value == nil {
		return "0.00"
	}
	// return format: 1500000.00
	return m.value.Text('f')
}

func (m Rate) RateToString() string {
	if m.value == nil {
		return "0.00"
	}
	// return format: 1500000.00
	return m.value.Text('f')
}
