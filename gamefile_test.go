package main

import (
	"os"
	"testing"
)

var testFile = "data/anchor.z8"

func testNewHeader(t *testing.T) {
	gf, err := os.Open(testFile)
	if err != nil {
		t.Error("Couldn't load file ", err)
	}

	var b []byte
	_, err = gf.Read(b)
	if err != nil {
		t.Error("Couldn't load file ", err)
	}

	h, _ := newHeader(b)
	if len(h.data) != 296 {
		t.Error("Header is the wrong size")
	}
}

func testGetVersion(t *testing.T) {
	gf, err := os.Open(testFile)
	if err != nil {
		t.Error("Couldn't load file ", err)
	}

	var b []byte
	_, err = gf.Read(b)
	if err != nil {
		t.Error("Couldn't load file ", err)
	}

	h, _ := newHeader(b)
	v := h.getVersion()
	expected := byte(8)
	if v != expected {
		t.Errorf("Wrong version. Got %x, expected %x.", v, expected)
	}
}
