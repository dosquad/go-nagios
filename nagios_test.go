package nagios_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/dosquad/go-nagios"
	"github.com/google/go-cmp/cmp"
)

func TestResults(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		exitCodes  []nagios.ExitCode
		expectExit nagios.ExitCode
	}{
		{"testing CRITICAL+OK", []nagios.ExitCode{nagios.CRITICAL, nagios.OK}, nagios.CRITICAL},
		{"testing WARNING+OK", []nagios.ExitCode{nagios.WARNING, nagios.OK}, nagios.WARNING},
		{"testing UNKNOWN+OK", []nagios.ExitCode{nagios.UNKNOWN, nagios.OK}, nagios.UNKNOWN},
		{"testing nothing (default OK)", []nagios.ExitCode{}, nagios.OK},
		{"testing OK", []nagios.ExitCode{nagios.OK}, nagios.OK},

		{"testing CRITICAL+OK+CRITICAL", []nagios.ExitCode{nagios.CRITICAL, nagios.OK, nagios.CRITICAL}, nagios.CRITICAL},
		{"testing WARNING+OK+WARNING", []nagios.ExitCode{nagios.WARNING, nagios.OK, nagios.WARNING}, nagios.WARNING},
		{"testing UNKNOWN+OK+UNKNOWN", []nagios.ExitCode{nagios.UNKNOWN, nagios.OK, nagios.UNKNOWN}, nagios.UNKNOWN},

		{
			"testing CRITICAL+WARNING+OK+CRITICAL",
			[]nagios.ExitCode{nagios.CRITICAL, nagios.WARNING, nagios.OK, nagios.CRITICAL},
			nagios.CRITICAL,
		},
		{"testing WARNING+OK+CRITICAL", []nagios.ExitCode{nagios.WARNING, nagios.OK, nagios.CRITICAL}, nagios.CRITICAL},
		{"testing WARNING+UNKNOWN", []nagios.ExitCode{nagios.WARNING, nagios.UNKNOWN}, nagios.UNKNOWN},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			res := &nagios.Result{}

			for _, applyCode := range tt.exitCodes {
				res.SetExitCode(applyCode)
			}

			if res.GetExitCode() != tt.expectExit.Int() {
				t.Errorf("nagios.Result.ExitCode: got '%d', want '%d'", res.ExitCode, tt.expectExit)
			}
		})
	}
}

func TestNagiosSprint(t *testing.T) {
	t.Parallel()
	tests := []struct {
		exitVal   nagios.ExitCode
		expectVal nagios.ExitCode
	}{
		{nagios.OK, nagios.OK},
		{nagios.WARNING, nagios.WARNING},
		{nagios.CRITICAL, nagios.CRITICAL},
		{nagios.UNKNOWN, nagios.UNKNOWN},
		{-1, nagios.UNKNOWN},
		{4, nagios.UNKNOWN},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("%d", tt.exitVal), func(t *testing.T) {
			t.Parallel()

			res := &nagios.Result{}
			exitVal := tt.exitVal
			res.SetExitCode(tt.exitVal)
			res.SetText("exit message")

			output, exitCode := nagios.Sprint(res)

			expectMessage := fmt.Sprintf("%s: exit message|\n", exitVal.String())
			if int(exitVal) != int(tt.expectVal) {
				expectMessage = fmt.Sprintf("%s: Exit code '%d' is not valid :exit message|\n", exitVal.String(), exitVal)
				exitVal = nagios.UNKNOWN
			}

			// t.Logf("nagios.Exit: exitVal got '%d'", exitCode)
			if !exitVal.Equal(exitCode) {
				t.Errorf("nagios.Exit: exitVal got '%d', want '%d'", exitCode, exitVal)
			}

			if diff := cmp.Diff(output, expectMessage); diff != "" {
				t.Errorf("nagios.Exit: output -got +want:\n%s", diff)
			}
		})
	}
}

