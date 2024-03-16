package test3

import (
	"regexp"
	"strings"
)

type Token struct {
	Type  string
	Token string
}

type StateFunc func(char byte) StateFunc

func start(char byte) StateFunc {
	if regexp.MustCompile(`\d`).MatchString(string(char)) {
		currentToken = &Token{Type: "Number", Token: string(char)}
		return inNumber
	} else if regexp.MustCompile(`\s`).MatchString(string(char)) {
		emitToken(Token{Type: "Space", Token: " "})
		return start
	} else if regexp.MustCompile(`"`).MatchString(string(char)) {
		currentToken = &Token{Type: "DoubleString", Token: string(char)}
		return inDoubleString
	} else if char == byte('\'') {
		currentToken = &Token{Type: "SingleString", Token: string(char)}
		return inSingleString
	} else if regexp.MustCompile(`[=+\-*\%]`).MatchString(string(char)) {
		emitToken(Token{Type: "Operator", Token: string(char)})
		return start
	} else if char == byte('_') {
		currentToken = &Token{Type: "Operator", Token: string(char)}
		return maybeIdentifier
	} else if char == byte('/') {
		currentToken = &Token{Type: "Operator", Token: string(char)}
		return maybeComment
	} else if regexp.MustCompile(`[a-zA-Z]`).MatchString(string(char)) {
		currentToken = &Token{Type: "Identifier", Token: string(char)}
		return inIdentifier
		//} else if char == byte(')') {
		//	//currentToken = &Token{Type: "xxxx", Token: string(char)}
		//	//currentToken.Token += string(char)
		//	//emitToken(*currentToken)
		//	//currentToken = nil
		//	return start
	} else {
		//currentToken = &Token{Type: "Bareword", Token: string(char)}
		//return inBareword
		return start
	}
	return nil
}

func maybeIdentifier(char byte) StateFunc {
	if char == '_' {
		currentToken.Type = "Operator"
		currentToken.Token += string(char)
		return maybeIdentifier
	} else if regexp.MustCompile(`[a-zA-Z]`).MatchString(string(char)) {
		currentToken.Type = "Identifier"
		currentToken.Token += string(char)
		return inIdentifier
	} else if regexp.MustCompile(`[0-9]`).MatchString(string(char)) {
		// 如果_后面不是数字且不是操作符
		currentToken.Type = "Bareword"
		currentToken.Token += string(char)
		return inBareword
	} else if !regexp.MustCompile(`[()=+\-*\%\s]`).MatchString(string(char)) {
		// 如果_后面不是数字且不是操作符
		currentToken.Type = "Bareword"
		currentToken.Token += string(char)
		return inBareword
	} else {
		emitToken(*currentToken)
		currentToken = nil
		return start(char)
	}
}

func maybeComment(char byte) StateFunc {
	if char == '*' {
		currentToken.Type = "Comment"
		currentToken.Token += string(char)
		return inComment
	} else {
		emitToken(*currentToken)
		currentToken = nil
		return start(char)
	}
}

func inComment(char byte) StateFunc {
	currentToken.Token += string(char)
	if char == '*' {
		return maybeEndComment
	}
	return inComment
}

func maybeEndComment(char byte) StateFunc {
	currentToken.Token += string(char)
	if char == '/' {
		emitToken(*currentToken)
		currentToken = nil
		return start
	} else {
		return inComment
	}
}

func maybeFunctionCallEnd(char byte) StateFunc {
	if regexp.MustCompile(`\s`).MatchString(string(char)) {
		emitToken(*currentToken)
		currentToken = nil
		return start(char)
	} else {
		currentToken.Type = "Bareword"
		currentToken.Token += string(char)
		return inBareword
	}
}

func inNumber(char byte) StateFunc {
	if regexp.MustCompile(`\d`).MatchString(string(char)) {
		currentToken.Token += string(char)
		return inNumber
		//} else if !regexp.MustCompile(`[=+\-*\%\s]`).MatchString(string(char)) {
		//	currentToken.Type = "Bareword"
		//	currentToken.Token += string(char)
		//	return inBareword
	} else {
		emitToken(*currentToken)
		currentToken = nil
		return start(char)
	}
}

func inString(char byte) StateFunc {
	if char == '\'' {
		currentToken.Token += string(char)
		emitToken(*currentToken)
		startRepairString = "'"
		currentToken = nil
		return start
	} else if char == '"' {
		currentToken.Token += string(char)
		emitToken(*currentToken)
		startRepairString = "\""
		currentToken = nil
		return start
	} else {
		currentToken.Token += string(char)
		return inString
	}
}

