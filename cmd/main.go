package main

import (
	"fmt"
	"github.com/dean2021/sqlparser/parser"
	"github.com/dean2021/sqlparser/parser/ast"
	"github.com/dean2021/sqlparser/parser/opcode"
	"github.com/dean2021/sqlparser/parser/test_driver"
	_ "github.com/dean2021/sqlparser/parser/test_driver"
	"strings"
)

type SQLiDetect struct {
	isRisk bool
	count  int
}

func (v *SQLiDetect) Enter(in ast.Node) (ast.Node, bool) {
	if _, ok := in.(*ast.CreateTableStmt); ok {
		v.isRisk = true
	}
	if _, ok := in.(*ast.WaitForStmt); ok {
		v.isRisk = true
	}
	if ce, ok := in.(*ast.CommonExpressionStmt); ok {
		if op, ok := ce.Expr.(*ast.FuncCallExpr); ok {
			strArray := []string{"sleep", "benchmark", "pg_sleep", "exec", "randomblob"}
			for _, value := range strArray {
				if value == strings.ToLower(op.FnName.String()) {
					v.isRisk = true
					break
				}
			}
		}
		if op, ok := ce.Expr.(*ast.BinaryOperationExpr); ok {
			if op.L != nil {
				if op, ok := op.L.(*ast.BinaryOperationExpr); ok {
					if op.Op == opcode.LogicOr || op.Op == opcode.LogicAnd {
						v.isRisk = true
					}
				}
			}
			if op.Op == opcode.LogicOr || op.Op == opcode.LogicAnd {
				v.isRisk = true
			}

		}

		if _, ok := ce.Expr.(*ast.SubqueryExpr); ok {
			v.isRisk = true
		}

		// TODO 可能存在误报
		if _, ok := ce.Expr.(*ast.ParenthesesExpr); ok {
			v.isRisk = true
		}
		if val, ok := ce.Expr.(*test_driver.ValueExpr); ok {
			if str, ok := val.GetValue().(string); ok {
				if check(str) {
					v.isRisk = true
				}
			}
		}

		//if v.isRisk == false {
		//	fmt.Println(reflect.TypeOf(ce.Expr))
		//}
	}

	if _, ok := in.(*ast.SelectStmt); ok {
		v.isRisk = true
	}

	if val, ok := in.(*test_driver.ValueExpr); ok {
		if str, ok := val.GetValue().(string); ok {
			if check(str) {
				v.isRisk = true
			}
		}
	}

	return in, false
}

func (v *SQLiDetect) Leave(in ast.Node) (ast.Node, bool) {

	//fmt.Println(reflect.TypeOf(in))
	return in, true
}

func parse(sql string) (*ast.StmtNode, error) {
	p := parser.New()
	stmtNodes, err := p.ParseOneStmt(sql, "", "")
	if err != nil {
		return nil, err
	}
	return &stmtNodes, nil
}

func check(sql string) bool {
	// 太短就不检查了
	if len(sql) < 3 {
		return false
	}
	astNode, err := parse(sql)
	if err != nil {
		fmt.Println("语法错误:", sql)
		return false
	}
	v := &SQLiDetect{}
	(*astNode).Accept(v)
	return v.isRisk
}

func main() {
	////
	//startTime := time.Now()
	lines := Load("/Users/user/Desktop/projects/sqlparser/tests")
	//i := 0
	for _, line := range lines {
		//astNode, err := parse(line)
		//if err != nil {
		//	//fmt.Printf("parse error: %v\n", err.Error())
		//	//fmt.Println("语法错误:", line)
		//	i++
		//	continue
		//}
		//v := &SQLiDetect{}
		//(*astNode).Accept(v)
		if check(line) {
			//fmt.Println("发现sql注入")
		} else {
			//fmt.Println("漏报:", line)
		}
	}

	//
	//endTime := time.Now()
	//elapsedTime := endTime.Sub(startTime)
	//fmt.Printf("代码执行耗时: %s \n", elapsedTime)

	//语法错误: sleep(__TIME__)#
	//语法错误: ;waitfor delay '0:0:__TIME__'--
	//语法错误: );waitfor delay '0:0:__TIME__'--
	//语法错误: ';waitfor delay '0:0:__TIME__'--
	//语法错误: ";waitfor delay '0:0:__TIME__'--
	//语法错误: ');waitfor delay '0:0:__TIME__'--
	//语法错误: ");waitfor delay '0:0:__TIME__'--
	//语法错误: ));waitfor delay '0:0:__TIME__'--
	//语法错误: '));waitfor delay '0:0:__TIME__'--
	//语法错误: "));waitfor delay '0:0:__TIME__'--
	//语法错误: benchmark(10000000,MD5(1))#
	//语法错误: 1 or benchmark(10000000,MD5(1))#
	//语法错误: " or benchmark(10000000,MD5(1))#
	//语法错误: ' or benchmark(10000000,MD5(1))#
	//语法错误: 1) or benchmark(10000000,MD5(1))#
	//语法错误: ") or benchmark(10000000,MD5(1))#
	//语法错误: ') or benchmark(10000000,MD5(1))#
	//语法错误: 1)) or benchmark(10000000,MD5(1))#
	//语法错误: ")) or benchmark(10000000,MD5(1))#
	//语法错误: ')) or benchmark(10000000,MD5(1))#
	//语法错误:  OR 3409=3409 AND ('pytW' LIKE 'pytW
	//语法错误:  OR 3409=3409 AND ('pytW' LIKE 'pytY

	// 优先按照函数解析，解析失败再按照字面量
	//line := `1=1 or 1=1 xx 1=1 or 1=1` // or 1 or name="(xx x)" and  1=sleep(1)  or age=(1 or 1=1) or 1=1/* 1 or 1=1`
	//// version( )1=1 or 1 or name="(xx x)" and  1=sleep(1)  or age=(1 or 1=1) or 1=1/* 1 or 1=1
	//line := `WAITFOR DELAY 'xx'`
	//astNode, err := parse(line)
	//if err != nil {
	//	fmt.Printf("parse error: %v\n", err.Error())
	//	return
	//} else {
	//	fmt.Println("语法正确")
	//}
	//
	//v := &SQLiDetect{}
	//(*astNode).Accept(v)
	//if v.isRisk {
	//	fmt.Println("发现sql注入")
	//} else {
	//	fmt.Println("漏报:", line)
	//}
}
