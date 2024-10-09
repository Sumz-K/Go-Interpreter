package token 


const (
    ASSIGN = "="
    PLUS = "+"
    MINUS = "-"
    ASTERISK = "*"
    SLASH = "/"
    BANG = "!"
    LT = "<"
    GT = ">"

    EQ = "=="
    NOTEQ = "!=`"
    IDENTIFIER="IDENT"
    INTEGER="INT"

    COMMA=","
    SEMICOLON=";"

    LPAREN="("
    RPAREN=")"
    LBRACE="{"
    RBRACE="}"

    EOF="EOF"
    ILLEGAL="ILLEGAL"


    //keywords

    LET="LET"
    FUNC="FUNCTION"
    IF="IF"
    ELSE="ELSE"
    TRUE="TRUE"
    FALSE="FALSE"
    RETURN="RETURN"


)
type TokenType string

type Token struct {
    Type TokenType
    Value string
}


var keywords = map[string]TokenType {
    "let":LET,
    "fn":FUNC,
    "return":RETURN,
    "if":IF,
    "else":ELSE,
    "true":TRUE,
    "false":FALSE,
}

func CheckID(id string) TokenType { //checks if the id is a keyword or not
    tokType,ok:=keywords[id]
    if ok {
        return tokType
    }
    return IDENTIFIER
}
