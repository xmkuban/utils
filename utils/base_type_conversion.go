package utils

func ToPointInt(i int) *int {
	return &i
}

func ToPointString(s string) *string {
	return &s
}

func ToPointBool(b bool) *bool {
	return &b
}
func ToPointInt64(i int64) *int64 {
	return &i
}

func ToPointFloat64(f float64) *float64 {
	return &f
}

func ToPointArrayInt(i []int) *[]int {
	if len(i) == 0 {
		return nil
	}
	return &i
}
func ToPointArrayInt64(i []int64) *[]int64 {
	if len(i) == 0 {
		return nil
	}
	return &i
}

func PointToInt(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

func PointToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func PointToBool(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

func PointToInt64(i *int64) int64 {
	if i == nil {
		return 0
	}
	return *i
}
func PointToFloat64(f *float64) float64 {
	if f == nil {
		return 0
	}
	return *f
}
