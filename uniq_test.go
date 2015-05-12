package main

import (
	"strings"
	"testing"
)

func TestNewScanner(t *testing.T) {
	u := NewScanner(strings.NewReader("Spring has come."))
	if u == nil {
		t.Fatal("Failed to create UniqScanner.")
	}
}

func TestScan(t *testing.T) {
	u := NewScanner(strings.NewReader("Lalala..."))
	if !u.Scan() {
		t.Error("Should advance the Scanner to the next token.")
	}
	if u.Scan() {
		t.Fatal("Should stop the scan.")
	}
}

func TestText(t *testing.T) {
	u := NewScanner(strings.NewReader("0\n0\n0\n0\n0\n1\n"))
	u.Scan()
	if u.Text() != "0" {
		t.Errorf("Should scan while the same token. (%v)", u.Text())
	}
	u.Scan()
	if u.Text() != "1" {
		t.Errorf("Should scan until the unique token. (%v)", u.Text())
	}
}
