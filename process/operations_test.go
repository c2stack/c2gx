package process

import "testing"

func TestAssignOp(t *testing.T) {
	script := &SubShell{}
	script.AddOp(&AssignOp{Var:"a", Val: &Value{Str:"hello"}})
	stack := &Stack{}
	err := script.Exec(stack)
	if err != nil {
		t.Error(err)
	}
	if actual, hasValue := stack.Vars["a"]; !hasValue {
		t.Error("no value")
	} else if actual != "hello" {
		t.Errorf("\nExpected:hello\n  Actual:%s", actual)
	}
}
