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

	c := flag.Bool("c", false, "print a number that how many times they occurred.")

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
	var repeats int

	scanner := bufio.NewScanner(r)
	if scanner.Scan() {
		prevLine = scanner.Text()
		if err = scanner.Err(); err != nil {
			panic(err)
		}
	}
	if !*c {
		fmt.Fprintln(w, prevLine)
	}
	for scanner.Scan() {
		thisLine = scanner.Text()
		if prevLine != thisLine {
			if *c {
				fmt.Fprintf(w, "%4d %s\n", repeats+1, prevLine)
			}
			prevLine = thisLine
			if !*c {
				fmt.Fprintln(w, prevLine)
			}
			repeats = 0
		} else {
			repeats++
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	if len(prevLine) != 0 {
		if *c {
			fmt.Fprintf(w, "%4d %s\n", repeats+1, prevLine)
		}
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: gouniq [-c] [input [output]]")
	os.Exit(1)
}
