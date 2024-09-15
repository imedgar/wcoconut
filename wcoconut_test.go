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
		expected string
		args     []string
	}{
		{output("247"), []string{"cmd", "-c", testFileName}},
		{output("3"), []string{"cmd", "-l", testFileName}},
		{output("45"), []string{"cmd", "-w", testFileName}},
		{output("247"), []string{"cmd", "-m", testFileName}},
		{output("3 45 247"), []string{"cmd", testFileName}},
	}

	for _, tc := range testCases {
		os.Args = tc.args

		output := captureOutput(main)
		if output != tc.expected {
			t.Errorf("case args %v, expected %q, but got %q", tc.args, tc.expected, output)
		}
	}
}

func output(out string) string {
	return out + " " + testFileName + "\n"
}

func captureOutput(f func()) string {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}
