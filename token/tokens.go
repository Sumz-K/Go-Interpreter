package token 


const (
    ASSIGN = "="
    ADD = "+"

    IDENTIFIER="IDENT"
    INTEGER="INT"

    COMMA=","
    SEMICOLON=";"

    LPARA="("
    RPARA=")"
    LBRACE="{"
    RBRACE="}"

    EOF="EOF"
    ILLEGAL="ILLEGAL"


    //keywords

    LET="LET"
    FUNC="FUNCTION"


)
type TokenType string

type Token struct {
    Type TokenType
    Value string
}
