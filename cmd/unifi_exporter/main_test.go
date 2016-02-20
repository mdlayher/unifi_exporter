package main

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/mdlayher/unifi"
)

func Test_pickSites(t *testing.T) {
	var tests = []struct {
		desc   string
		choose string
		sites  []*unifi.Site
		pick   []*unifi.Site
		err    error
	}{
		{
			desc:   "no site chosen",
			choose: "",
			sites: []*unifi.Site{
				{Description: "foo"},
				{Description: "bar"},
				{Description: "baz"},
			},
			pick: []*unifi.Site{
				{Description: "foo"},
				{Description: "bar"},
				{Description: "baz"},
			},
		},
		{
			desc:   "one valid site chosen",
			choose: "bar",
			sites: []*unifi.Site{
				{Description: "foo"},
				{Description: "bar"},
				{Description: "baz"},
			},
			pick: []*unifi.Site{
				{Description: "bar"},
			},
		},
		{
			desc:   "one invalid site chosen",
			choose: "qux",
			sites: []*unifi.Site{
				{Description: "foo"},
				{Description: "bar"},
				{Description: "baz"},
			},
			err: errors.New("was not found in UniFi Controller"),
		},
	}

	for i, tt := range tests {
		t.Logf("[%02d] test %q", i, tt.desc)

		pick, err := pickSites(tt.choose, tt.sites)
		if want, got := errStr(tt.err), errStr(err); !strings.Contains(got, want) {
			t.Fatalf("unexpected error:\n- want: %v\n-  got: %v",
				want, got)
		}
		if err != nil {
			continue
		}

		if want, got := tt.pick, pick; !reflect.DeepEqual(want, got) {
			t.Fatalf("unexpected sites:\n- want: %v\n-  got: %v",
				want, got)
		}
	}
}

func errStr(err error) string {
	if err == nil {
		return ""
	}

	return err.Error()
}

func Test_sitesString(t *testing.T) {
	var tests = []struct {
		sites []*unifi.Site
		out   string
	}{
		{
			out: "",
		},
		{
			sites: []*unifi.Site{{
				Description: "Foo",
			}},
			out: "Foo",
		},
		{
			sites: []*unifi.Site{
				{Description: "Foo"},
				{Description: "Bar"},
				{Description: "Baz"},
			},
			out: "Foo, Bar, Baz",
		},
	}

	for i, tt := range tests {
		t.Logf("[%02d] out: %q", i, tt.out)

		out := sitesString(tt.sites)
		if want, got := tt.out, out; want != got {
			t.Fatalf("unexpected output:\n- want: %v\n-  got: %v",
				want, got)
		}
	}
}
