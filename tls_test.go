package nagios_test

import (
	"testing"

	"github.com/dosquad/go-nagios"
)

func TestStringToTLSVersion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		tlsVersion    string
		expectVersion nagios.TLSVersion
		expectError   bool
	}{
		{"TLS1.0", "tlsv1.0", nagios.VersionTLS10, false},
		{"TLS1.0", "tls10", nagios.VersionTLS10, false},

		{"TLS1.1", "tlsv1.1", nagios.VersionTLS11, false},
		{"TLS1.1", "tls11", nagios.VersionTLS11, false},

		{"TLS1.2", "tlsv1.2", nagios.VersionTLS12, false},
		{"TLS1.2", "tls12", nagios.VersionTLS12, false},

		{"TLS1.3", "tlsv1.3", nagios.VersionTLS13, false},
		{"TLS1.3", "tls13", nagios.VersionTLS13, false},

		{"Invalid Fallback Default", "something", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			v, err := nagios.StringToTLSVersion(tt.tlsVersion)
			if !tt.expectError && err != nil {
				t.Errorf("nagios.StringToTLSVersion: error want nil, got '%s'", err)
			}

			if v != tt.expectVersion {
				t.Errorf("nagios.StringToTLSVersion: string to tls want '%d', got '%d'", tt.expectVersion, v)
			}

			vt := nagios.StringToTLSVersionDefault(tt.tlsVersion, nagios.VersionTLS10)
			if tt.expectError {
				if vt != nagios.VersionTLS10 {
					t.Errorf("nagios.StringToTLSVersionDefault: string to tls want '%d', got '%d'", nagios.VersionTLS10, vt)
				}
			} else {
				if vt != tt.expectVersion {
					t.Errorf("nagios.StringToTLSVersionDefault: string to tls want '%d', got '%d'", tt.expectVersion, vt)
				}
			}
		})
	}
}

func TestStringToTLSVersion_Invalid(t *testing.T) {
	t.Parallel()

	v, err := nagios.StringToTLSVersion("sslv3.0")
	if err == nil {
		t.Error("nagios.StringToTLSVersion: expected error, received nil")
	}

	if v != 0 {
		t.Errorf("nagios.StringToTLSVersion: got '%d' want '%d'", v, 0)
	}
}
