package rc

import (
	"testing"
	"strings"
)

func hereDoc(s string) string {
	return strings.Replace(s, "@BackTick@", "`", -1)
}

func TestRcScripts(t *testing.T) {
	tests := []struct{
		label string
		script string
	}{
		//{
		//	label: "assignment",
		//	script : "a = 'hello'",
		//},
//		{
//			label: "func",
//			script : `fn sayHello {
//  echo 'hello'
//}
//a = @BackTick@sayHello
//`,
//		},
	}
	for _, test := range tests {
		vars, err := Run(hereDoc(test.script), nil)
		if err != nil {
			t.Errorf("test:%s - ERR=%v", test.label, err)
		} else {
			actual, hasVar := vars["a"]
			if ! hasVar {
				t.Errorf("test:%s - variable 'a' not found", test.label)
			} else if actual != "hello" {
				t.Errorf("test:%s\nExpected:'hello'\n  Actual:'%s'", test.label, actual)
			}
		}

	}
}
