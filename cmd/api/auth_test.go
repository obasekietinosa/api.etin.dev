package main

import "testing"

func TestParseBearerToken(t *testing.T) {
	tests := []struct {
		name    string
		header  string
		want    string
		wantErr bool
	}{
		{name: "standard", header: "Bearer token", want: "token"},
		{name: "lowercase", header: "bearer token", want: "token"},
		{name: "extra spaces", header: "  BEARER   token   ", want: "token"},
		{name: "missing token", header: "Bearer", wantErr: true},
		{name: "wrong scheme", header: "Basic token", wantErr: true},
		{name: "empty", header: "", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseBearerToken(tt.header)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tt.want {
				t.Fatalf("expected %q got %q", tt.want, got)
			}
		})
	}
}
