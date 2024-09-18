package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	usage           = "usage: wcoconut [-c -l -w -m] file"
	timeoutDuration = 100 * time.Millisecond
)

var (
	errInvalidOption = errors.New("invalid option")
	validOpts        = []string{"-c", "-l", "-w", "-m"}
)

func main() {
	if err := run(os.Args[1:], os.Stdin, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string, stdin io.Reader, stdout io.Writer) error {
	if len(args) < 1 {
		return fmt.Errorf(usage)
	}

	opt, path, err := parseArgs(args)
	if err != nil {
		return err
	}

	data, err := readInput(path, stdin)
	if err != nil {
		return err
	}
	return count(opt, path, data, stdout)
}

func parseArgs(args []string) (string, string, error) {
	fArg := args[0]
	path := ""

	if len(args) == 2 {
		path = args[1]
	} else if !isValidOpt(fArg) {
		path = fArg
		fArg = ""
	}

	if isValidOpt(fArg) {
		return "", "", fmt.Errorf("%w: %s", errInvalidOption, fArg)
	}

	return fArg, path, nil
}

func isValidOpt(opt string) bool {
	return len(opt) == 2 && strings.HasPrefix(opt, "-") && !slices.Contains(validOpts, opt)
}

func readInput(path string, stdin io.Reader) ([]byte, error) {
	input, err := readStdin(stdin)
	if err != nil || len(input) == 0 {
		return os.ReadFile(path)
	}
	return input, nil
}

func count(opt, path string, data []byte, w io.Writer) error {
	out := path
	if slices.Contains(validOpts, out) {
		out = string(data)
	}
	switch opt {
	case "-c":
		fmt.Fprintf(w, "%d %s\n", len(data), out)
	case "-l":
		fmt.Fprintf(w, "%d %s\n", linesInFile(data), out)
	case "-w":
		fmt.Fprintf(w, "%d %s\n", wordsInFile(data), out)
	case "-m":
		fmt.Fprintf(w, "%d %s\n", utf8.RuneCount(data), out)
	default:
		fmt.Fprintf(w, "%d %d %d %s\n", linesInFile(data), wordsInFile(data), len(data), out)
	}
	return nil
}

func linesInFile(data []byte) int {
	return bytes.Count(data, []byte{'\n'})
}

func wordsInFile(data []byte) int {
	return len(bytes.Fields(data))
}

func readStdin(r io.Reader) ([]byte, error) {
	var input strings.Builder
	done := make(chan struct{})

	go func() {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			input.WriteString(scanner.Text() + "\n")
		}
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(timeoutDuration):
	}

	return []byte(strings.TrimSpace(input.String())), nil
}
