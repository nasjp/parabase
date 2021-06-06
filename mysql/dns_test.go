package mysql

import (
	"testing"
)

func TestDNS(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		rawDataSourceName string
		want              string
		wantErr           bool
	}{
		{"Default", "user:password@/raw", "user:password@/replace", false},
		{"IP", "user:password@tcp(host:port)/raw", "user:password@tcp(host:port)/replace", false},
		{"Config", "user:password@/raw?parseTime=true&tls=true&multiStatements=true&charset=utf8mb4", "user:password@/replace?parseTime=true&tls=true&multiStatements=true&charset=utf8mb4", false},
		{"NoSelectedDB", "user:password@/", "user:password@/replace", false},
		{"Invalid", "user:password@", "", true},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := dns(tt.rawDataSourceName, "replace")
			if (err != nil) && !tt.wantErr {
				t.Fatal(err)
			}

			if got != tt.want {
				t.Errorf("dns replace failed, got: '%s', want: '%s'", got, tt.want)
			}
		})

	}
}
