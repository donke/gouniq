package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

    "github.com/donke/gouniq"
)

var (
	count       = flag.Bool("c", false, "")
	duplicate   = flag.Bool("d", false, "")
	unique      = flag.Bool("u", false, "")
	insensitive = flag.Bool("i", false, "")

	fields = flag.Int("f", 0, "")
	chars  = flag.Int("s", 0, "")

	reader = os.Stdin
	writer = os.Stdout
)

var usage = `usage: gouniq [options...] [input [output]]

Options:
  -c       Print a number that how may times they occured.
  -d       Print only duplicated lines.
  -u       Print only unique lines.
  -i       Case insensitive comparison on line.
  -f num   Ignore the first num fields.
  -s chars Ignore the first chars characters.
`

func usageAndExit() {
	flag.Usage()
	fmt.Fprintf(os.Stderr, "\n")
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
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usage)
	}

	flag.Parse()
	if flag.NArg() > 2 {
		usageAndExit()
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
		usageAndExit()
	}
	if *duplicate && *unique {
		usageAndExit()
	}

	scanner := gouniq.NewScanner(reader)
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
