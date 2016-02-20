package unifiexporter

import (
	"testing"
)

func Test_siteDescription(t *testing.T) {
	var tests = []struct {
		in  string
		out string
	}{
		{
			in:  "Foo Bar",
			out: "foobar",
		},
		{
			in:  "Foo-Bar_Baz",
			out: "foobarbaz",
		},
		{
			in:  "Foo bar  Baz - _ .qux",
			out: "foobarbazqux",
		},
	}

	for i, tt := range tests {
		t.Logf("[%02d] in: %q, out: %q", i, tt.in, tt.out)

		got := siteDescription(tt.in)
		if want := tt.out; want != got {
			t.Fatalf("unexpected output:\n- want: %v\n-  got: %v",
				want, got)
		}
	}
}
