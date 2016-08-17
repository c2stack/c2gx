package rc

import (
	"github.com/dhubler/c2gx/process"
	"github.com/dhubler/c2g/c2"
)

func Run(script string, env map[string]interface{}) (map[string]interface{}, error) {
	n, err := Loads(script)
c2.Debug.Printf(process.Dump(n))
	if err != nil {
		return nil, err
	}
	return process.Run(n, env)
}

func Loads(script string) (process.CodeNode, error) {
	l := lex(script)
	errCode := yyNewParser().Parse(l)
	if errCode != 0 {
		return nil, l.lastError
	}
	if len(l.stack.code) > 0 {
		return l.stack.code[0], nil
	}
	return nil, nil
}
