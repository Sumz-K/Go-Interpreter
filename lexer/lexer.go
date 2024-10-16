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

func (l* Lexer) peek() byte {
    if l.readPosition>=len(l.input) {
        return 0
    }
    return l.input[l.readPosition]
}

func (l* Lexer) readID() string{
    start:=l.position
    for l.isLetter() {
        l.readChar()
    }
    return l.input[start:l.position]
}

func (l* Lexer) readNumber() string {
    start:=l.position
    for l.isDigit() {
        l.readChar()
    }
    return l.input[start:l.position]
}

func (l* Lexer) isLetter() bool{
    return l.char>='a' && l.char<='z' || l.char>='A' && l.char<='Z'
}
func (l* Lexer) isDigit() bool {
    return l.char>='0' && l.char<='9'
}
func (l* Lexer) ignoreWhiteSpace() {
    for l.char==' ' || l.char=='\t' || l.char=='\n' || l.char=='\r'{
        l.readChar()
    }
}
func (l* Lexer) NextToken() token.Token {
    
    l.ignoreWhiteSpace()
    currChar:=l.char

    var tok token.Token

    switch currChar {
    case '=':
        if l.peek() == '=' {
            l.readChar()
            tok=token.Token{Type: token.EQ,Value: "=="}
        } else{
            tok=createToken(token.ASSIGN,l.char)
        }
    case '+':
        tok=createToken(token.PLUS,l.char)
    case '-':
        tok=createToken(token.MINUS,l.char)
    case '*':
        tok=createToken(token.ASTERISK,l.char)
    case '/':
        tok=createToken(token.SLASH,l.char)
    case '!':
        if l.peek() == '=' {
            l.readChar()
            tok=token.Token{Type: token.NOTEQ,Value: "!="}
        } else {
            tok=createToken(token.BANG,l.char)
        }
    case '<':
        tok=createToken(token.LT,l.char)
    case '>':
        tok=createToken(token.GT,l.char)
    case ';':
        tok=createToken(token.SEMICOLON,l.char)
    case ',':
        tok=createToken(token.COMMA,l.char)
    case '(':
        tok=createToken(token.LPAREN,l.char)
    case ')':
        tok=createToken(token.RPAREN,l.char)
    case '{':
        tok=createToken(token.LBRACE,l.char)
    case '}':
        tok=createToken(token.RBRACE,l.char)
    case 0:
        tok.Value=""
        tok.Type=token.EOF
    default:
        if l.isLetter() {
            tok.Value=l.readID()
            tok.Type=token.CheckID(tok.Value)
            return tok
        }  else if l.isDigit() {
            tok.Value=l.readNumber()
            tok.Type=token.INTEGER
            return tok
        } else {
            tok=createToken(token.ILLEGAL,l.char)
        }

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
