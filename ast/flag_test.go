// Copyright 2016 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package ast

import (
	"github.com/dean2021/sqlparser"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHasAggFlag(t *testing.T) {
	expr := &BetweenExpr{}
	flagTests := []struct {
		flag   uint64
		hasAgg bool
	}{
		{FlagHasAggregateFunc, true},
		{FlagHasAggregateFunc | FlagHasVariable, true},
		{FlagHasVariable, false},
	}
	for _, tt := range flagTests {
		expr.SetFlag(tt.flag)
		require.Equal(t, tt.hasAgg, HasAggFlag(expr))
	}
}

func TestFlag(t *testing.T) {
	flagTests := []struct {
		expr string
		flag uint64
	}{
		{
			"1 between 0 and 2",
			FlagConstant,
		},
		{
			"case 1 when 1 then 1 else 0 end",
			FlagConstant,
		},
		{
			"case 1 when 1 then 1 else 0 end",
			FlagConstant,
		},
		{
			"case 1 when a > 1 then 1 else 0 end",
			FlagConstant | FlagHasReference,
		},
		{
			"1 = ANY (select 1) OR exists (select 1)",
			FlagHasSubquery,
		},
		{
			"1 in (1) or 1 is true or null is null or 'abc' like 'abc' or 'abc' rlike 'abc'",
			FlagConstant,
		},
		{
			"row (1, 1) = row (1, 1)",
			FlagConstant,
		},
		{
			"(1 + a) > ?",
			FlagHasReference | FlagHasParamMarker,
		},
		{
			"trim('abc ')",
			FlagHasFunc,
		},
		{
			"now() + EXTRACT(YEAR FROM '2009-07-02') + CAST(1 AS UNSIGNED)",
			FlagHasFunc,
		},
		{
			"substring('abc', 1)",
			FlagHasFunc,
		},
		{
			"sum(a)",
			FlagHasAggregateFunc | FlagHasReference,
		},
		{
			"(select 1) as a",
			FlagHasSubquery,
		},
		{
			"@auto_commit",
			FlagHasVariable,
		},
		{
			"default(a)",
			FlagHasDefault,
		},
		{
			"a is null",
			FlagHasReference,
		},
		{
			"1 is true",
			FlagConstant,
		},
		{
			"a in (1, count(*), 3)",
			FlagConstant | FlagHasReference | FlagHasAggregateFunc,
		},
		{
			"'Michael!' REGEXP '.*'",
			FlagConstant,
		},
		{
			"a REGEXP '.*'",
			FlagHasReference,
		},
		{
			"-a",
			FlagHasReference,
		},
	}
	p := sqlparser.New()
	for _, tt := range flagTests {
		stmt, err := p.ParseOneStmt("select "+tt.expr, "", "")
		require.NoError(t, err)
		selectStmt := stmt.(*SelectStmt)
		SetFlag(selectStmt)
		expr := selectStmt.Fields.Fields[0].Expr
		require.Equalf(t, tt.flag, expr.GetFlag(), "For %s", tt.expr)
	}
}
