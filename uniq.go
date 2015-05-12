package main

import (
	"bufio"
	"io"
)

// UniqScanner ...
type UniqScanner struct {
	scanner *bufio.Scanner
	prev    string
	first   bool
}

// NewScanner ...
func NewScanner(r io.Reader) *UniqScanner {
	scanner := &UniqScanner{scanner: bufio.NewScanner(r), first: true}
	return scanner
}

// Scan ...
func (u *UniqScanner) Scan() bool {
	if u.first {
		if u.scanner.Scan() {
			u.prev = u.scanner.Text()
			u.first = false
			return true
		}
	} else {
		for u.scanner.Scan() {
			this := u.scanner.Text()
			if u.prev != this {
				u.prev = this
				return true
			}
		}
		return false
	}
	return false
}

// Text ...
func (u *UniqScanner) Text() string {
	return u.prev
}
