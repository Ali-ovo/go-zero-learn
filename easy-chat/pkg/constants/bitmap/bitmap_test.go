package bitmap

import (
	"testing"
)

func TestIsSet(t *testing.T) {
	b := NewBitmap(20)

	b.Set("ppp")
	b.Set("222")
	b.Set("ppp")
	b.Set("ccc")

	for _, bit := range b.bits {
		t.Logf("%b, %v", bit, bit)
	}

}
