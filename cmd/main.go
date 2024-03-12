package main

import (
	"fmt"
	"github.com/dean2021/sqlparser/parser/opcode"
	"time"

	"github.com/dean2021/sqlparser/parser"
	"github.com/dean2021/sqlparser/parser/ast"
	_ "github.com/dean2021/sqlparser/parser/test_driver"
)

type SQLiDetect struct {
	isRisk bool
}

func (v *SQLiDetect) Enter(in ast.Node) (ast.Node, bool) {

	if selectStmt, ok := in.(*ast.SelectStmt); ok {
		// 子查询
		if selectStmt.AfterSetOperator != nil {
			v.isRisk = true
		}
		if op, ok := selectStmt.Where.(*ast.BinaryOperationExpr); ok {
			if op.Op == opcode.LogicOr || op.Op == opcode.LogicAnd {
				v.isRisk = true
			}

			if callStmt, ok := op.R.(*ast.FuncCallExpr); ok {
				// TODO 检查函数名是否在黑名单列表内，如果在则标识存在风险
				fmt.Println(callStmt.FnName)
			}
		}

		if selectStmt.GroupBy != nil ||
			selectStmt.OrderBy != nil ||
			selectStmt.Limit != nil ||
			selectStmt.WindowSpecs != nil ||
			selectStmt.Having != nil {
			v.isRisk = true
		}

		// 函数调用
	}

	return in, false
}

func (v *SQLiDetect) Leave(in ast.Node) (ast.Node, bool) {

	return in, true
}

func parse(sql string) (*ast.StmtNode, error) {
	p := parser.New()

	startTime := time.Now()

	stmtNodes, err := p.ParseOneStmt(sql, "", "")
	if err != nil {
		return nil, err
	}
	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)
	fmt.Printf("代码执行耗时: %s \n", elapsedTime)

	return &stmtNodes, nil
}

func main() {

	// ' OR '' = '
	// 1' ORDER BY 1--+
	// ' and 1 in (select min(name) from sysobjects where xtype = 'U' and name > '.') --
	// AND 1083=1083 AND (1427=1427
	// ;delete from xx

	//astNode, err := parse("'1' or 1=1")
	//if err != nil {
	//	fmt.Printf("parse error: %v\n", err.Error())
	//	return
	//}
	//fmt.Println(astNode)

	p := parser.New()

	stmts, _, err := p.Parse("1 or 1=1", "", "")
	if err != nil {
		fmt.Println("解析错误:", err)
		return
	}
	fmt.Println(stmts)
	//v := &SQLiDetect{}
	//(*astNode).Accept(v)
	//
	//if v.isRisk {
	//	fmt.Println("发现sql注入")
	//}

	//scanner := parser.NewScanner(`' or 1=1`)
	//scanner.LexLiteral()
	//
	//fmt.Println(scanner)
	//// 使用reflect.ValueOf获取结构体的值
	//v := reflect.ValueOf(*scanner)
	//v.FieldByName("buf").(bytes.Buffer)
	//fmt.Println()

	//fmt.Printf("%v\n", *astNode)
}
