package main

import (
	"fmt"
	"github.com/dean2021/sqlparser/parser"
	_ "github.com/dean2021/sqlparser/parser/test_driver"
)

func isOK(sql string) bool {
	p := parser.New()
	_, err := p.ParseOneStmt(sql, "", "")
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func main() {
	//lines := Load("/Users/user/Desktop/projects/sqlparser/tests")
	////i := 0
	//
	//for _, line := range lines {
	//	if !isOK(line) {
	//		code := test3.Fix(line)
	//		if !isOK(code) {
	//			fmt.Println("修复前:", line)
	//			fmt.Println("修复后:", code)
	//		}
	//	}
	//}
	//
	//p := parser.New()
	//_, err := p.ParseOneStmt(`=1 or 1=1`, "", "")
	//if err != nil {
	//	fmt.Println(err)
	//}

	fmt.Println(isOK("1'='2'"))
	//fmt.Println(isOK(test3.Fix(`1) or sleep(__TIME__)#`)))
}
