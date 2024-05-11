package utilities

import (
	"os"
	"testing"
)

func TestRepeatFn(t *testing.T) {
	str1 := Repeat("..", 4)
	str2 := Repeat(".", 2)

	if str1 != "../../../.." {
		t.Fatalf(`repeat fn does not create 4 sub-paths for str1`)
	}

	if str2 != "./." {
		t.Fatalf(`repeat fn does not create 2 sub-paths for str2`)
	}
}

func TestIsInvalidPath(t *testing.T) {
	invalidPath := "doanfkjzx/sdfj931/sdfkjal"
	err := os.Chdir(invalidPath)
	if err == nil {
		t.Fatalf(`Path doanfkjzx/sdfj931/sdfkjal exists`)
	}

	val := IsInvalidPath(invalidPath)

	if val != true {
		t.Fatalf(`InvalidPath returns true`)
	}
}

func TestPrependStrSlice(t *testing.T) {
	sl := []string{"apple", "banana"}

	sl = prependStrSlice(sl, "cookie")
	if len(sl) != 3 {
		t.Fatalf(`Slice of string is not incremented`)
	}

	if sl[0] != "cookie" {
		t.Fatalf(`cookie is not prepended to slice`)
	}
}
