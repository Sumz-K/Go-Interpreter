package parser

import (
	"testing"

	"github.com/Sumz-K/Go-Interpreter/ast"
	"github.com/Sumz-K/Go-Interpreter/lexer"
)

func TestLetStatements(t *testing.T) {
    input:=`
    let x=5;
    let y=10;
    let foobar=838383;
    `

    l:=lexer.New(input)
    p:=New(l)

    program:=p.ParseProgram()
    checkErrors(t,p)

    if program==nil {
        t.Fatalf("ParseProgram() returned nil")
    }

    if len(program.Statements)!=3 {
        t.Fatalf("Needed 3 elements after parsing got %d",len(program.Statements))
    }

    

    tests:=[]struct{
        expectedIdentifier string
    }{
        {"x"},
        {"y"},
        {"foobar"},
    }

    for i,tt:=range tests {
        stmt:=program.Statements[i]
        if !testLetStmtHelper(t,stmt,tt.expectedIdentifier) {
            return 
        }
    }
}

func checkErrors(t *testing.T, p * Parser) {
    errors:=p.showErrors()
    if len(errors) ==0 {
        return
    }

    t.Errorf("Parser has %d errors",len(errors))

    for _,err:= range errors {
        t.Errorf("Parser error: %q",err)
    }   
    t.FailNow()
}

func testLetStmtHelper(t *testing.T, stmt ast.Statement,id string) bool{
    if stmt.TokenValue() !="let" {
        t.Errorf("Token type: expected let got %v",stmt.TokenValue())
        return false
    }

    letStmt,ok:=stmt.(*ast.LetStmt)

    if !ok {
        t.Errorf("Not *ast.LetStmt got %v",stmt)
        return false 
    }

    if letStmt.Name.Value!=id {
        t.Errorf("Variable name Expected %v got %v",id,letStmt.Name.Value)
        return false 
    }

    if letStmt.Name.TokenValue() !=id {
        t.Errorf("Expected %v got %v",id,letStmt.Name)
    }

    return true 


}



func TestReturnStatements(t *testing.T) {
        input := `
    return 5;
    return 10;
    return 993322;
    `

    l:=lexer.New(input)

    p:=New(l)

    program:=p.ParseProgram()

    checkErrors(t,p)
    
    if program==nil {
        t.Fatalf("Parse Program() returned nil for return statements")
    }

    if len(program.Statements)!=3 {
        t.Fatalf("Needed 3 elements after parsing returns got %d",len(program.Statements))
    }

    for _,stmt:= range program.Statements {
        returnStmt,ok:=stmt.(*ast.ReturnStmt)
        
        if !ok {
            t.Errorf("Statement not of type return, got %T",stmt)
            continue
        }

        if returnStmt.TokenValue()!="return" {
            t.Errorf("Expected token name return got %q",returnStmt.TokenValue())
        }
    }
    
}


func TestIdentifierExpression(t *testing.T) { 
        input := "foobar;"
        l := lexer.New(input)
        p := New(l)
        program := p.ParseProgram()
        checkErrors(t, p)
        if len(program.Statements) != 1 {
            t.Fatalf("program has not enough statements. got=%d",len(program.Statements))
        }
        stmt, ok := program.Statements[0].(*ast.ExpressionStmt) 
        if !ok {
                t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
                    program.Statements[0])
        }
        ident, ok := stmt.Expression.(*ast.Identifier) 
        if !ok {
            t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
        }
        if ident.Value != "foobar" {
            t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
        }
        if ident.TokenValue() != "foobar" {
                t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar", ident.TokenValue())
        } 
}


func TestIntLiteral(t *testing.T) {
    input:="5;"
    l:=lexer.New(input)
    p:=New(l)

    program:=p.ParseProgram()
    checkErrors(t,p)
    if len(program.Statements)!=1{
        t.Fatalf("Program does not have the right number of elements, expected 1 got %q",len(program.Statements))
    }

    stmt,ok:=program.Statements[0].(*ast.ExpressionStmt)
    if !ok {
        t.Fatalf("The statement is not an expression statement, got %T",program.Statements[0])
    }

    intLit,ok:=stmt.Expression.(*ast.IntegerLiteral)
    if !ok  {
        t.Fatalf("The statement is not an integer literal, got %T",stmt)
    }

    if intLit.TokenValue()!="5" {
        t.Errorf("intLit.TokenValue() wrong, expedcted %q got %q","5",intLit.TokenValue())
    }

}

func TestPrefixParse(t *testing.T) {
    tests:=[]struct{
        input string 
        operator string 
        value int64 
    }{
        {"!5","!",5},
        {"-5","-",5},
    }

    for _,tt:= range tests {
        l:=lexer.New(tt.input)
        p:=New(l)
        program:=p.ParseProgram()
        checkErrors(t,p)

        if len(program.Statements)!=1 {
            t.Fatalf("Expected 1 statement got %d",len(program.Statements))
        }

        stmt,ok:=program.Statements[0].(*ast.ExpressionStmt)
        if !ok  {
            t.Fatalf("Statement not an expression statement, got %T",program.Statements[0])
        }

        prefixStmt,ok:=stmt.Expression.(*ast.PrefixExpression)
        if !ok {
            t.Fatalf("Statement not a prefix expression, got %T",stmt.Expression)
        }

        if prefixStmt.Operator!=tt.operator{
            t.Errorf("Expected operator to be %s but got %s",tt.operator,prefixStmt.Operator)
        }

        if !compareInt(t,prefixStmt.Right,tt.value) {
            return 
        }
    }
}

func compareInt(t *testing.T,il ast.Expression,val int64) bool{
    intVal,ok:=il.(*ast.IntegerLiteral)
    if !ok {
        t.Errorf("Expected IntegerLiteral got %T",il)
        return false 
    }

    if intVal.Value!=val {
        t.Errorf("integer value wrong, expected %d got %d",val,intVal.Value)
        return false 
    }
    return true
}


func TestParsingInfixExpressions(t *testing.T) { 
    infixTests := []struct {
        input      string
        leftValue  int64
        operator   string
        rightValue int64
    }{
        {"5 + 5;", 5, "+", 5},
        {"5 - 5;", 5, "-", 5},
        {"5 * 5;", 5, "*", 5},
        {"5 / 5;", 5, "/", 5},
        {"5 > 5;", 5, ">", 5},
        {"5 < 5;", 5, "<", 5},
        {"5 == 5;", 5, "==", 5},
        {"5 != 5;", 5, "!=", 5},
    }

    for _,tt:= range infixTests {
        l:=lexer.New(tt.input)
        p:=New(l)
        program:=p.ParseProgram()
        checkErrors(t,p)

        if len(program.Statements)!=1 {
            t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
        }

        stmt,ok:=program.Statements[0].(*ast.ExpressionStmt)
        if !ok {
            t.Fatalf("The statement is not an expression statement got %T",program.Statements[0])
        }

        expr,ok:=stmt.Expression.(*ast.InfixExpression)
        if !ok {
            t.Fatalf("The expression is not an infix expression, got %T",stmt.Expression)
        }

        if !compareInt(t,expr.LeftExpr,tt.leftValue) {
            return 
        }

        if tt.operator!=expr.Operator {
            t.Errorf("Incorrect operator,expected %v got %v",tt.operator,expr.Operator)
        }

        if !compareInt(t,expr.RightExpr,tt.rightValue) {
            return 
        }
    }

}