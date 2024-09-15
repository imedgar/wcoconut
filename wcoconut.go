package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strings"
	"time"
	"unicode/utf8"
)

const invOpt = "invalid option. usage wcoconut [-c -l -w -m] file"

var (
	path      string
	data      []byte
	validOpts = []string{"-c", "-l", "-w", "-m"}
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		log.Fatal(invOpt)
	}

	path = args[0]
	if len(args) == 2 {
		path = args[1]
	}

	if isValidOpt(args[0]) {
		log.Fatal("invalid option ", args[0])
	}

	input := readStdin()
	if input != "" {
		data = []byte(input)
		path = ""
	} else {
		dat, err := os.ReadFile(path)
		if err != nil {
			log.Fatal(err)
		}
		data = dat
	}
	count(args[0])
}

func isValidOpt(opt string) bool {
	return len(opt) == 2 && strings.Contains(opt, "-") && !slices.Contains(validOpts, opt)
}

func count(opt string) {
	switch opt {
	case "-c":
		fmt.Println(bytesInFile(), path)
	case "-l":
		fmt.Println(linesInFile(), path)
	case "-w":
		fmt.Println(wordsInFile(), path)
	case "-m":
		fmt.Println(utf8.RuneCount(data), path)
	default:
		fmt.Println(linesInFile(), wordsInFile(), bytesInFile(), path)
	}
}

func bytesInFile() int {
	return len(data)
}

func linesInFile() int {
	return len(strings.Split(string(data), "\n")) - 1
}

func wordsInFile() int {
	return len(bytes.Fields(data))
}

func readStdin() string {
	var input strings.Builder
	reader := bufio.NewReader(os.Stdin)

	done := make(chan struct{})
	go func() {
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					log.Fatal(err)
				}
				break
			}
			input.WriteString(line)
		}
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(100 * time.Millisecond):
	}

	return strings.TrimSpace(input.String())
}
