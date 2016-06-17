package process

var builtins map[string]interface{}
func init() {
	builtins = map[string]interface{} {
		"echo" : func(stack *Stack) int {
			stack.Out.WriteString("Stuff")
			return 0
		},
	}
}
