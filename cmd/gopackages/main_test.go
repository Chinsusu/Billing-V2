package main

import "testing"

func TestIncludePackage(t *testing.T) {
	tests := []struct {
		name string
		pkg  string
		want bool
	}{
		{
			name: "billing package",
			pkg:  "github.com/Chinsusu/Billing-V2/internal/modules/wallet",
			want: true,
		},
		{
			name: "frontend dependency package",
			pkg:  "github.com/Chinsusu/Billing-V2/frontend/node_modules/flatted/golang/pkg/flatted",
			want: false,
		},
		{
			name: "windows style frontend dependency package",
			pkg:  "github.com\\Chinsusu\\Billing-V2\\frontend\\node_modules\\flatted\\golang\\pkg\\flatted",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := includePackage(tt.pkg); got != tt.want {
				t.Fatalf("includePackage(%q) = %v, want %v", tt.pkg, got, tt.want)
			}
		})
	}
}
