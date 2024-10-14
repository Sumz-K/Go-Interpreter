package ast


import(
    "github.com/Sumz-K/Go-Interpreter/token"
)
type Node interface {
    TokenValue() string //every node in our AST has to return the token value corresponding to it
}

type Statement interface {
    Node  //all statement entities must also be nodes
    StatementNode() //dummy marker method
}

type Expression interface {
    Node 
    ExpressionNode()
}


type Program struct {
    Statements []Statement
}


func (p* Program) TokenValue() string {
    if len(p.Statements) > 0 {
        return p.Statements[0].TokenValue()
    } else {
        return ""
    }
}

/*

To represent a let statement like 
let x=5

we need 3 things, 
the token it itself
the variable name i.e x
the expression on the rhs i.e. 5
*/



type LetStmt struct {
    Token token.Token // the LET token 
    Name *Identifier  
    Value Expression
}


type Identifier struct {
    Token token.Token //the IDENTIFIER token
    Value string 
}


func (ls *LetStmt) TokenValue() string {
    return ls.Token.Value
}

func (ls *LetStmt) StatementNode() {}

func (id *Identifier) TokenValue() string {
    return id.Token.Value
}


func (id *Identifier) ExpressionNode() {}

type ReturnStmt struct {
    Token token.Token //the return token
    ReturnValue Expression
}


func (rs *ReturnStmt) StatementNode() {}

func (rs *ReturnStmt) TokenValue() string {
    return rs.Token.Value
}

