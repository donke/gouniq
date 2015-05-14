package main

import (
	"bufio"
	"io"
)

// UniqScanner ...
type UniqScanner struct {
	scanner *bufio.Scanner
	first   bool
	last    bool
	prev    string
	this    string
	text    string
	repeats int
}

// NewScanner ...
func NewScanner(r io.Reader) *UniqScanner {
	scanner := &UniqScanner{scanner: bufio.NewScanner(r), first: true, last: true}
	return scanner
}

// Scan ...
func (u *UniqScanner) Scan() bool {
	if u.first {
		u.first = false
		if u.scanner.Scan() {
			u.prev = u.scanner.Text()
			u.text = u.prev
			return true
		}
		return false
	}

	for u.scanner.Scan() {
		u.this = u.scanner.Text()
		if u.prev != u.this {
			u.prev = u.this
			u.text = u.prev
			return true
		}
		u.repeats++
	}

	return false
}

// ScanDuplicate ...
func (u *UniqScanner) ScanDuplicate() bool {
	if u.first {
		u.first = false
		if u.scanner.Scan() {
			u.prev = u.scanner.Text()
		}
	}

	for u.scanner.Scan() {
		u.this = u.scanner.Text()
		if u.prev != u.this {
			if u.repeats != 0 {
				u.text = u.prev
				u.prev = u.this
				u.repeats = 0
				return true
			}
			u.prev = u.this
			u.repeats = 0
		} else {
			u.repeats++
		}
	}

	if u.last {
		u.last = false
		if u.repeats != 0 {
			u.text = u.prev
			return true
		}
	}

	return false
}

// ScanUnique ...
func (u *UniqScanner) ScanUnique() bool {
	if u.first {
		u.first = false
		if u.scanner.Scan() {
			u.prev = u.scanner.Text()
		}
	}

	for u.scanner.Scan() {
		u.this = u.scanner.Text()
		if u.prev != u.this {
			if u.repeats == 0 {
				u.text = u.prev
				u.prev = u.this
				return true
			}
			u.prev = u.this
			u.repeats = 0
		} else {
			u.repeats++
		}
	}

	if u.last {
		u.last = false
		if u.repeats == 0 {
			u.text = u.prev
			return true
		}
	}

	return false
}

// Text ...
func (u *UniqScanner) Text() string {
	return u.text
}
