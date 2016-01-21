package str

import (
	"testing"
)

func TestIsMail1(t *testing.T) {
	got := IsMail("ulric.qin-_abc@a.com")
	expect := true
	if got != expect {
		t.Errorf("IsMail1:\n Expect => %v\n Got    => %v\n", expect, got)
	}
}

func TestIsMail2(t *testing.T) {
	got := IsMail("i@a.com")
	expect := true
	if got != expect {
		t.Errorf("IsMail2:\n Expect => %v\n Got    => %v\n", expect, got)
	}
}

func TestIsMail3(t *testing.T) {
	got := IsMail("ia.com")
	expect := false
	if got != expect {
		t.Errorf("IsMail3:\n Expect => %v\n Got    => %v\n", expect, got)
	}
}

func TestIsPhone1(t *testing.T) {
	got := IsPhone("18655555555")
	expect := true
	if got != expect {
		t.Errorf("TestIsPhone1:\n Expect => %v\n Got    => %v\n", expect, got)
	}
}

func TestIsPhone2(t *testing.T) {
	got := IsPhone("+8618688885555")
	expect := true
	if got != expect {
		t.Errorf("TestIsPhone2:\n Expect => %v\n Got    => %v\n", expect, got)
	}
}

func TestIsPhone3(t *testing.T) {
	got := IsPhone("1861218552")
	expect := false
	if got != expect {
		t.Errorf("TestIsPhone3:\n Expect => %v\n Got    => %v\n", expect, got)
	}
}
