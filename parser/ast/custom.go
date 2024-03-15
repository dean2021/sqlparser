package ast

import (
	"github.com/dean2021/sqlparser/parser/format"
)

type CommonExpressionStmt struct {
	ddlNode
	Expr ExprNode
}

// Restore implements Node interface.
func (n *CommonExpressionStmt) Restore(ctx *format.RestoreCtx) error {

	return nil
}

// Accept implements Node Accept interface.
func (n *CommonExpressionStmt) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	return v.Leave(n)
}

// WaitForStmt 由于WAITFOR不是SQL的标准语句，所以它只适用于SQL Server数据库。
type WaitForStmt struct {
	ddlNode
	DELAY string
	TIME  string
}

// Restore implements Node interface.
func (n *WaitForStmt) Restore(ctx *format.RestoreCtx) error {
	//ctx.WriteKeyWord("WAITFOR DELAY ")

	return nil
}

// Accept implements Node Accept interface.
func (n *WaitForStmt) Accept(v Visitor) (Node, bool) {
	newNode, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNode)
	}
	return v.Leave(n)
}
