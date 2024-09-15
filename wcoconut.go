package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

const invOpt = "invalid option. usage wcoconut [-c -l -w -m] file"

var (
	path string
	data []byte
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

func count(opt string) {
	switch opt {
	case "-c":
		fmt.Println(bytesInFile(), path)
	case "-l":
		fmt.Println(linesInFile(), path)
	case "-w":
		fmt.Println(wordsInFile(), path)
	case "-m":
		fmt.Println(bytes.Count(data, []byte(""))-1, path)
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
