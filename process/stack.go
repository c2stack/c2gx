package process

import "bytes"

type Stack struct {
	Global *Stack
	Env     map[string]interface{}
	Vars    map[string]interface{}
	Out     bytes.Buffer
	Parent  *Stack
	Funcs   map[string]*FuncDefOp
	LastErr error
}

func NewGlobalStack(env map[string]interface{}) *Stack {
	g := &Stack{
		Funcs: make(map[string]*FuncDefOp),
		Env: env,
	}
	g.Global = g
	return g
}

func (stack *Stack) NewStack() *Stack {
	return &Stack{
		Parent: stack,
		Global: stack.Global,
	}
}

func (stack *Stack) Assign(key string, value interface{}) {
	if stack.Vars == nil {
		stack.Vars = make(map[string]interface{}, 1)
	}
	stack.Vars[key] = value
}

func (stack *Stack) SetFunc(key string, fdef *FuncDefOp) {
	stack.Global.Funcs[key] = fdef
}

func (stack *Stack) Func(key string) (fdef *FuncDefOp, found bool) {
	fdef, found = stack.Global.Funcs[key]
	return
}

func (stack *Stack) ResolveValue(valname string) interface{} {
	if v, hasValue := stack.Vars[valname]; hasValue {
		return v
	}
	return nil
}

func (stack *Stack) ClearLastError() error {
	e := stack.LastErr
	stack.LastErr = nil
	return e
}
