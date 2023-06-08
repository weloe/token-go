package util

import (
	"testing"
)

func TestAppendStr(t *testing.T) {
	strings := []string{"1", "2"}
	str := AppendStr(strings, "1")
	str = AppendStr(str, "3")
	for i, s := range strings {
		if s == "1" && i == 2 {
			t.Errorf("AppendStr() = %v, want %v", str, []string{"1", "2"})
		}

	}
	for i, s := range strings {
		if i == 2 && s != "3" {
			t.Errorf("AppendStr() = %v, want %v", str, []string{"1", "2", "3"})
		}
	}
}
