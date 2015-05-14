package main

import (
	"strings"
	"testing"
)

func TestNewScanner(t *testing.T) {
	s := NewScanner(strings.NewReader("Spring has come."))
	if s == nil {
		t.Fatal("Failed to create UniqScanner.")
	}
}

func TestScan(t *testing.T) {
	s := NewScanner(strings.NewReader("Lalala..."))
	if !s.Scan() {
		t.Error("Should advance the Scanner to the next token.")
	}
	if s.Scan() {
		t.Fatal("Should stop the scan.")
	}
}

func TestScanDuplicate(t *testing.T) {
	s := NewScanner(strings.NewReader("0\n0\n1\n"))
	if !s.ScanDuplicate() {
		t.Error("Should advance the Scanner to the next token.")
	}
	if s.ScanDuplicate() {
		t.Error("Should stop the scan.")
	}
}

func TestScanUniq(t *testing.T) {
	s := NewScanner(strings.NewReader("0\n1\n1\n"))
	if !s.ScanUnique() {
		t.Error("Should advance the Scanner to the next token.")
	}
	if s.ScanUnique() {
		t.Error("Should stop the scan.")
	}
}

func TestText(t *testing.T) {
	s := NewScanner(strings.NewReader("0\n0\n1\n"))
	s.Scan()
	if s.Text() != "0" {
		t.Errorf("Should scan while the same token. (%v)", s.Text())
	}
	s.Scan()
	if s.Text() != "1" {
		t.Errorf("Should scan until the unique token. (%v)", s.Text())
	}
}

func TestDuplicateText(t *testing.T) {
	var s *UniqScanner

	s = NewScanner(strings.NewReader("0\n1\n1\n"))
	s.ScanDuplicate()
	if s.Text() != "1" {
		t.Errorf("Should remove unique token. (%v)", s.Text())
	}
	if s.ScanDuplicate() {
		t.Error("Should stop the scan.")
	}

	s = NewScanner(strings.NewReader("0\n0\n1\n"))
	s.ScanDuplicate()
	if s.Text() != "0" {
		t.Errorf("Should remove unique token. (%v)", s.Text())
	}
	if s.ScanDuplicate() {
		t.Error("Should stop the scan.")
	}

	s = NewScanner(strings.NewReader("0\n"))
	if s.ScanDuplicate() {
		t.Errorf("Should skip all tokens.")
	}

	s = NewScanner(strings.NewReader("0\n1\n"))
	if s.ScanDuplicate() {
		t.Errorf("Should skip all tokens.")
	}
}

func TestUniqueText(t *testing.T) {
	var s *UniqScanner

	s = NewScanner(strings.NewReader("0\n0\n1\n"))
	s.ScanUnique()
	if s.Text() != "1" {
		t.Errorf("Should remove duplica token. (%v)", s.Text())
	}
	if s.ScanUnique() {
		t.Error("Should stop the scan.")
	}

	s = NewScanner(strings.NewReader("0\n1\n1\n"))
	s.ScanUnique()
	if s.Text() != "0" {
		t.Errorf("Should remove duplica token. (%v)", s.Text())
	}
	if s.ScanUnique() {
		t.Error("Should stop the scan.")
	}

	s = NewScanner(strings.NewReader("0\n"))
	s.ScanUnique()
	if s.Text() != "0" {
		t.Error("Should advance the Scanner to the next token.")
	}
	if s.ScanUnique() {
		t.Error("Should stop the scan.")
	}

	s = NewScanner(strings.NewReader("0\n1\n"))
	s.ScanUnique()
	if s.Text() != "0" {
		t.Error("Should advance the Scanner to the next token.")
	}
	s.ScanUnique()
	if s.Text() != "1" {
		t.Error("Should advance the Scanner to the next token.")
	}
	if s.ScanUnique() {
		t.Error("Should stop the scan.")
	}
}
