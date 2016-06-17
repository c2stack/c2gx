package process

type Value struct {
	Var  string
	Str  string
	Num  string
	Func *FuncCallOp
}

func (self *Value) Eval(stack *Stack) (interface{}, error) {
	if len(self.Var) > 0 {
		return stack.Vars[self.Var], nil
	} else if len(self.Str) > 0 {
		return self.Str, nil
	}
	// decide how to return number
	return self.Num, nil
}

func (self *Value) String() string {
	if len(self.Var) > 0 {
		return self.Var
	} else if len(self.Str) > 0 {
		return self.Str
	}
	return self.Num
}

func Variable(ident string) *Value {
	return &Value{Var: ident}
}
