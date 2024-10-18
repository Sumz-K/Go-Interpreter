package repl

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/Sumz-K/Go-Interpreter/lexer"
	"github.com/Sumz-K/Go-Interpreter/parser"
	//"github.com/Sumz-K/Go-Interpreter/token"
	// "github.com/Sumz-K/Go-Interpreter/ast"
)
func Start(in *os.File, out io.Writer) { 
    fmt.Print("Starting parser\n")
    scanner := bufio.NewScanner(in)
    for { 
        scanned := scanner.Scan() 
        if !scanned {
            return
        }
        line := scanner.Text()
        l := lexer.New(line)
        p := parser.New(l)
        program := p.ParseProgram()
        if len(p.ShowErrors()) != 0 {
            printParserErrors(out, p.ShowErrors())
            continue
        }
        io.WriteString(out, program.String())
        io.WriteString(out, "\n")
    }
}

func printParserErrors(out io.Writer, errors []string) { 
    for _, msg := range errors {
        io.WriteString(out, "\t"+msg+"\n")
    }
}

    