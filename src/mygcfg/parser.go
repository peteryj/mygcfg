//
// This file defines Parser 
//

package mygcfg

import (
	"bufio"
    "os"
    "fmt"
)

const (
	PARSE_SUCCESS = 0
	ERR_TOKEN     = iota
	ERR_PARSE
	ERR_LINE_LEN0
)

const (
	T_UNKNOWN = 0
	T_LBRACE  = iota
	T_RBRACE
	T_EQUAL
	T_STRING
	T_SECTION
	T_KEY
	T_VALUE
)

func errStr (err int) string {
    switch err {
    case ERR_TOKEN:
        return "token error"
    case ERR_PARSE:
        return "parse error"
    case ERR_LINE_LEN0:
        return "line length 0"
    default:
        return ""
    }
    return ""
}

func gotError (err int) error {
    return fmt.Errorf("error(%d):%s!", err, errStr(err))
}

type ConfigSec map[string]string

type Parser struct {
	tokenList []Token
	nToken    uint
	state     uint
	output    map[string]ConfigSec
}

////////////////////////////////////////////
// Level 1 functions

//
// This function parse the file, fill in the tokenList and output.
//
func (pr *Parser) ParseFile(fileName string) (err error) {
    fmt.Printf("handling %s...\n", fileName);

    f, err := os.OpenFile(fileName, os.O_RDONLY, 0666)
	if  err != nil {
		return err
	}

	bufReader := bufio.NewReader(f)
    pr.tokenList = make([]Token, 100, 100);

	for {
        line, err := bufReader.ReadBytes('\n')
		if err != nil {
			break
		}

		if err = pr.parseTokenOneLine(line); err != nil {
			pr.state = ERR_TOKEN
			return err
		}
	}

	return pr.genOutput()
}

//
// This function return the out ConfigSet
//
func (pr *Parser) Output() (r *map[string]ConfigSec, err error) {
	if pr.state != PARSE_SUCCESS {
		return nil, gotError(ERR_PARSE)
	}

	return &pr.output, nil
}

/////////////////////////////////////////////
// Level 2 functions

//
// This function parse tokens for oneline
//
func (pr *Parser) parseTokenOneLine(line []byte) (err error) {
	if len(line) == 0 {
		return gotError(ERR_LINE_LEN0)
	}

	cToken := &pr.nToken
	curToken := T_UNKNOWN
	handleLastToken := false
	curType := T_UNKNOWN
	lastPos := 0

	for curPos, c := range line {
		switch {
		case '[' == c:
			curToken = T_LBRACE
			curType = T_SECTION
			lastPos = curPos + 1
		case ']' == c:
			if curToken == T_STRING {
				curToken = T_RBRACE
				handleLastToken = true
			}
		case '=' == c:
			curToken = T_EQUAL
			curType = T_KEY
			handleLastToken = true
		case '\n' == c:
			if curToken == T_STRING {
				curType = T_VALUE
				handleLastToken = true
			}
		case '0' <= c && c <= '9':
			fallthrough
		case 'a' <= c && c <= 'z':
			fallthrough
		case 'A' <= c && c <= 'Z':
			curToken = T_STRING
		default:
			// just ignore all other undefined chars
			continue
		}

		if handleLastToken {
			pr.tokenList[*cToken] = Token{curType, line[lastPos:curPos]}
			lastPos = curPos + 1

			*cToken++
			handleLastToken = false
		}
	}

    return nil
}

//
// This function go through the tokenList, check the syntax
// and fill the out ConfigSet
//
func (pr *Parser) genOutput() (err error) {

    pr.output = make(map[string]ConfigSec)

    fmt.Printf("%v\n", pr.tokenList);

    var curSection, curKey, curValue string
	for _, tk := range pr.tokenList {
		switch tk.tokenType {
		case T_SECTION:
			curSection = string(tk.data)

            fmt.Printf("%v %v\n", tk, curSection)
			pr.output[curSection] = make(map[string]string)
		case T_KEY:
			curKey = string(tk.data)
		case T_VALUE:
			curValue = string(tk.data)

            fmt.Printf("%v %v\n", tk, curSection)
			pr.output[curSection][curKey] = curValue
		}
	}
    return nil
}
