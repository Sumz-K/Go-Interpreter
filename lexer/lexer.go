package lexer

import (
	"github.com/Sumz-K/Go-Interpreter/token"
    
)

type Lexer struct {
    input string  //input to lex
    position int  //where we read from before  position
    readPosition int  //nextPosition to "peek"
    char byte //represents the character at the current position
}

// To return a lexer type given an input 
func New(input string) *Lexer {
    l:=&Lexer{
        input: input,
    }
    l.readChar()
    return l
}


func (l *Lexer) readChar() {
    if l.readPosition>=len(l.input) {
        l.char=0
    } else {
        l.char=l.input[l.readPosition]
    }

    l.position=l.readPosition
    l.readPosition+=1
}


func (l* Lexer) NextToken() token.Token {
    currChar:=l.char

    var tok token.Token

    switch currChar {
    case '=':
        tok=createToken(token.ASSIGN,l.char)
    case '+':
        tok=createToken(token.ADD,l.char)
    case ';':
        tok=createToken(token.SEMICOLON,l.char)
    case ',':
        tok=createToken(token.COMMA,l.char)
    case '(':
        tok=createToken(token.LPARA,l.char)
    case ')':
        tok=createToken(token.RPARA,l.char)
    case '{':
        tok=createToken(token.LBRACE,l.char)
    case '}':
        tok=createToken(token.RBRACE,l.char)
    case 0:
        tok.Value=""
        tok.Type=token.EOF
    }
    
    l.readChar()
    return tok
}


func createToken(tokenType token.TokenType, ch byte) token.Token {
    return token.Token{
        Type: tokenType,
        Value: string(ch),
    }
}
