package utils

import (
	"strings"
	"testing"
)

func TestValidateSecret(t *testing.T) {
	tests := []struct {
		name    string
		secret  string
		wantErr bool
	}{
		{name: "empty", secret: "", wantErr: true},
		{name: "default example secret", secret: "secret", wantErr: true},
		{name: "31 bytes", secret: strings.Repeat("a", 31), wantErr: true},
		{name: "32 bytes", secret: strings.Repeat("a", 32), wantErr: false},
		{name: "longer than 32", secret: strings.Repeat("b", 64), wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateSecret(tt.secret)
			if (err != nil) != tt.wantErr {
				t.Fatalf("validateSecret(%q) error = %v, wantErr %v", tt.secret, err, tt.wantErr)
			}
		})
	}
}
