package rc

import (
	"testing"
	"fmt"
)

func TestRcParse(t *testing.T) {
	tests := []struct{
		label string
		script string
		expected string
	} {
		{
			"function",
			"fn x",
			"fn x",
		},
		{
			"for statement",
			"for ( i ) switch",
			"for (i in <nil>) switch",
		},
		{
			"assigment",
			"a = b",
			"a = b",
		},
	}
	for _, test := range tests {
		l := lex(test.script)
		errCode := yyNewParser().Parse(l)
		label := fmt.Sprintf("%s - '%s'", test.label, test.script)
		if errCode != 1 {
			t.Errorf("%s - err %d", label, errCode)
		} else  if l.tree == nil {
			t.Errorf("%s - no tree", label)
		} else {
			actual := l.tree.String()
			if test.expected != actual {
				t.Errorf("%s\nExpected:%s\n  Actual:%s", label, test.expected, actual)
			}
		}
	}
}