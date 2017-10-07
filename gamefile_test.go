package main

import (
	"os"
	"testing"
)

var testFile = "data/anchor.z8"

const maxFileSize = 524288

func getData() []byte {
	gf, _ := os.Open(testFile)

	b := make([]byte, maxFileSize)
	_, _ = gf.Read(b)
	return b
}

func TestNewHeader(t *testing.T) {
	b := getData()
	h, err := newHeader(b)
	if err != nil {
		t.Error("Couldn't load file ", err)
	}
	expected := 37
	if len(h.data) != expected {
		t.Errorf("Header is the wrong size. Expected %d, got %d.", expected, len(h.data))
	}
}

func getHeader() (*Header, error) {
	return newHeader(getData())
}

func TestGetVersion(t *testing.T) {
	h, _ := getHeader()
	v := h.getVersion()
	expected := byte(8)
	if v != expected {
		t.Errorf("Wrong version. Got %x, expected %x.", v, expected)
	}
}

func TestGetHiMemBaseAddr(t *testing.T) {
	h, _ := getHeader()
	a := h.getHiMemAddr()
	expected := uint16(0x18fe)
	if a != expected {
		t.Errorf("Wrong Hi Mem Address. Got %x, expected %x.", a, expected)
	}
}

func TestGetObjTableAddr(t *testing.T) {
	h, _ := getHeader()
	a := h.getObjTableAddr()
	expected := uint16(0x0a01)
	if a != expected {
		t.Errorf("Wrong Object Table Address. Got %x, expected %x.", a, expected)
	}
}
