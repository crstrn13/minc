package providers

import "testing"

func TestPublishPort(t *testing.T) {
	cases := []struct {
		host, container int
		want            string
	}{
		{0, 6443, "127.0.0.1::6443"},
		{0, 80, "127.0.0.1::80"},
		{9443, 443, "127.0.0.1:9443:443"},
		{443, 443, "443:443"}, // <1024 binds all interfaces
	}
	for _, c := range cases {
		if got := publishPort(c.host, c.container); got != c.want {
			t.Errorf("publishPort(%d,%d)=%q want %q", c.host, c.container, got, c.want)
		}
	}
}

func TestParseHostPort(t *testing.T) {
	cases := []struct {
		in      string
		want    int
		wantErr bool
	}{
		{"127.0.0.1:49153", 49153, false},
		{"0.0.0.0:32768\n[::]:32768", 32768, false},
		{"  127.0.0.1:12345  \n", 12345, false},
		{"", 0, true},
		{"garbage", 0, true},
	}
	for _, c := range cases {
		got, err := ParseHostPort(c.in)
		if c.wantErr {
			if err == nil {
				t.Errorf("ParseHostPort(%q) expected error", c.in)
			}
			continue
		}
		if err != nil || got != c.want {
			t.Errorf("ParseHostPort(%q)=%d,%v want %d", c.in, got, err, c.want)
		}
	}
}
