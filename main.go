package main

import (
	"log"
	"os"

	"github.com/Sumz-K/Go-Interpreter/repl"
)

func main() {
    file,err:=os.Open("monkey/code1.monkey")
    if err!=nil {
        log.Fatal("Could not open file")
    }
    defer file.Close()
    repl.Start(file)

}