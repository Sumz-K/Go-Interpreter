package ast


import(
    "github.com/Sumz-K/Go-Interpreter/token"
    "bytes"
)
type Node interface {
    TokenValue() string //every node in our AST has to return the token value corresponding to it
    String() string //to print out nodes
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

func (p *Program) String() string {
    var buffer bytes.Buffer
    for _,stmt:=range p.Statements {
        buffer.WriteString(stmt.String())
    }
    return buffer.String()
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

func (ls *LetStmt) String() string {
    var buf bytes.Buffer
    buf.WriteString(ls.TokenValue()+ " ")
    buf.WriteString(ls.Name.String())
    buf.WriteString(" = ")
    if ls.Value!=nil {
        buf.WriteString(ls.Value.String())
    }
    buf.WriteString(";")
    return buf.String()
}

func (id *Identifier) TokenValue() string {
    return id.Token.Value
}


func (id *Identifier) ExpressionNode() {}

func(id *Identifier) String() string {
    return id.Value
}

type ReturnStmt struct {
    Token token.Token //the return token
    ReturnValue Expression
}


func (rs *ReturnStmt) StatementNode() {}

func (rs *ReturnStmt) TokenValue() string {
    return rs.Token.Value
}


func(rs *ReturnStmt) String() string {
    var buf bytes.Buffer
    buf.WriteString(rs.TokenValue()+" ")

    if rs.ReturnValue!=nil {
        buf.WriteString(rs.ReturnValue.String())
    }
    buf.WriteString(";")
    return buf.String()
}


type ExpressionStmt struct { // wrapper to denote stmts like x+10 which are valid standalone statements
    Token token.Token
    Expression Expression
}

func (es *ExpressionStmt) TokenValue() string {
    return es.Token.Value
}

func (es *ExpressionStmt) StatementNode() {}

func (es *ExpressionStmt) String() string {
    if es.Expression!=nil {
        return es.Expression.String()
    }
    return ""
}



type IntegerLiteral struct {
    Token token.Token
    Value int64
}

func (il *IntegerLiteral) ExpressionNode() {}

func (il *IntegerLiteral) TokenValue() string {
    return il.Token.Value
}

func (il *IntegerLiteral) String() string {
    return il.Token.Value
}


// two types 
// -5 and !5
type PrefixExpression struct {
    Token token.Token
    Operator string 
    Right Expression // the rest of the expression
}

func (pe* PrefixExpression) ExpressionNode() {}
func (pe *PrefixExpression) TokenValue() string {
    return pe.Token.Value
}

func (pe *PrefixExpression) String() string {
    var buf bytes.Buffer

    buf.WriteString("(")
    buf.WriteString(pe.Operator)
    buf.WriteString(pe.Right.String())
    buf.WriteString(")")

    return buf.String()
}

/*
to represent:
5+5
5-5
5*5
5/5
5>5
5<5
5==5
5!=5
*/
type InfixExpression struct {
    Token token.Token // the operator token, like +
    LeftExpr Expression
    Operator string
    RightExpr Expression
}

func (ie *InfixExpression) TokenValue() string {
    return ie.Token.Value
}

func (ie *InfixExpression) ExpressionNode() {}

func (ie* InfixExpression) String() string {
    var buf bytes.Buffer
    buf.WriteString("(")
    buf.WriteString(ie.LeftExpr.String())
    buf.WriteString(" "+ ie.Operator+ " ")
    buf.WriteString(ie.RightExpr.String())
    buf.WriteString(")")

    return buf.String()
}

type Boolean struct {
    Token token.Token
    Value bool 
}

func (b* Boolean) TokenValue() string {
    return b.Token.Value
}

func (b *Boolean) ExpressionNode() {}

func(b* Boolean) String() string {
    return b.TokenValue()
}


type BlockStmt struct {
    Token token.Token
    Statements []Statement
}

func (bs *BlockStmt) StatementNode() {}

func (bs *BlockStmt) TokenValue() string {
    return bs.Token.Value
}

func (bs *BlockStmt) String() string{
    var buf bytes.Buffer
    buf.WriteString("{")
    for _,ele:=range bs.Statements {
        buf.WriteString(ele.String())
    }
    buf.WriteString("}")
    return buf.String()
}

type IfExpression struct {
    Token token.Token
    Condition Expression
    Consequence *BlockStmt
    Alternative *BlockStmt
}

func (ifexp *IfExpression) ExpressionNode() {}

func (ifexp *IfExpression) TokenValue() string {
    return ifexp.Token.Value
}

func (ifexp *IfExpression) String() string {
    var buf bytes.Buffer
    buf.WriteString("if ")
    buf.WriteString(ifexp.Condition.String())
    buf.WriteString(" ")
    buf.WriteString(ifexp.Condition.String())

    if ifexp.Alternative!=nil {
        buf.WriteString(" else ")
        buf.WriteString(ifexp.Alternative.String())
    }
    return buf.String()
}


type Function struct {
    Token token.Token
    Params []*Identifier
    Body *BlockStmt
}

func (fn *Function) TokenValue() string {
    return fn.Token.Value
}

func (fn *Function) ExpressionNode() {}

func (fn *Function) String() string {
    var buf bytes.Buffer
    buf.WriteString("fn")
    buf.WriteString("(")
    for _,param:=range fn.Params {
        buf.WriteString(param.String()+", ")
    }
    buf.WriteString(") {")
    buf.WriteString(fn.Body.String())
    buf.WriteString("}")
    return buf.String()

}

// for function calls
type CallExpr struct {
    Token token.Token // The '(' token
    Function Expression // a call can be of the form add(2,3) but it can also be of the form fn(x,y){x+y;} (2,3). Now if the function part is made an expression instead of identifier(which also is an expression) we can validate both these calls. Both "add" and "fn(x,y) {x+y;}" ar eessentially still 'expressions'. One is an Identifier struct other is a Function struct but both fulfill the Expression interface
    Arguments []Expression  // no * because Expression is an interface and interfaces already hold references to concrete types
}

func (ce *CallExpr) ExpressionNode() {}

func (ce *CallExpr) TokenValue() string {
    return ce.Token.Value
}

func (ce *CallExpr) String() string {
    var buf bytes.Buffer
    buf.WriteString(ce.Function.String())
    buf.WriteString("(")
    for _,arg:=range ce.Arguments {
        buf.WriteString(arg.String()+", ")
    }
    buf.WriteString(")")

    return buf.String()
}


