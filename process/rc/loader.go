package rc

import (
	"github.com/c2gx/process"
	"github.com/c2g/c2"
	"fmt"
)

func RcScript(script string) (process.Op, error) {
	l := lex(script)
	errCode := yyNewParser().Parse(l)
	if errCode > 1 {
		return nil, c2.NewErr(fmt.Sprintf("Could not parse script, err=%d", errCode))
	}
	return l.tree, nil
}
