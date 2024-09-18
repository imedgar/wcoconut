package main

import (
	"bytes"
	"io"
	"os"
	"testing"
)

const (
	testFileName = "iwtbyc.txt"
	fileContent  = `BLANK: Bereft of father! Bereft of mother! Marcus! Thou hast lost even thy love!
CINNA: Fortune hath escape'd thee! For what end shalt thou live?
ZIDANE: For the sake of our friend... Let us bury our steel in the heart of the wretched King Leo!
  `
	pipedInput = "this is a test"
)

func TestMain(m *testing.M) {
	err := os.WriteFile(testFileName, []byte(fileContent), 0644)
	if err != nil {
		panic("error creating test file: " + err.Error())
	}
	oldArgs := os.Args
	exitCode := m.Run()
	os.Args = oldArgs
	err = os.Remove(testFileName)
	if err != nil {
		panic("error deleting test file: " + err.Error())
	}
	os.Exit(exitCode)
}

func TestWcoconut(t *testing.T) {
	testCases := []struct {
		name     string
		expected string
		input    string
		args     []string
	}{
		{"Count bytes", output("247"), "", []string{"-c", testFileName}},
		{"Count lines", output("3"), "", []string{"-l", testFileName}},
		{"Count words", output("45"), "", []string{"-w", testFileName}},
		{"Count characters", output("247"), "", []string{"-m", testFileName}},
		{"Default count", output("3 45 247"), "", []string{testFileName}},
		{"Piped input with -c", "0 4 14 " + pipedInput + "\n", pipedInput, []string{"-c"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := runWithInput(tc.args, tc.input)
			if err != nil {
				t.Fatalf("Error running test: %v", err)
			}
			if output != tc.expected {
				t.Errorf("case %s with args %v, expected %q, but got %q", tc.name, tc.args, tc.expected, output)
			}
		})
	}
}

func output(out string) string {
	return out + " " + testFileName + "\n"
}

func runWithInput(args []string, input string) (string, error) {
	oldArgs := os.Args
	oldStdin := os.Stdin
	oldStdout := os.Stdout

	os.Args = append([]string{"cmd"}, args...)

	ir, iw, err := os.Pipe()
	if err != nil {
		return "", err
	}
	or, ow, err := os.Pipe()
	if err != nil {
		return "", err
	}

	os.Stdin = ir
	os.Stdout = ow

	go func() {
		defer iw.Close()
		io.WriteString(iw, input)
	}()

	go func() {
		defer ow.Close()
		main()
	}()

	var buf bytes.Buffer
	io.Copy(&buf, or)

	os.Args = oldArgs
	os.Stdin = oldStdin
	os.Stdout = oldStdout

	return buf.String(), nil
}
