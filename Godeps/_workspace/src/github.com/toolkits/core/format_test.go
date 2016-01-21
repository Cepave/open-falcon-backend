package core

import (
	"testing"
)

func TestReadableSize(t *testing.T) {
	got := ReadableSize(12.123)
	expect := "12.1B"
	if got != expect {
		t.Errorf("ReadableSize:\n Expect => %v\n Got    => %v\n", expect, got)
	}

	got = ReadableSize(1124.123)
	expect = "1.1K"
	if got != expect {
		t.Errorf("ReadableSize:\n Expect => %v\n Got    => %v\n", expect, got)
	}

	got = ReadableSize(1024 * 1024 * 8)
	expect = "8.0M"
	if got != expect {
		t.Errorf("ReadableSize:\n Expect => %v\n Got    => %v\n", expect, got)
	}

	got = ReadableSize(1024 * 1024 * 1024 * 8)
	expect = "8.0G"
	if got != expect {
		t.Errorf("ReadableSize:\n Expect => %v\n Got    => %v\n", expect, got)
	}

	got = ReadableSize(1024 * 1024 * 1024 * 1024 * 8)
	expect = "8.0T"
	if got != expect {
		t.Errorf("ReadableSize:\n Expect => %v\n Got    => %v\n", expect, got)
	}

	got = ReadableSize(1024 * 1024 * 1024 * 1024 * 1024 * 8)
	expect = "8.0P"
	if got != expect {
		t.Errorf("ReadableSize:\n Expect => %v\n Got    => %v\n", expect, got)
	}

	got = ReadableSize(1024 * 1024 * 1024 * 1024 * 1024 * 1024 * 8)
	expect = "TooLarge"
	if got != expect {
		t.Errorf("ReadableSize:\n Expect => %v\n Got    => %v\n", expect, got)
	}

	got = ReadableSize(0)
	expect = "0.0B"
	if got != expect {
		t.Errorf("ReadableSize:\n Expect => %v\n Got    => %v\n", expect, got)
	}
}
