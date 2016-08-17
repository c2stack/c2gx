package rc

import (
	"fmt"
	"github.com/c2stack/c2gx/process"
	"testing"
)

func TestRcParse(t *testing.T) {
	tests := []struct {
		label    string
		script   string
		expected string
	}{
		//{
			//"function",
//			"fn x {\necho 'hi'\n}\n",
//			`@
//  fn x
//`,
//		},
//		{
//			"for statement",
//			"for i\nswitch",
//			`@
//  for i in []
//    switch
//`,
//		},
//		{
//			"assignment",
//			"a = b",
//			`@
//  a = b
//`,
//		},
	}
	for _, test := range tests {
		l := lex(test.script)
		errCode := yyNewParser().Parse(l)
		label := fmt.Sprintf("%s - '%s'", test.label, test.script)
		if errCode != 0 {
			t.Errorf("%s - err %s", label, l.lastError)
		} else {
			head := l.stack.headNode()
			if head == nil {
				t.Errorf("%s - no tree", label)
			} else {
				actual := process.Dump(head)
				if test.expected != actual {
					t.Errorf("%s\nExpected:\n%s\nActual:\n%s", label, test.expected, actual)
				}
			}
		}
	}
}
