package main

import (
	"os"
	"testing"
)

var testFile = "data/anchor.z8"

const maxFileSize = 524288

func testGetVersion(t *testing.T) {
	h := setupHeader(t, testFile)
	v := h.getVersion()
	expected := uint8(8)
	if v != expected {
		t.Errorf("Wrong version. Got %x, expected %x.", v, expected)
	}
}

func testGetHiMemBase(t *testing.T) {
	h := setupHeader(t, testFile)
	a := h.getHiMemBase()
	expected := uint16(0xFE18)
	if a != expected {
		t.Errorf("Wrong Hi Mem Base. Got %x, expected %x.", a, expected)
	}

}

func testGetInitialPC(t *testing.T) {
	h := setupHeader(t, testFile)
	a := h.getHiMemBase()
	expected := uint16(0xFE19)
	if a != expected {
		t.Errorf("Wrong Initial PC. Got %x, expected %x.", a, expected)
	}

}

func setupHeader(t *testing.T, filename string) *Header {
	gf, err := os.Open(filename)
	if err != nil {
		t.Fatal("Couldn't load file ", err)
	}
	var b []uint8
	_, err = gf.Read(b)
	if err != nil {
		t.Fatal("Couldn't read file ", err)
	}
	h, err := newHeader(b)
	if err != nil {
		t.Fatal("Couldn't create header ", err)
	}

	return h
}
