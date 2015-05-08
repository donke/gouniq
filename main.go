package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	var err error

	r := os.Stdin
	w := os.Stdout

	c := flag.Bool("c", false, "print a number that how many times they occurred.")
	d := flag.Bool("d", false, "only print duplicated lines.")
	i := flag.Bool("i", false, "case insensitive comparison.")
	u := flag.Bool("u", false, "only print unique lines.")

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

	if *c {
		if *d || *u {
			usage()
		}
	} else if !*d && !*u {
		*d = true
		*u = true
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
	} else {
		os.Exit(0)
	}

	if !*c && *d && *u {
		show(w, prevLine, *c, *d, *u, repeats)
	}
	for scanner.Scan() {
		thisLine = scanner.Text()
		if compare(prevLine, thisLine, *i) {
			if *c || !*d || !*u {
				show(w, prevLine, *c, *d, *u, repeats)
			}
			prevLine = thisLine
			if !*c && *d && *u {
				show(w, prevLine, *c, *d, *u, repeats)
			}
			repeats = 0
		} else {
			repeats++
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	if *c || !*d || !*u {
		show(w, prevLine, *c, *d, *u, repeats)
	}
}

func compare(s1 string, s2 string, i bool) bool {
	if i {
		return strings.ToLower(s1) != strings.ToLower(s2)
	}
	return s1 != s2
}

func show(w io.Writer, s string, c bool, d bool, u bool, repeats int) {
	if c {
		fmt.Fprintf(w, "%4d %s\n", repeats+1, s)
	}
	if (d && repeats != 0) || (u && repeats == 0) {
		fmt.Fprintln(w, s)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: gouniq [-c | -d | -u] [-i] [input [output]]")
	os.Exit(1)
}
