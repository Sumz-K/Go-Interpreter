package ast

import (
	"github.com/Sumz-K/Go-Interpreter/token"
	"testing"
)


func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStmt{
				Token: token.Token{Type: token.LET, Value: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENTIFIER, Value: "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENTIFIER, Value: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}
	if program.String() != "let myVar = anotherVar;" { 
		t.Errorf("program.String() wrong. got=%q", program.String())
	}		
}