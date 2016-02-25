package unifiexporter

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mdlayher/unifi"
	"github.com/prometheus/client_golang/prometheus"
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

func testUniFiClient(t *testing.T, input []byte) (*unifi.Client, func()) {
	unifiServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		_, _ = w.Write(input)
	}))

	c, err := unifi.NewClient(unifiServer.URL, nil)
	if err != nil {
		t.Fatalf("failed to create UniFi client: %v", err)
	}

	return c, func() { unifiServer.Close() }
}

func testCollector(t *testing.T, collector prometheus.Collector) []byte {
	if err := prometheus.Register(collector); err != nil {
		t.Fatalf("failed to register Prometheus collector: %v", err)
	}
	defer prometheus.Unregister(collector)

	promServer := httptest.NewServer(prometheus.Handler())
	defer promServer.Close()

	resp, err := http.Get(promServer.URL)
	if err != nil {
		t.Fatalf("failed to GET data from prometheus: %v", err)
	}
	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read server response: %v", err)
	}

	return buf
}
