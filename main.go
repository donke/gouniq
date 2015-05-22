package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
)

var (
	count       = flag.Bool("c", false, "print a number that how many times they occured.")
	duplicate   = flag.Bool("d", false, "only print duplicated lines.")
	unique      = flag.Bool("u", false, "only print unique lines.")
	insensitive = flag.Bool("i", false, "case insensitive comparison.")
	fields      = flag.Int("f", 0, "ignore the first [num] fields.")
	chars       = flag.Int("s", 0, "ignore the first [num] characters")
	reader      = os.Stdin
	writer      = os.Stdout
)

const (
	help = "usage: gouniq [-c | -d | -u] [-i] [-f fields] [-s chars] [input [output]]"
)

func usage() {
	fmt.Fprintln(os.Stderr, help)
	os.Exit(1)
}

func skip(s string) string {
	r := s
	for i := 0; i < *fields; i++ {
		r = r[strings.IndexAny(r, " \t")+1:]
	}
	if *chars != 0 {
		if *chars >= utf8.RuneCountInString(r) {
			return r
		}
		ru := []rune(r)
		return string(ru[*chars:])
	}
	return r
}

func main() {
	flag.Parse()

	if flag.NArg() > 2 {
		usage()
	}

	if flag.NArg() > 1 {
		w, err := os.Create(flag.Args()[1])
		if err != nil {
			fmt.Fprintln(os.Stderr, flag.Args()[1]+":", err)
			os.Exit(1)
		}
		writer = w
		defer writer.Close()
	}

	if flag.NArg() > 0 {
		r, err := os.Open(flag.Args()[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, flag.Args()[0]+":", err)
			os.Exit(1)
		}
		reader = r
		defer reader.Close()
	}

	if *count && (*duplicate || *unique) {
		usage()
	}
	if *duplicate && *unique {
		usage()
	}

	scanner := NewScanner(reader)
	scanner.Equal(func(s1, s2 string) bool {
		if *insensitive {
			return strings.ToLower(skip(s1)) == strings.ToLower(skip(s2))
		}
		return skip(s1) == skip(s2)
	})

	switch {
	case *count:
		scanner.ScanFunc(scanner.ScanCount)
	case *duplicate:
		scanner.ScanFunc(scanner.ScanDuplicate)
	case *unique:
		scanner.ScanFunc(scanner.ScanUnique)
	default:
		scanner.ScanFunc(scanner.ScanOriginal)
	}

	for scanner.Scan() {
		if *count {
			fmt.Fprintf(writer, "%4d %s\n", scanner.Count(), scanner.Text())
		} else {
			fmt.Fprintln(writer, scanner.Text())
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
