package process

func Run(n CodeNode, env map[string]interface{}) (map[string]interface{}, error) {
	stack := NewGlobalStack(env)
	err := n.Exec(stack)
	return stack.Vars, err
}
