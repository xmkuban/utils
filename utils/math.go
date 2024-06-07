package utils

import (
	"math"

	"github.com/shopspring/decimal"
)

const (
	Round     = 0
	RoundUp   = 1
	RoundDown = 2
)

type Math struct {
	Round int
}

func NewMath() *Math {
	return &Math{}
}

// FloatAddNum 以高精度加法的方式精确地添加两个浮点数。
// 这样做可以避免浮点数运算带来的精度损失。
func (m *Math) FloatAddNum(a, b float64) float64 {
	_a := decimal.NewFromFloat(a)
	_b := decimal.NewFromFloat(b)
	c := _a.Add(_b)
	_c, _ := c.Float64()
	return _c
}

// FloatSubNum 以高精度减法的方式精确地减去两个浮点数。
func (m *Math) FloatSubNum(a, b float64) float64 {
	_a := decimal.NewFromFloat(a)
	_b := decimal.NewFromFloat(b)
	c := _a.Sub(_b)
	_c, _ := c.Float64()
	return _c
}

// FloatMulNum 精确地计算两个浮点数的乘积，支持指定的小数位数。
func (m *Math) FloatMulNum(a, b float64, exp int32) float64 {
	if a == 0 || b == 0 {
		return 0
	}
	_a := decimal.NewFromFloat(a)
	_b := decimal.NewFromFloat(b)
	var c decimal.Decimal
	if exp == -1 {
		c = _a.Mul(_b)
	} else {
		c = _a.Mul(_b).Truncate(exp + 1)
	}
	var _c float64
	switch m.Round {
	case Round:
		_c, _ = c.Round(exp).Float64()
	case RoundUp:
		_c, _ = c.RoundUp(exp).Float64()
	case RoundDown:
		_c, _ = c.RoundDown(exp).Float64()
	}
	return _c
}

// FloatDivNum 精确地计算两个浮点数的商，支持指定的小数位数。
func (m *Math) FloatDivNum(a, b float64, exp int32) float64 {
	if a == 0 || b == 0 {
		return 0
	}
	_a := decimal.NewFromFloat(a)
	_b := decimal.NewFromFloat(b)
	var c decimal.Decimal
	if exp == -1 {
		c = _a.Div(_b)
	} else {
		c = _a.Div(_b).Truncate(exp + 1)
	}
	var _c float64
	switch m.Round {
	case Round:
		_c, _ = c.Round(exp).Float64()
	case RoundUp:
		_c, _ = c.RoundUp(exp).Float64()
	case RoundDown:
		_c, _ = c.RoundDown(exp).Float64()
	}
	return _c
}

func (m *Math) FloatExponent(a float64, exp int32) float64 {
	if a == 0 {
		return 0
	}
	_c, _ := decimal.NewFromFloat(a).Round(exp).Float64()
	return _c
}

func (m *Math) FloatExponentUp(a float64, exp int32) float64 {
	if a == 0 {
		return 0
	}
	_c, _ := decimal.NewFromFloat(a).RoundUp(exp).Float64()
	return _c
}

func (m *Math) FloatExponentDown(a float64, exp int32) float64 {
	if a == 0 {
		return 0
	}
	_c, _ := decimal.NewFromFloat(a).RoundDown(exp).Float64()
	return _c
}

// FloatExponentToStr 计算浮点数的幂，并以字符串形式返回。
func (m *Math) FloatExponentToStr(a float64, exp int32) string {
	if a == 0 {
		return "0"
	}
	_c := decimal.NewFromFloatWithExponent(a, -exp).String()
	if _c == "" {
		_c = "0"
	}
	return _c
}

// DecimalToFloat 将 decimal.Decimal 转换为 float64。
func DecimalToFloat(a decimal.Decimal) float64 {
	b, _ := a.Float64()
	return b
}

// MaxInt 返回两个整数中的最大值。
func MaxInt(a, b int) int {
	return int(math.Max(float64(a), float64(b)))
}

// MinInt 返回两个整数中的最小值。
func MinInt(a, b int) int {
	return int(math.Min(float64(a), float64(b)))
}
