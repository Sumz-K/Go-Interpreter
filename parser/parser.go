package parser

import (
	"fmt"

	"github.com/Sumz-K/Go-Interpreter/ast"
	"github.com/Sumz-K/Go-Interpreter/lexer"
	"github.com/Sumz-K/Go-Interpreter/token"
)


type Parser struct {
    l *lexer.Lexer
    currToken token.Token
    peekToken token.Token
    errors []string
}

func New(l *lexer.Lexer) *Parser {
    p:=&Parser{
        l:l,
        errors: []string{},
    }
    //Read two tokens to set the current and peek tokens
    p.next()
    p.next()

    return p
}

func (p *Parser) next() {
    p.currToken=p.peekToken
    p.peekToken=p.l.NextToken()
}

func (p* Parser) showErrors() []string {
    return p.errors
}

func (p *Parser) addError(t token.TokenType) {
    msg:=fmt.Sprintf("expected next token to be %s, got %s instead",t, p.peekToken.Type)

    p.errors = append(p.errors,msg)
}


func (p *Parser) ParseProgram() *ast.Program {
    program:=&ast.Program{}
    program.Statements=[]ast.Statement{}

    for p.currToken.Type!=token.EOF {
        statement:=p.parseStmt()
        if statement!=nil {
            program.Statements = append(program.Statements, statement)
        }
        p.next()
    }
    return program
}




func (p* Parser) parseStmt() ast.Statement {
    switch p.currToken.Type {
        case token.LET:
            return p.parseLetStmt()
        default:
            return nil
    }
}


func (p* Parser) parseLetStmt() ast.Statement {
    stmt:=&ast.LetStmt{}

    stmt.Token=p.currToken //token.LET

    if !p.expected(token.IDENTIFIER) {
        return nil //nect token has to be an ID``
    }
    stmt.Name=&ast.Identifier{
        Token: p.currToken, //IDENTIFIER token
        Value: p.currToken.Value, //name of that variable
    }

   if !p.expected(token.ASSIGN) { //check if next token is "="
        return nil
    }

    for !p.isCurr(token.SEMICOLON) {
        p.next()
    }


    return stmt


}

func (p *Parser) isCurr(tok token.TokenType) bool {
    return p.currToken.Type==tok
}



func (p *Parser) expected(tok token.TokenType) bool {
    if p.peekToken.Type==tok {
        p.next()
        return true 
    }
    p.addError(tok)
    return false 
}
