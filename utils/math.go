package utils

import (
	"math"

	"github.com/shopspring/decimal"
)

func FloatAddNum(a, b float64) float64 {
	_a := decimal.NewFromFloat(a)
	_b := decimal.NewFromFloat(b)
	c := _a.Add(_b)
	_c, _ := c.Float64()
	return _c
}

func FloatSubNum(a, b float64) float64 {
	_a := decimal.NewFromFloat(a)
	_b := decimal.NewFromFloat(b)
	c := _a.Sub(_b)
	_c, _ := c.Float64()
	return _c
}

func FloatMulNum(a, b float64, exp int32) float64 {
	_a := decimal.NewFromFloat(a)
	_b := decimal.NewFromFloat(b)
	var c decimal.Decimal
	if exp == -1 {
		c = _a.Mul(_b)
	} else {
		c = _a.Mul(_b).Truncate(exp)
	}
	_c, _ := c.Float64()
	return _c
}

func FloatDivNum(a, b float64, exp int32) float64 {
	_a := decimal.NewFromFloat(a)
	_b := decimal.NewFromFloat(b)
	var c decimal.Decimal
	if exp == -1 {
		c = _a.Div(_b)
	} else {
		c = _a.Div(_b).Truncate(exp)
	}
	_c, _ := c.Float64()
	return _c
}

func FloatExponent(a float64, exp int32) float64 {
	if a == 0 {
		return 0
	}
	_c, _ := decimal.NewFromFloatWithExponent(a, -exp).Float64()
	return _c
}

func DecimalToFloat(a decimal.Decimal) float64 {
	b, _ := a.Float64()
	return b
}

func MaxInt(a, b int) int {
	return int(math.Max(float64(a), float64(b)))
}

func MinInt(a, b int) int {
	return int(math.Min(float64(a), float64(b)))
}
