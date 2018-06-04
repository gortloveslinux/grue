package main

import (
	"io/ioutil"
	"testing"
)

var testFile = "data/anchor.z8"

func setupMachine(t *testing.T, filename string) ZMachine {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	var zm ZMachine
	zm.Init(buf)
	return zm
}

func TestVerify(t *testing.T) {
	zm := setupMachine(t, testFile)

	zm.verify()
	//TODO more
}
