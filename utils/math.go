package utils

import (
	"math"
	"strconv"

	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
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
	exp = exp + 1
	if a == 0 || b == 0 {
		return 0
	}
	_a := decimal.NewFromFloat(a)
	_b := decimal.NewFromFloat(b)
	var c decimal.Decimal
	if exp == -1 {
		c = _a.Mul(_b)
	} else {
		c = _a.Mul(_b).Truncate(exp)
	}
	_c, _ := c.Float64()
	_c = FormatFloat(_c, int(exp-1))
	return _c
}

func FloatDivNum(a, b float64, exp int32) float64 {
	exp = exp + 1
	if a == 0 || b == 0 {
		return 0
	}
	_a := decimal.NewFromFloat(a)
	_b := decimal.NewFromFloat(b)
	var c decimal.Decimal
	if exp == -1 {
		c = _a.Div(_b)
	} else {
		c = _a.Div(_b).Truncate(exp)
	}
	_c, _ := c.Float64()
	_c = FormatFloat(_c, int(exp-1))
	return _c
}

func FloatExponent(a float64, exp int32) float64 {
	exp = exp + 1
	if a == 0 {
		return 0
	}
	_c, _ := decimal.NewFromFloatWithExponent(a, -exp).Float64()
	_c = FormatFloat(_c, int(exp-1))
	return _c
}

func FloatExponentToStr(a float64, exp int32) string {
	if a == 0 {
		return "0"
	}
	_c := decimal.NewFromFloatWithExponent(a, -exp).String()
	if _c == "" {
		_c = "0"
	}
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

func FormatFloat(f float64, dig int) float64 {
	result := cast.ToFloat64(strconv.FormatFloat(f, 'f', dig+1, 64))

	pow := math.Pow(10, float64(dig))

	return math.Round(result*pow) / pow
}
