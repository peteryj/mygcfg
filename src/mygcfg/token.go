//
// This is the package for Token
//

package mygcfg

import (
    "fmt"
)

const (
	TOKEN_UNKNOWN = iota
	TOKEN_SECTION
	TOKEN_FIELD
	TOKEN_VALUE
	TOKEN_OP_LBRCE
	TOKEN_OP_RBRCE
	TOKEN_OP_EQUAL
)

type Token struct {
	tokenType int
	data      []byte
}

//
// New Token if needed
//
func NewToken(tokenType int, data []byte) (r *Token) {
	return &Token{tokenType, data}
}

func (tk *Token) TokenType() int {
	return tk.tokenType
}

func (tk Token) String() string {
    return fmt.Sprintf("(%d)%s", tk.tokenType, string(tk.data));
}

