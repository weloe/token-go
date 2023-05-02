package ctx

import (
	"testing"
)

func TestDefaultRespImplement_AddCookie(t *testing.T) {
	type args struct {
		name    string
		value   string
		path    string
		domain  string
		timeout int64
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &DefaultRespImplement{}
			r.AddCookie(tt.args.name, tt.args.value, tt.args.path, tt.args.domain, tt.args.timeout)
		})
	}
}

func TestDefaultRespImplement_AddHeader(t *testing.T) {
	type args struct {
		name  string
		value string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &DefaultRespImplement{}
			r.AddHeader(tt.args.name, tt.args.value)
		})
	}
}

func TestDefaultRespImplement_DeleteCookie(t *testing.T) {
	type args struct {
		name   string
		path   string
		domain string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &DefaultRespImplement{}
			r.DeleteCookie(tt.args.name, tt.args.path, tt.args.domain)
		})
	}
}

func TestDefaultRespImplement_SetServer(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &DefaultRespImplement{}
			r.SetServer(tt.args.value)
		})
	}
}
