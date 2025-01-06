package nagios_test

import (
	"testing"

	"github.com/dosquad/go-nagios"
)

func TestExitCodeFromString(t *testing.T) {
	t.Parallel()
	tests := []struct {
		exitCode string
		expect   nagios.ExitCode
	}{
		{"OK", nagios.OK},
		{"WARNING", nagios.WARNING},
		{"CRITICAL", nagios.CRITICAL},
		{"critical", nagios.CRITICAL},
		{"Critical", nagios.CRITICAL},
		{"UNKNOWN", nagios.UNKNOWN},
		{"invalid", nagios.UNKNOWN},
	}

	for _, tt := range tests {
		t.Run(tt.exitCode, func(t *testing.T) {
			t.Parallel()

			v := nagios.ExitCodeFromString(tt.exitCode)

			if !v.Equal(int(tt.expect)) {
				t.Errorf("nagios.ExitCodeFromString: got '%d' want '%d'", v, tt.expect)
			}
		})
	}
}
