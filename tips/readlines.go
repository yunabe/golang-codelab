package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"os"
	"strconv"
)

var (
	input = flag.String("input", "", "Input file path.")
)

// How to read lines from a file in Go.
func readWithScanner() {
	f, err := os.Open(*input)
	if err != nil {
	 	log.Fatal(err)
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		log.Print(strconv.Quote(s.Text()))
	}
	if s.Err() != nil {
		// non-EOF error.
		log.Fatal(s.Err())
	}
}

// You can use bufio.Reader though bufio.Scanner is simpler.
func readWithReader() {
	f, err := os.Open(*input)
	if err != nil {
	 	log.Fatal(err)
	}
	defer f.Close()
	r := bufio.NewReader(f)
	for {
		// line includes '\n'.
		line, err := r.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		log.Print(strconv.Quote(line))
	}
}

func readLinesFromStdin() {
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		log.Print(strconv.Quote(s.Text()))
	}
	if s.Err() != nil {
		// non-EOF error.
		log.Fatal(s.Err())
	}
}

func main() {
	flag.Parse()

	if len(*input) != 0 {
		// Read lines from a file.
		readWithScanner()
		readWithReader()
	} else {
		// Read lines from stdin.
		readLinesFromStdin()
	}
}
