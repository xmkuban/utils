package utils

import (
	"testing"
)

func TestGenRandomString(t *testing.T) {
	t.Log(GenRandomString(32))
}
