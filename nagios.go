package nagios

import (
	"fmt"
	"os"
	"strconv"
)

// Result contains the output used by Nagios.
type Result struct {
	ExitCode  ExitCode
	Text      string
	Prefix    string
	Perfdata  Perfdata
	Multiline []string
}

// SetExitCode sets the exit code only if it is greater than the current value
// (so you can set WARNING early in a function and you won't accidentally set it OK later).
func (r *Result) SetExitCode(exitCode ExitCode) {
	if exitCode.IsGreaterCode(r.ExitCode) {
		r.ExitCode = exitCode
	}
}

func (r *Result) SetExitCodeInt(exitCode int) {
	r.SetExitCode(ExitCode(exitCode))
}

// GetExitCode retrieves the exit code.
func (r *Result) GetExitCode() int {
	return r.ExitCode.Int()
}

// AddPerfData adds performance data to the string slice.
func (r *Result) AddPerfData(name, value string) {
	r.Perfdata = append(r.Perfdata, fmt.Sprintf("%s=%s", name, value))
}

// AddLine adds a line to multiline output.
func (r *Result) AddLine(line string) {
	r.Multiline = append(r.Multiline, line)
}

// SetText sets the main text message.
func (r *Result) SetText(msg string) {
	r.Text = msg
}

// Error uses an error to update a Result.
func (r *Result) Error(err error) {
	r.ExitCode = CRITICAL
	r.Text = err.Error()
}

// ErrorExit uses an error to update the Result and passes it to Exit().
func (r *Result) ErrorExit(err error) {
	r.ExitCode = CRITICAL
	r.Text = err.Error()
	Exit(r)
}

func Sprint(res *Result) (string, int) {
	var (
		output, text string
		exitCode     ExitCode
	)

	switch {
	case res.ExitCode == OK:
		text = "OK: " + res.Prefix + res.Text
		exitCode = res.ExitCode
	case res.ExitCode == WARNING:
		text = "WARNING: " + res.Prefix + res.Text
		exitCode = res.ExitCode
	case res.ExitCode == CRITICAL:
		text = "CRITICAL: " + res.Prefix + res.Text
		exitCode = res.ExitCode
	case res.ExitCode == UNKNOWN:
		text = "UNKNOWN: " + res.Prefix + res.Text
		exitCode = res.ExitCode
	default:
		text = "UNKNOWN: Exit code '" + strconv.Itoa(int(res.ExitCode)) + "' is not valid :" + res.Prefix + res.Text
		exitCode = 3
	}

	perfData := res.Perfdata.String()

	if len(res.Multiline) > 0 {
		multiline := ""
		for _, l := range res.Multiline {
			multiline = multiline + l + "\n"
		}

		output = fmt.Sprintf("%s|%s\n%s\n", text, perfData, multiline)
	} else {
		output = fmt.Sprintf("%s|%s\n", text, perfData)
	}

	return output, exitCode.Int()
}

// Render formats and prints the output, returning the exitcode to be passed to
// os.Exit
//
//nolint:forbidigo // returns output to console.
func Print(res *Result) int {
	output, exitCode := Sprint(res)
	fmt.Print(output)

	return exitCode
}

// Exit uses the Result struct to output Nagios plugin compatible output and exit codes.

func Exit(res *Result) {
	os.Exit(Print(res))
}
