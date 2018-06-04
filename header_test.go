package main

import (
	"io/ioutil"
	"testing"
)

func setupHeader(t *testing.T, filename string) ZHeader {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	var zh ZHeader
	zh.load(buf)
	return zh
}

func TestLoadHeader(t *testing.T) {
	zh := setupHeader(t, testFile)
	t.Logf("Header: %x", zh)
}
