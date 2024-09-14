package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
)

const invOpt = "invalid option. usage wcoconut [-c -l] file"

var path string

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		log.Fatal(invOpt)
	}
	path = args[0]
	if len(args) == 2 {
		path = args[1]
	}
	dat, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	switch opt := args[0]; opt {
	case "-c":
		fmt.Println(bytesInFile(dat), path)
	case "-l":
		fmt.Println(linesInFile(dat), path)
	case "-w":
		fmt.Println(wordsInFile(dat), path)
	case "-m":
		fmt.Println(bytes.Count(dat, []byte(""))-1, path)
	default:
		fmt.Println(linesInFile(dat), wordsInFile(dat), bytesInFile(dat), path)
	}
}

func bytesInFile(dat []byte) int {
	return len(dat)
}

func linesInFile(dat []byte) int {
	return len(strings.Split(string(dat), "\n")) - 1
}

func wordsInFile(dat []byte) int {
	return len(bytes.Fields(dat))
}
