// Copyright 2023 PingCAP, Inc.
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

package sqlparser_test

import (
	"github.com/dean2021/sqlparser"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestKeywords(t *testing.T) {
	// Test for the first keyword
	require.Equal(t, "ADD", sqlparser.Keywords[0].Word)
	require.Equal(t, true, sqlparser.Keywords[0].Reserved)

	// Make sure TiDBKeywords are included.
	found := false
	for _, kw := range sqlparser.Keywords {
		if kw.Word == "ADMIN" {
			found = true
		}
	}
	require.Equal(t, found, true, "TiDBKeyword ADMIN is part of the list")
}

func TestKeywordsLength(t *testing.T) {
	require.Equal(t, 644, len(sqlparser.Keywords))

	reservedNr := 0
	for _, kw := range sqlparser.Keywords {
		if kw.Reserved {
			reservedNr += 1
		}
	}
	require.Equal(t, 233, reservedNr)
}

func TestKeywordsSorting(t *testing.T) {
	for i, kw := range sqlparser.Keywords {
		if i > 1 && sqlparser.Keywords[i-1].Word > kw.Word && sqlparser.Keywords[i-1].Section == kw.Section {
			t.Errorf("%s should come after %s, please update parser.y and re-generate keywords.go\n",
				sqlparser.Keywords[i-1].Word, kw.Word)
		}
	}
}
