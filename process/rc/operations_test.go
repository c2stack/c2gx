package rc

import (
	"testing"
)

func TestProcessOperations(t *testing.T) {
	op, err := RcScript(`a=hello
`)
	if err != nil {
		t.Error(err.Error())
	} else {
		t.Log(op)
	}
}
