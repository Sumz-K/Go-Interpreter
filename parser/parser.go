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

    prefixFunc map[token.TokenType]prefixParseFn
    infixFunc map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
    p:=&Parser{
        l:l,
        errors: []string{},
    }

    // initialise the prefix map and register a function to parse ids 
    p.prefixFunc=make(map[token.TokenType]prefixParseFn)
    p.registerPrefixFunc(token.IDENTIFIER,p.parseIdentifier)
    p.registerPrefixFunc(token.INTEGER,p.parseIntLiteral)
    p.registerPrefixFunc(token.MINUS,p.parsePrefixExpression)
    p.registerPrefixFunc(token.BANG,p.parsePrefixExpression)
    p.registerPrefixFunc(token.TRUE,p.parseBoolean)
    p.registerPrefixFunc(token.FALSE,p.parseBoolean)
    p.registerPrefixFunc(token.LPAREN,p.parseGrouped)
    p.registerPrefixFunc(token.IF,p.parseIfExpression)
    p.registerPrefixFunc(token.FUNC,p.parseFunction)

    p.infixFunc=make(map[token.TokenType]infixParseFn)
    p.registerInfixFunc(token.PLUS,p.parseInfixExpression)
    p.registerInfixFunc(token.MINUS,p.parseInfixExpression)
    p.registerInfixFunc(token.ASTERISK,p.parseInfixExpression)
    p.registerInfixFunc(token.SLASH,p.parseInfixExpression)
    p.registerInfixFunc(token.GT,p.parseInfixExpression)
    p.registerInfixFunc(token.LT,p.parseInfixExpression)
    p.registerInfixFunc(token.EQ,p.parseInfixExpression)
    p.registerInfixFunc(token.NOTEQ,p.parseInfixExpression)
    p.registerInfixFunc(token.LPAREN,p.parseCallExpression)
    //Read two tokens to set the current and peek tokens
    p.next()
    p.next()

    return p
}

func (p *Parser) next() {
    p.currToken=p.peekToken
    p.peekToken=p.l.NextToken()
}

func (p* Parser) ShowErrors() []string {
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
        case token.RETURN:
            return p.parseReturnStmt()
        default:
            return p.parseExpressionStmt()
    }
}


func (p* Parser) parseReturnStmt() ast.Statement {
    stmt:=&ast.ReturnStmt{}
    stmt.Token=p.currToken // the token.RETURN

    p.next()

    stmt.ReturnValue=p.parseExpression(LOWEST)

    if p.isNext(token.SEMICOLON) {
        p.next()
    }

    return stmt 
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

    p.next()

    stmt.Value=p.parseExpression(LOWEST)
    if p.isNext(token.SEMICOLON) {
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