func inDoubleString(char byte) StateFunc {
	if char == '"' {
		currentToken.Token += string(char)
		emitToken(*currentToken)
		currentToken = nil
		return start
	} else {
		currentToken.Token += string(char)
		return inDoubleString
	}
}

func inSingleString(char byte) StateFunc {
	if char == '\'' {
		currentToken.Token += string(char)
		emitToken(*currentToken)
		currentToken = nil
		return start
	} else {
		currentToken.Token += string(char)
		return inSingleString
	}
}

func inBareword(char byte) StateFunc {
	// 如果没遇到空格一直是bareword
	if !regexp.MustCompile(`\s`).MatchString(string(char)) {
		currentToken.Token += string(char)
		return inBareword
	} else {
		emitToken(*currentToken)
		currentToken = nil
		return start(char)
	}
}

func inIdentifier(char byte) StateFunc {
	if regexp.MustCompile(`[a-zA-Z0-9_]`).MatchString(string(char)) {
		currentToken.Token += string(char)
		return inIdentifier
	} else if char == '(' {
		currentToken.Type = "Function"
		currentToken.Token += string(char)
		return inFunction
	} else if char == '(' {
		currentToken.Type = "Function"
		currentToken.Token += string(char)
		return inFunction
	} else if !regexp.MustCompile(`[=+\-*\%\s]`).MatchString(string(char)) {
		currentToken.Type = "Bareword"
		currentToken.Token += string(char)
		return inBareword
	} else {
		emitToken(*currentToken)
		currentToken = nil
		return start(char)
	}
}

func inFunction(char byte) StateFunc {
	if char == ')' {
		currentToken.Token += string(char)
		emitToken(*currentToken)
		currentToken = nil
		return start
	}
	currentToken.Token += string(char)
	return inFunction
}

func emitToken(token Token) {
	tokens = append(tokens, token)
}

func tokenizer(input string) []Token {
	state := start
	if strings.Contains(input, "'") || strings.Contains(input, "\"") {
		if !strings.HasPrefix(input, "'") && !strings.HasPrefix(input, `""`) {
			currentToken = &Token{Type: "String", Token: ``}
			state = inString
		} else if strings.HasPrefix(input, "'") {
			currentToken = &Token{Type: "SingleString", Token: "'"}
			state = inSingleString
		} else if strings.HasPrefix(input, `"`) {
			currentToken = &Token{Type: "DoubleString", Token: `"`}
			state = inDoubleString
		}
	}

	for i := 0; i < len(input); i++ {
		state = state(input[i])
	}

	if currentToken != (nil) {
		emitToken(*currentToken)
	}
	return tokens
}

var tokens []Token
var currentToken *Token
var startRepairString = ""
var endRepairString = ""

func Fix(code string) string {

	result := ""
	startRepairString = ""
	endRepairString = ""
	tokens = []Token{}
	currentToken = nil
	codeTokens := tokenizer(code)

	for i, token := range codeTokens {
		//	fmt.Println(token)
		//result = result + token.Token
		//fmt.Println(token)
		//if token.Type == "Bareword" {
		//	continue
		//}
		result = result + token.Token
		if token.Type == "SingleString" && i == len(codeTokens)-1 && !strings.HasSuffix(token.Token, "'") {
			endRepairString = "'"
		}

		if token.Type == "DoubleString" && i == len(codeTokens)-1 && !strings.HasSuffix(token.Token, `"`) {
			endRepairString = `"`
		}

		// ======
		if token.Type == "Operator" && i == len(codeTokens)-1 && !strings.HasSuffix(token.Token, `"`) {
			endRepairString = `T`
		}

		if token.Type == "Space" && i == 0 {
			startRepairString = "T"
		}
	}

	//fmt.Println("before", code)
	//fmt.Println("repair", startRepairString+result+endRepairString)
	//fmt.Println("####")
	s := startRepairString + result + endRepairString
	//fmt.Println(s)
	return s
}

//code := "1231312 31' or '1'='1"
//code := "version/**/() )/**/1=1 or 1 or name=\"(xx x)\" and  1=sleep(1)  or age=(1 or 1=1) or 1=1/* 1 or 1=1"

//code := ") 213"
//codeTokens := tokenizer(code)
//result := ""
//for i, token := range codeTokens {
//	//result = result + token.Token
//	//fmt.Println(token)
//	if token.Type == "Bareword" {
//		continue
//	}
//	result = result + token.Token
//	if token.Type == "String" && i == len(codeTokens)-1 && !strings.HasSuffix(token.Token, "'") {
//		result += "'"
//	}
//}
//
//fmt.Println(result)
