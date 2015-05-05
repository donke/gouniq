package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

func main() {
	var err error

	r := os.Stdin
	w := os.Stdout

    flag.Parse()

	if flag.NArg() > 2 {
		usage()
	}

	if flag.NArg() > 1 {
		w, err = os.Create(flag.Args()[1])
		if err != nil {
			panic(err)
		}
		defer w.Close()
	}

	if flag.NArg() > 0 {
		r, err = os.Open(flag.Args()[0])
		if err != nil {
			panic(err)
		}
		defer r.Close()
	}

	var prevLine string
	var thisLine string

	scanner := bufio.NewScanner(r)
	if scanner.Scan() {
		prevLine = scanner.Text()
		if err = scanner.Err(); err != nil {
			panic(err)
		}
	}
	for scanner.Scan() {
		thisLine = scanner.Text()
		if prevLine != thisLine {
            fmt.Fprintln(w, prevLine)
			prevLine = thisLine
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	if len(prevLine) != 0 {
        fmt.Fprintln(w, prevLine)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: gouniq [input [output]]")
	os.Exit(1)
}
