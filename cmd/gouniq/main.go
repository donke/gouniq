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
	c = flag.Bool("c", false, "")
	d = flag.Bool("d", false, "")
	u = flag.Bool("u", false, "")
	i = flag.Bool("i", false, "")

	f = flag.Int("f", 0, "")
	s = flag.Int("s", 0, "")

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

func skip(str string, fields int, chars int) string {
	result := str
	for i := 0; i < fields; i++ {
		result = result[strings.IndexAny(result, " \t")+1:]
	}
	if chars != 0 {
		if chars >= utf8.RuneCountInString(result) {
			return result
		}
		ru := []rune(result)
		return string(ru[chars:])
	}
	return result
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

	if *c && (*d || *u) {
		usageAndExit()
	}
	if *d && *u {
		usageAndExit()
	}

	scanner := gouniq.NewScanner(reader)
	scanner.Equal(func(s1, s2 string) bool {
		if *i {
			return strings.ToLower(skip(s1, *f, *s)) == strings.ToLower(skip(s2, *f, *s))
		}
		return skip(s1, *f, *s) == skip(s2, *f, *s)
	})

	switch {
	case *c:
		scanner.ScanFunc(scanner.ScanCount)
	case *d:
		scanner.ScanFunc(scanner.ScanDuplicate)
	case *u:
		scanner.ScanFunc(scanner.ScanUnique)
	default:
		scanner.ScanFunc(scanner.ScanOriginal)
	}

	for scanner.Scan() {
		if *c {
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
