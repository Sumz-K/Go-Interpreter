package parser

import (
	"fmt"
	"strconv"

	"github.com/Sumz-K/Go-Interpreter/ast"
	"github.com/Sumz-K/Go-Interpreter/token"
)

// for operator precedence ig
const (
	_ int = iota
	LOWEST
	EQUALS // == LESSGREATER // > or <
	SUM //+
	PRODUCT //*
	PREFIX //-Xor!X
	CALL // myFunction(X)
	)
	

// In Pratt parsing, each token is associated with some parsing functions, in our case infix and prefix
// Prefix parsing refers to cases where the expressions has to be parsed conisdering the token to preceed some expression like --5
// Infix parsing refers to cases where the token/operator is present in between some other tokens like 5*5

type (
	prefixParseFn func() ast.Expression // in prefix cases nothing preceeds the token so no parameter
	infixParseFn func(ast.Expression) ast.Expression  // in infix cases the tokens to the left of the operator need to be passed in as a parameter
)

func (p* Parser) registerInfixFunc(tt token.TokenType, fn infixParseFn) {
	p.infixFunc[tt]=fn
}

func (p* Parser) registerPrefixFunc(tt token.TokenType, fn prefixParseFn) {
	p.prefixFunc[tt]=fn
}

func (p* Parser) parseExpressionStmt() *ast.ExpressionStmt {
	stmt:=&ast.ExpressionStmt{}
	stmt.Token=p.currToken

	stmt.Expression=p.parseExpression(LOWEST)

	if p.isNext(token.SEMICOLON) { // the semicolon is optional in the langauge 
		p.next()
	}
	return stmt 
}

func (p* Parser) isNext(tok token.TokenType) bool {
	return p.peekToken.Type==tok 
}

func (p* Parser) noPrefixFuncErr(tok token.TokenType) {
	msg:=fmt.Sprintf("There exists no prefix parse function for token %s",tok)
	p.errors = append(p.errors, msg)
}


func (p* Parser) parseExpression(precedence int) ast.Expression {
	prefix:=p.prefixFunc[p.currToken.Type] //check if a function corresponding to when the token is in prefix position
	if prefix==nil {
		p.noPrefixFuncErr(p.currToken.Type)
		return nil 
	}

	leftExpr:=prefix() // if the function is present call the function and obtain the expression
	return  leftExpr

}



func (p *Parser) parseIdentifier() ast.Expression {
	ident := &ast.Identifier{Token: p.currToken, Value: p.currToken.Value}
	return ident 
}


func (p* Parser) parseIntLiteral() ast.Expression {
	il:=&ast.IntegerLiteral{}
	il.Token=p.currToken

	intVal,err:=strconv.Atoi(p.currToken.Value)
	if err!=nil {
		fmt.Printf("Error converting %q string to integer",p.currToken.Value)
		p.addError(token.INTEGER)
		return nil 
	}
	il.Value=int64(intVal)

	return il 

}


// actually parses exprs of type -5 and !5 etc
func (p* Parser) parsePrefixExpression() ast.Expression{
	expr:=&ast.PrefixExpression{}
	expr.Token=p.currToken
	expr.Operator=p.currToken.Value

	p.next()

	expr.Right=p.parseExpression(PREFIX)
	return expr 
}