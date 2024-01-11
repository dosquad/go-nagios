package nagios

import (
	"crypto/tls"
	"errors"
	"fmt"
	"strings"
)

// TLSVersion is a list of supported TLS versions.
type TLSVersion uint16

const (
	// VersionTLS10 TLSv1.0 .
	VersionTLS10 = TLSVersion(tls.VersionTLS10)
	// VersionTLS11 TLSv1.1 .
	VersionTLS11 = TLSVersion(tls.VersionTLS11)
	// VersionTLS12 TLSv1.2 .
	VersionTLS12 = TLSVersion(tls.VersionTLS12)
	// VersionTLS13 TLSv1.3 .
	VersionTLS13 = TLSVersion(tls.VersionTLS13)
)

// StringToTLSVersionDefault converts a string to a TLS version or returns a default.
func StringToTLSVersionDefault(s string, defaultVersion TLSVersion) TLSVersion {
	if v, err := StringToTLSVersion(s); err == nil {
		return v
	}

	return defaultVersion
}

// ErrUnknownTLSVersion is returned when the TLS version string can not be parsed.
var ErrUnknownTLSVersion = errors.New("unknown TLS version string")

// StringToTLSVersion converts a string to a TLS version or returns a default.
func StringToTLSVersion(tlsName string) (TLSVersion, error) {
	switch strings.ToLower(tlsName) {
	case "tls10", "tlsv1.0":
		return VersionTLS10, nil
	case "tls11", "tlsv1.1":
		return VersionTLS11, nil
	case "tls12", "tlsv1.2":
		return VersionTLS12, nil
	case "tls13", "tlsv1.3":
		return VersionTLS13, nil
	}

	return 0, fmt.Errorf("%w: %s", ErrUnknownTLSVersion, tlsName)
}
