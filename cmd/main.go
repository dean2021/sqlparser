package main

import (
	"fmt"
	"github.com/dean2021/sqlparser/parser"
	"github.com/dean2021/sqlparser/parser/ast"
	"github.com/dean2021/sqlparser/parser/opcode"
	"github.com/dean2021/sqlparser/parser/test_driver"
	_ "github.com/dean2021/sqlparser/parser/test_driver"
	"reflect"
	"strings"
)

type SQLiDetect struct {
	isRisk bool
	count  int
}

func (v *SQLiDetect) Enter(in ast.Node) (ast.Node, bool) {

	if ce, ok := in.(*ast.CommonExpressionStmt); ok {
		fmt.Println(reflect.TypeOf(ce.Expr))

		if op, ok := ce.Expr.(*ast.FuncCallExpr); ok {
			strArray := []string{"sleep"}
			for _, value := range strArray {
				if value == strings.ToLower(op.FnName.String()) {
					v.isRisk = true
					break
				}
			}
		}

		if op, ok := ce.Expr.(*ast.BinaryOperationExpr); ok {
			fmt.Println(op.L)
			if op.L != nil {
				if op, ok := op.L.(*ast.BinaryOperationExpr); ok {

					if op.Op == opcode.LogicOr || op.Op == opcode.LogicAnd {
						v.isRisk = true
					}
				}
			}
			//
			//if op.R != nil {
			//	if op, ok := op.R.(*ast.BinaryOperationExpr); ok {
			//		if op.Op == opcode.LogicOr || op.Op == opcode.LogicAnd {
			//			v.isRisk = true
			//		}
			//	}
			//}

			if val, ok := op.L.(*test_driver.ValueExpr); ok {
				fmt.Println("val:", val.GetType(), val.GetValue())
			}
			//fmt.Println(reflect.TypeOf(op.L))
			if op.Op == opcode.LogicOr || op.Op == opcode.LogicAnd {
				v.isRisk = true
			}
		}
	}

	if _, ok := in.(*ast.SelectStmt); ok {
		v.isRisk = true
		//// 子查询
		//v.count++
		//if v.count > 1 {
		//	//fmt.Println("##子查询")
		//	v.isRisk = true
		//}
		//
		//if op, ok := selectStmt.Where.(*ast.BinaryOperationExpr); ok {
		//	if op.Op == opcode.LogicOr || op.Op == opcode.LogicAnd {
		//		v.isRisk = true
		//	}
		//}
		//
		////
		////if v.sqlType == "SubSQL" {
		////	if op, ok := selectStmt.Where.(*ast.BinaryOperationExpr); ok {
		////		//where id=1 exec()
		////		if op.Op == opcode.LogicOr || op.Op == opcode.LogicAnd {
		////			//fmt.Println(reflect.TypeOf(op.R))
		////			if _, ok := op.R.(*ast.BinaryOperationExpr); ok {
		////				v.isRisk = true
		////			}
		////
		////		}
		////	}
		////}
		//
		//if selectStmt.GroupBy != nil ||
		//	selectStmt.OrderBy != nil ||
		//	selectStmt.Limit != nil ||
		//	selectStmt.WindowSpecs != nil ||
		//	selectStmt.Having != nil {
		//	v.isRisk = true
		//}

	}

	return in, false
}

func (v *SQLiDetect) Leave(in ast.Node) (ast.Node, bool) {

	//fmt.Println(reflect.TypeOf(in))
	return in, true
}

func parse(sql string) (*ast.StmtNode, error) {
	p := parser.New()

	//startTime := time.Now()

	stmtNodes, err := p.ParseOneStmt(sql, "", "")
	if err != nil {
		return nil, err
	}
	//endTime := time.Now()
	//elapsedTime := endTime.Sub(startTime)
	//fmt.Printf("代码执行耗时: %s \n", elapsedTime)

	return &stmtNodes, nil
}

func main() {
	////
	//startTime := time.Now()
	//
	//lines := Load("/Users/user/Desktop/projects/sqlparser/tests")
	//i := 0
	//for _, line := range lines {
	//	_, err := parse(line)
	//	if err != nil {
	//		//fmt.Printf("parse error: %v\n", err.Error())
	//		fmt.Println("语法错误:", line)
	//		i++
	//		continue
	//	}
	//	//v := &SQLiDetect{}
	//	//(*astNode).Accept(v)
	//	//if v.isRisk {
	//	//	//fmt.Println("发现sql注入")
	//	//} else {
	//	//	fmt.Println("漏报:", line)
	//	//}
	//}
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

	line := `select (1)`
	//line := "sleep(1)"
	astNode, err := parse(line)
	if err != nil {
		fmt.Printf("parse error: %v\n", err.Error())
		return
	} else {
		fmt.Println("语法正确")
	}

	v := &SQLiDetect{}
	(*astNode).Accept(v)
	if v.isRisk {
		fmt.Println("发现sql注入")
	} else {
		fmt.Println("漏报:", line)
	}
	//fmt.Println("语法解析错误个数:", i)
}
