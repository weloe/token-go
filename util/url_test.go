package util

import (
	"log"
	"testing"
)

func TestSpliceUrl(t *testing.T) {
	type args struct {
		u1 string
		u2 string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "success1",
			args: args{
				u1: "http://domain.com",
				u2: "/sso/auth",
			},
			want: "http://domain.com/sso/auth",
		},
		{
			name: "success2",
			args: args{
				u1: "",
				u2: "http://domain.com/sso/auth",
			},
			want: "http://domain.com/sso/auth",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SpliceUrl(tt.args.u1, tt.args.u2); got != tt.want {
				t.Errorf("SpliceUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasUrl(t *testing.T) {
	strings := []string{"1", "2", "3"}
	hasUrl := HasUrl(strings, "2")
	if !hasUrl {
		t.Errorf("HasUrl() = %v, want %v", false, true)
	}
}

func TestMatchUrl(t *testing.T) {
	all := "*"
	if !MatchUrl(all, "123") {
		t.Errorf("HasUrl() = %v, want %v", false, true)
	}
}

func TestIsValidUrl(t *testing.T) {
	type args struct {
		u1 string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{

		{
			name: "error url1",
			args: args{u1: "htp:23/asd"},
			want: false,
		},
		{
			name: "validated",
			args: args{u1: "http://123.com:90//ac"},
			want: true,
		},
		{
			name: "error url2",
			args: args{u1: "http:/123.com:90//ac"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidUrl(tt.args.u1); got != tt.want {
				t.Errorf("IsValidUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddQueryMap(t *testing.T) {
	m := make(map[string]string)
	m["1"] = "2"
	m["2"] = "3"

	if got := AddQueryMap("/sso/auth", m); got != "/sso/auth?1=2&2=3" && got != "/sso/auth?2=3&1=2" {
		t.Errorf("AddQueryMap() = %v, want %v", got, "/sso/auth?1=2&2=3")
	}

	if got := AddQueryMap("/sso/auth?", m); got != "/sso/auth?1=2&2=3" && got != "/sso/auth?2=3&1=2" {
		t.Errorf("AddQueryMap() = %v, want %v", got, "/sso/auth?1=2&2=3")
	}
}

func TestAddQuery(t *testing.T) {
	if got := AddQuery("/sso", "1", "2"); got != "/sso?1=2" {
		t.Errorf("AddQuery() = %v, want %v", got, "/sso?1=2")
	}
	if got := AddQuery("/sso?", "1", "2"); got != "/sso?1=2" {
		t.Errorf("AddQuery() = %v, want %v", got, "/sso?1=2")
	}
}

func TestMapToQuery(t *testing.T) {
	m := make(map[string]string)
	m["1"] = "2"
	query := MapToQuery(m)
	if query != "1=2" {
		t.Errorf("MapToQuery() = %v, want %v", query, "1=2")
	}
}

func TestEncode(t *testing.T) {
	encode := Encode("abc123==123")
	log.Print("Encode(\"abc123==123\") = " + encode)
}

func TestAddQueryValue(t *testing.T) {
	if got := AddQueryValue("/sso/auth?back=http://123.com/login", "ticket=23324"); got != "/sso/auth?back=http://123.com/login&ticket=23324" {
		t.Errorf("AddQueryValue() = %v, want %v", got, "/sso/auth?back=http://123.com/login&ticket=23324")
	}
	if got := AddQueryValue("/sso/auth?back=http://123.com/login?", "ticket=23324"); got != "/sso/auth?back=http://123.com/login?ticket=23324" {
		t.Errorf("AddQueryValue() = %v, want %v", got, "/sso/auth?back=http://123.com/login?ticket=23324")
	}
	if got := AddQueryValue("/sso/auth?back=http://123.com/login?ticket=123", "redirect=23324"); got != "/sso/auth?back=http://123.com/login?ticket=123&redirect=23324" {
		t.Errorf("AddQueryValue() = %v, want %v", got, "/sso/auth?back=http://123.com/login?ticket=123&redirect=23324")
	}
	if got := AddQueryValue("/sso/auth?back=http://123.com/login?ticket=123&", "redirect=23324"); got != "/sso/auth?back=http://123.com/login?ticket=123&redirect=23324" {
		t.Errorf("AddQueryValue() = %v, want %v", got, "/sso/auth?back=http://123.com/login?ticket=123&redirect=23324")
	}
}
