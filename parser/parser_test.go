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