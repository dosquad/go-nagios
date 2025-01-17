package nagios

import (
	"strings"
)

// ExitCode is an extension of int with some helpers for
// use with Nagios exit codes.
type ExitCode int

const (
	// OK is the Nagios OK Exitcode.
	//nolint:varnamelen // it is the constant for OK.
	OK ExitCode = 0

	// WARNING is the Nagios WARNING Exitcode.
	WARNING ExitCode = 1

	// CRITICAL is the Nagios CRITICAL Exitcode.
	CRITICAL ExitCode = 2

	// UNKNOWN is the Nagios UNKNOWN Exitcode.
	UNKNOWN ExitCode = 3
)

const (
	okString       = "OK"
	warningString  = "WARNING"
	criticalString = "CRITICAL"
	unknownString  = "UNKNOWN"
)

// Int returns the ExitCode as a nagios compatible exit code.
func (e ExitCode) Int() int {
	switch e {
	case OK, WARNING, CRITICAL, UNKNOWN:
		return int(e)
	}

	return int(UNKNOWN)
}

// IsGreater returns true if the ExitCode is a greater than
// the supplied exit code.
func (e ExitCode) IsGreater(a int) bool {
	return e.IsGreaterCode(ExitCode(a))
}

// IsGreaterCode returns true if the ExitCode is a greater than
// the supplied exit code.
func (e ExitCode) IsGreaterCode(exitCode ExitCode) bool {
	if e < 0 && exitCode == 0 {
		return true
	}

	if exitCode < 0 {
		return false
	}

	return e > exitCode
}

// Equal returns true if a is equal to ExitCode.
func (e ExitCode) Equal(a int) bool {
	return int(e) == a
}

// String returns the ExitCode as a string.
func (e ExitCode) String() string {
	switch e {
	case OK:
		return okString
	case WARNING:
		return warningString
	case CRITICAL:
		return criticalString
	case UNKNOWN:
		return unknownString
	}

	return unknownString
}

// ExitCodeFromString converts a supplied string to an ExitCode.
func ExitCodeFromString(exitCode string) ExitCode {
	switch strings.ToUpper(exitCode) {
	case okString:
		return OK
	case warningString:
		return WARNING
	case criticalString:
		return CRITICAL
	case unknownString:
		return UNKNOWN
	}

	return UNKNOWN
}

// ExitCodeFromInt converts a supplied int to an ExitCode.
func ExitCodeFromInt(exitCode int) ExitCode {
	v := ExitCode(exitCode)
	switch v {
	case OK, WARNING, CRITICAL, UNKNOWN:
		return v
	}

	return UNKNOWN
}
