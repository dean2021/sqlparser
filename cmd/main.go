package main

import (
	"fmt"
	"github.com/dean2021/sqlparser"
	"github.com/dean2021/sqlparser/ast"
	"github.com/dean2021/sqlparser/opcode"
	"github.com/dean2021/sqlparser/test_driver"
	_ "github.com/dean2021/sqlparser/test_driver"
	"reflect"
	"strings"
)

type SQLiDetect struct {
	isRisk bool
	count  int
}

func (v *SQLiDetect) Enter(in ast.Node) (ast.Node, bool) {
	if _, ok := in.(*ast.BRIEStmt); ok {
		v.isRisk = true
	}
	if _, ok := in.(*ast.LoadDataStmt); ok {
		v.isRisk = true
	}
	if _, ok := in.(*ast.CreateTableStmt); ok {
		v.isRisk = true
	}
	if _, ok := in.(*ast.WaitForStmt); ok {
		v.isRisk = true
	}
	if ce, ok := in.(*ast.CommonExpressionStmt); ok {
		//fmt.Println(reflect.TypeOf(ce.Expr))
		if op, ok := ce.Expr.(*ast.FuncCallExpr); ok {
			strArray := []string{"sleep", "benchmark", "pg_sleep", "exec", "randomblob", "substring", "lower", "ascii", "version", "databases"}
			for _, value := range strArray {
				if value == strings.ToLower(op.FnName.String()) {
					v.isRisk = true
					break
				}
			}
		}
		if op, ok := ce.Expr.(*ast.BinaryOperationExpr); ok {
			//fmt.Println(reflect.TypeOf(op.Op))
			if op.L != nil {
				if op, ok := op.L.(*ast.BinaryOperationExpr); ok {

					if op.Op == opcode.LogicOr || op.Op == opcode.LogicAnd {
						v.isRisk = true
					}
				}
			}

			fmt.Println("一级:", reflect.TypeOf(op.L), op.Op, reflect.TypeOf(op.R))

			fmt.Println("reflect.TypeOf(op.R)", reflect.TypeOf(op.R))
			fmt.Println(op.R.Text())
			if l, ok := op.L.(*ast.BinaryOperationExpr); ok {
				fmt.Println("	L:", l.L.Text(), l.Op, l.R.Text())
			} else if l, ok := op.L.(*test_driver.ValueExpr); ok {
				fmt.Println("	V:", l.GetValue())
			}
			if r, ok := op.R.(*ast.BinaryOperationExpr); ok {
				fmt.Println("	R:", r.L.Text(), r.Op, r.R.Text())
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

		//if val, ok := ce.Expr.(*test_driver.ValueExpr); ok {
		//	fmt.Println("字符串:", val.GetValue())
		//}
		//if val, ok := ce.Expr.(*test_driver.ValueExpr); ok {
		//	//fmt.Println(reflect.TypeOf(ce.Expr))
		//	//fmt.Println(val.GetValue())
		//	//if str, ok := val.GetValue().(string); ok {
		//	//if check(str) {
		//	v.isRisk = true
		//	//}
		//	//}
		//}

		//if v.isRisk == false {
		//	fmt.Println(reflect.TypeOf(ce.Expr))
		//}
	}

	if s, ok := in.(*ast.SelectStmt); ok {
		for _, field := range s.Fields.Fields {
			if call, ok := field.Expr.(*ast.FuncCallExpr); ok {
				fmt.Println(call.FnName, call.Args)
			}
		}
		v.isRisk = true
	}

	//if val, ok := in.(*test_driver.ValueExpr); ok {
	//
	//	if str, ok := val.GetValue().(string); ok {
	//		if check(str) {
	//			v.isRisk = true
	//		}
	//	}
	//}

	return in, false
}

func (v *SQLiDetect) Leave(in ast.Node) (ast.Node, bool) {

	//fmt.Println(reflect.TypeOf(in))
	return in, true
}

func parse(sql string) (*ast.StmtNode, error) {
	p := sqlparser.New()
	p.EnableWindowFunc(false)
	stmtNodes, err := p.ParseOneStmt(sql, "", "")
	if err != nil {
		return nil, err
	}
	return &stmtNodes, nil
}

//	func check(sql string) bool {
//		// 太短就不检查了
//		if len(sql) < 6 {
//			return false
//		}
//		astNode, err := parse(sql)
//		if err != nil {
//			code := test3.Fix(sql)
//			astNode, err := parse(code)
//			if err != nil {
//				fmt.Println("语法错误:", sql)
//				fmt.Println("修复语法:", code)
//				return false
//			}
//			v := &SQLiDetect{}
//			(*astNode).Accept(v)
//			return v.isRisk
//		}
//		v := &SQLiDetect{}
//		(*astNode).Accept(v)
//		return v.isRisk
//	}

func main() {

	//startTime := time.Now()
	//lines := Load("/Users/user/Desktop/projects/sqlparser/tests")
	//
	//for _, line := range lines {
	//	if len(line) < 6 {
	//		continue
	//	}
	node, err := parse(`select * from syscat.tabauth where grantee = current `)
	//node, err := parse(`'string' and substring(password/textz(),1,1)='string'`)
	//node, err := parse(`select substring(1)`)
	//node, err := parse(`1 union select 1;`)
	//node, err := parse(`x union select 1;`)
	//node, err := parse(`BACKUP database master to disk='stringx'`)
	//node, err := parse(`'string' union select 1; `)
	//node, err := parse(`create table myfile(x TEXT)`)
	//select substring(1)
	//select substrings(1)
	if err != nil {
		fmt.Println("语法错误:", err)
		return
	}
	v := &SQLiDetect{}
	(*node).Accept(v)
	if v.isRisk {
		fmt.Println("发现sql注入")
	}
}
