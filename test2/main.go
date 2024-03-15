package main

import (
	"fmt"
	"regexp"
	"strings"
)

type Token struct {
	Type  string
	Token string
}

var currentToken *Token
var tokens []Token

func main() {
	code := ")version( )/**/1=1 or 1 or name=(xx x) and 1=sleep(1) or age=(1 or 1=1) or 1=1/* 1 or 1=1"
	tokens = tokenizer(code)
	fmt.Println(tokens)
}

func tokenizer(input string) []Token {
	state := start
	for _, r := range input {
		char := string(r)
		state = state(char)
	}
	if currentToken != nil {
		emitToken(*currentToken)
	}
	return tokens
}

func start(char string) func(string) func(string) {
	if isDigit(char) {
		currentToken = &Token{Type: "Number", Token: char}
		return inNumber
	} else if strings.TrimSpace(char) == "" {
		return start
	} else if isOperator(char) {
		emitToken(Token{Type: "Operator", Token: char})
		return start
	} else if char == "/" {
		currentToken = &Token{Type: "Operator", Token: char}
		return maybeComment
	} else if isAlpha(char) {
		currentToken = &Token{Type: "Identifier", Token: char}
		return inIdentifier
	}
	return nil
}

func maybeComment(char string) func(string) func(string) {
	if char == "*" {
		currentToken.Type = "Comment"
		currentToken.Token += char
		return inComment
	} else {
		emitToken(*currentToken)
		currentToken = nil
		return start(char)
	}
}

func inComment(char string) func(string) func(string) {
	currentToken.Token += char
	if char == "*" {
		return maybeEndComment
	}
	return inComment
}

func maybeEndComment(char string) func(string) func(string) {
	currentToken.Token += char
	if char == "/" {
		emitToken(*currentToken)
		currentToken = nil
		return start
	} else {
		return inComment
	}
}

func inNumber(char string) func(string) func(string) {
	if isDigit(char) {
		currentToken.Token += char
		return inNumber
	} else {
		emitToken(*currentToken)
		currentToken = nil
		return start(char)
	}
}

func inIdentifier(char string) func(string) func(string) {
	if isAlpha(char) {
		currentToken.Token += char
		return inIdentifier
	} else if char == "(" {
		currentToken.Type = "Function"
		currentToken.Token += char
		return inFunction
	} else {
		emitToken(*currentToken)
		currentToken = nil
		return start(char)
	}
}

func inFunction(char string) func(string) func(string) {
	if char != ")" && currentToken != nil {
		currentToken.Token += char
	}
	if char == ")" && currentToken != nil {
		emitToken(*currentToken)
		currentToken = nil
	}
	return inFunction
}

func emitToken(token Token) {
	tokens = append(tokens, token)
}

func isDigit(char string) bool {
	return regexp.MustCompile(`\d`).MatchString(char)
}

func isOperator(char string) bool {
	return regexp.MustCompile(`["'()=+-]`).MatchString(char)
}

func isAlpha(char string) bool {
	return regexp.MustCompile(`[a-z]`).MatchString(strings.ToLower(char))
}
