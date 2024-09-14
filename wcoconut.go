package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) < 2 {
		log.Fatal("no file or command specified")
	}
	if args[0] != "-c" {
		log.Fatal("invalid command ", args[0])
	}
	dat, err := os.ReadFile(args[1])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(dat), args[1])
}
