package name

import (
	"testing"
)

func TestCreateFromString(t *testing.T) {
	cases := []struct {
		in string
		test string
		match bool
	}{
		{"^<foo><bar>$", "/foo/bar", true},
	}

	for _, c := range cases {
		regularName, err := CreateFromString(c.in)
		if err != nil {
			t.Errorf("Failed to create a RegularName from %s", c.in)
		}

		testName, err := Parse(c.test)
		if err != nil {
			t.Errorf("Failed to parse name string %s", c.test)
		}

		matched := regularName.Accepts(testName)
		if matched != c.match {
			if c.match {
				t.Errorf("Expected %s to accept %s, but it did not", c.in, c.test)
			} else {
				t.Errorf("Did not expect %s to accept %s, but it did", c.in, c.test)
			}
		}
	}
}
