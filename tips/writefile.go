package main

import (
	"bufio"
	"flag"
	"io/ioutil"
	"log"
	"os"
)

var (
	output    = flag.String("output", "", "Output file path.")
	useIOUtil = flag.Bool("use_ioutil", false, "Use ioutil.WriteFile")
)

const message = "Hello Golang!\n"

func main() {
	flag.Parse()

	var err error
	var f *os.File
	if len(*output) != 0 {
		// Write to a file.
		f, err = os.Create(*output)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
	} else {
		// Write to Stdout.
		f = os.Stdout
	}

	if !*useIOUtil {
		w := bufio.NewWriter(f)
		// Don't forget to call Flush!
		defer w.Flush()
		// You can ignore the first return value n (the number of bytes written) because
		// err != nil when n != len(message).
		_, err = w.WriteString(message)
	} else {
		err = ioutil.WriteFile(*output, []byte(message), 0666)
	}
	if err != nil {
		log.Fatal(err)
	}
}
