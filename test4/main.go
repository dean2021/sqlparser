package main

import (
	"github.com/pingcap/parser"
	_ "github.com/pingcap/tidb/pkg/types/parser_driver"
)

func main() {

	sql := `select user() from x`
	p := parser.New()
	_, err := p.ParseOneStmt(sql, "", "")
	if err != nil {
		panic(err)
	}
}
