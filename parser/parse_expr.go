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
	EQUALS // ==
	LESSGREATER // > or <
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
var precedences = map[token.TokenType]int {
	token.EQ: EQUALS,
	token.NOTEQ:EQUALS,
	token.LT:LESSGREATER,
	token.GT:LESSGREATER,
	token.PLUS:SUM,
	token.MINUS:SUM,
	token.ASTERISK:PRODUCT,
	token.SLASH:PRODUCT,
	token.LPAREN:CALL,

}

func (p *Parser) currPrecedence() int {
	ans,ok:=precedences[p.currToken.Type]
	if ok {
		return ans
	}
	return LOWEST
}

func (p *Parser) peekPrecedence() int {
	ans,ok:=precedences[p.peekToken.Type]
	if ok {
		return ans
	}
	return LOWEST
}


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

func (p* Parser) parseExpression(precedence int) ast.Expression {
	prefix:=p.prefixFunc[p.currToken.Type] //check if a function corresponding to when the token is in prefix position
	// 5*5
	if prefix==nil {
		p.noPrefixFuncErr(p.currToken.Type)
		return nil 
	}

	leftExpr:=prefix() // if the function is present call the function and obtain the expression

	for p.peekToken.Type!=token.SEMICOLON && precedence<p.peekPrecedence(){
		infix:=p.infixFunc[p.peekToken.Type]
		if infix==nil {
			return leftExpr
		}
		p.next()
		leftExpr=infix(leftExpr)
	}
	return  leftExpr

}


func (p* Parser) isNext(tok token.TokenType) bool {
	return p.peekToken.Type==tok 
}

func (p* Parser) noPrefixFuncErr(tok token.TokenType) {
	msg:=fmt.Sprintf("There exists no prefix parse function for token %s",tok)
	p.errors = append(p.errors, msg)
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



func (p* Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	ie:=&ast.InfixExpression{
		Token: p.currToken,
		Operator: p.currToken.Value,
		LeftExpr: left,
	}
	prec:=p.currPrecedence()
	p.next()
	ie.RightExpr=p.parseExpression(prec)
	return ie
}

func (p* Parser) parseBoolean() ast.Expression{
	expr:=&ast.Boolean{}
	expr.Token=p.currToken
	if expr.Token.Value=="true"{
		expr.Value=true 
	} else {
		expr.Value=false
	}
	return expr
}


// (1+2)*3
func (p* Parser) parseGrouped() ast.Expression {
	p.next()
	expr:=p.parseExpression(LOWEST)

	if !p.expected(token.RPAREN) {
		return nil 
	}
	return expr 
}

// if (a>b) {x} else {y}
func(p *Parser) parseIfExpression() ast.Expression {
	expr:=&ast.IfExpression{}
	expr.Token=p.currToken
	if !p.expected(token.LPAREN) {
		return nil 
	}
	p.next()
	expr.Condition=p.parseExpression(LOWEST)
	if !p.expected(token.RPAREN) {
		return nil 
	}

	if !p.expected(token.LBRACE) {
		return nil 
	}

	expr.Consequence=p.parseBlock()

	if p.isNext(token.ELSE) {
		p.next()
		if !p.expected(token.LBRACE) {
			return nil 
		}
		expr.Alternative=p.parseBlock()
	}

	return expr
}

func(p *Parser) parseBlock() *ast.BlockStmt {
	block:=&ast.BlockStmt{}
	block.Token=p.currToken
	block.Statements=[]ast.Statement{}

	p.next()

	for !p.isCurr(token.RBRACE) && !p.isCurr(token.EOF) {
		stmt:=p.parseStmt()
		if stmt!=nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.next()
	}
	return block 
}

// fn(x,y) {x+y;}
func (p* Parser) parseFunction() ast.Expression {
	expr:=&ast.Function{}
	expr.Token=p.currToken

	if !p.expected(token.LPAREN) {
		return nil 
	}

	expr.Params=p.parseFunctionParams()
	

	if !p.expected(token.RPAREN) {
		return nil 
	}

	if !p.expected(token.LBRACE) {
		return nil 
	}
	expr.Body=p.parseBlock()

	return expr 

}

// a,b,c,d)
func (p* Parser) parseFunctionParams() []*ast.Identifier {
	ids:=[]*ast.Identifier{}

	if p.isNext(token.RPAREN) {
		p.next()
		return ids 
	}

	p.next()
	id1:=&ast.Identifier{
		Token: p.currToken,
		Value: p.currToken.Value,
	}

	ids = append(ids, id1)

	for !p.isNext(token.RPAREN) {
		p.next()
		p.next() //to skip comma 

		ident:=&ast.Identifier{
			Token: p.currToken,
			Value: p.currToken.Value,
		}
		ids = append(ids, ident)
	}

	
	return ids
}


func (p* Parser) parseCallExpression(function ast.Expression) ast.Expression {
	call:=&ast.CallExpr{
		Token: p.currToken,
		Function: function,
	}
	call.Arguments=p.parseCallArgs()
	return call 
}

// add(2,3) currToken at (
func(p* Parser) parseCallArgs() []ast.Expression{
	var args []ast.Expression
	if p.isNext(token.RPAREN) {
		p.next()
		return args 
	}
	p.next()
	args = append(args, p.parseExpression(LOWEST))
	for !p.isNext(token.RPAREN) {
		p.next()
		p.next()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expected(token.RPAREN) {
		return nil 
	}

	return args 
	

}