func TestNagiosSprint_Multiline(t *testing.T) {
	t.Parallel()

	res := &nagios.Result{}
	exitVal := nagios.WARNING
	res.SetExitCode(exitVal)
	res.SetText("exit message")
	res.AddLine("line 1")
	res.AddLine("line 2")

	output, exitCode := nagios.Sprint(res)

	expectMessage := fmt.Sprintf("%s: exit message|\nline 1\nline 2\n\n", exitVal.String())
	if !exitVal.Equal(exitCode) {
		t.Errorf("nagios.Exit: exitVal got '%d', want '%d'", exitCode, exitVal)
	}

	if diff := cmp.Diff(output, expectMessage); diff != "" {
		t.Errorf("nagios.Exit: output -got +want:\n%s", diff)
	}
}

func TestNagiosSprint_Perfdata(t *testing.T) {
	t.Parallel()

	res := &nagios.Result{}
	exitVal := nagios.WARNING
	res.SetExitCode(exitVal)
	res.SetText("exit message")
	res.AddPerfData("key1", "value1")
	res.AddPerfData("key2", "value2")

	output, exitCode := nagios.Sprint(res)

	expectMessage := fmt.Sprintf("%s: exit message|key1=value1 key2=value2\n", exitVal.String())
	if !exitVal.Equal(exitCode) {
		t.Errorf("nagios.Exit: exitVal got '%d', want '%d'", exitCode, exitVal)
	}

	if diff := cmp.Diff(output, expectMessage); diff != "" {
		t.Errorf("nagios.Exit: output -got +want:\n%s", diff)
	}
}

func TestNagiosSprint_Multiline_Perfdata(t *testing.T) {
	t.Parallel()

	res := &nagios.Result{}
	exitVal := nagios.WARNING
	res.SetExitCodeInt(exitVal.Int())
	res.SetText("exit message")
	res.AddLine("line 1")
	res.AddLine("line 2")
	res.AddPerfData("key1", "value1")
	res.AddPerfData("key2", "value2")

	output, exitCode := nagios.Sprint(res)

	expectMessage := fmt.Sprintf("%s: exit message|key1=value1 key2=value2\nline 1\nline 2\n\n", exitVal.String())
	if !exitVal.Equal(exitCode) {
		t.Errorf("nagios.Exit: exitVal got '%d', want '%d'", exitCode, exitVal)
	}

	if diff := cmp.Diff(output, expectMessage); diff != "" {
		t.Errorf("nagios.Exit: output -got +want:\n%s", diff)
	}
}

func TestNagiosSprint_Error(t *testing.T) {
	t.Parallel()

	err := errors.New("I'm a little error teapot")
	exitVal := nagios.CRITICAL
	res := &nagios.Result{}
	res.Error(err)

	output, exitCode := nagios.Sprint(res)

	expectMessage := fmt.Sprintf("%s: I'm a little error teapot|\n", exitVal.String())
	if !exitVal.Equal(exitCode) {
		t.Errorf("nagios.Exit: exitVal got '%d', want '%d'", exitCode, exitVal)
	}

	if diff := cmp.Diff(output, expectMessage); diff != "" {
		t.Errorf("nagios.Exit: output -got +want:\n%s", diff)
	}
}

func TestNagiosPrint(t *testing.T) {
	t.Parallel()

	err := errors.New("I'm a little error teapot")
	exitVal := nagios.CRITICAL
	res := &nagios.Result{}
	res.Error(err)

	output, exitCode := nagios.Sprint(res)

	exitPrintCode := nagios.Print(res)

	expectMessage := fmt.Sprintf("%s: I'm a little error teapot|\n", exitVal.String())
	if !exitVal.Equal(exitCode) {
		t.Errorf("nagios.Exit: exitVal got '%d', want '%d'", exitCode, exitVal)
	}

	if exitCode != exitPrintCode {
		t.Errorf("nagios.Exit: exitPrintCode got '%d', want '%d'", exitPrintCode, exitCode)
	}

	if diff := cmp.Diff(output, expectMessage); diff != "" {
		t.Errorf("nagios.Exit: output -got +want:\n%s", diff)
	}
}
