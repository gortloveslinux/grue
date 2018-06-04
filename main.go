package main

import (
	"io/ioutil"
	"log"
	"os"
)

func main() {
	fileName := os.Args[1]
	buf, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	logger := log.New(os.Stdout, "", 0)

	var zm ZMachine
	zm.Init(buf, logger.Printf)

	zm.InterpretInstruction()
}
