package process

import (
	"fmt"
	"bytes"
)

func Dump(op Op) string {
	var buff bytes.Buffer
	op.Dump("", &buff)
	return buff.String()
}

func ArgAsString(arg interface{}) string {
	if tostr, ok := arg.(fmt.Stringer); ok {
		return tostr.String()
	}

	return fmt.Sprintf("%v", arg)
}

func ConcatFunction(stack *Stack, args ...interface{}) string {
	var buff bytes.Buffer
	for _, arg := range args {
		buff.WriteString(ArgAsString(arg))
	}
	return buff.String()
}
