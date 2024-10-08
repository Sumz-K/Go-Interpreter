package lexer

import (
	"testing"

	"github.com/Sumz-K/Go-Interpreter/token"
)


func TestToken(t *testing.T) {
    input := `=+(){},;`
    tests:=[]struct{
        expectedType token.TokenType
        expectedValue string
    }{
        {token.ASSIGN,"="},
        {token.ADD,"+"},
        {token.LPARA,"("},
        {token.RPARA,")"},
        {token.LBRACE,"{"},
        {token.RBRACE,"}"},
        {token.COMMA,","},
        {token.SEMICOLON,";"},
        {token.EOF,""},
    }

    l:=New(input)

    for i,test:=range tests {
        tok:=l.NextToken()

        if tok.Type!= test.expectedType {
            t.Fatalf("The type of the token %d is wrong, expected%q got%q",i,test.expectedType,tok.Type) 
        }

        if tok.Value!=test.expectedValue {
            t.Fatalf("The value of the token %d is wrong, expected%q got%q",i,test.expectedValue,tok.Value)
        }

        
    }
}






