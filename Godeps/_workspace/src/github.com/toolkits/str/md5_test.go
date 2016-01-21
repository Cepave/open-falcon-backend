package str

import (
	"testing"
)

func TestMd5Encode(t *testing.T) {
	got := Md5Encode("abc")
	expect := "900150983cd24fb0d6963f7d28e17f72"
	if got != expect {
		t.Error("Md5Encode:\n Expect => %v\n Got    => %v\n", expect, got)
	}
}
