package repl

import (
	"bufio"
	"fmt"
	"os"

	"github.com/Sumz-K/Go-Interpreter/lexer"
	"github.com/Sumz-K/Go-Interpreter/token"
)
func Start(file *os.File) {
    scanner:=bufio.NewScanner(file)
    for scanner.Scan() {
        line:=scanner.Text()
        l:=lexer.New(line)

        for tok:=l.NextToken();tok.Type!=token.EOF;tok=l.NextToken() {
            fmt.Printf("%v\n",tok)
        }
    }
}
