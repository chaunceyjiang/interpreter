package object

func NewContext() *Context {
	s := make(map[string]Object)
	return &Context{store: s, stack: nil}
}

// Context 上下文环境，用来保存上下文环境变量
type Context struct {
	store map[string]Object // 这里采用map来保存
	stack *Context          // 函数栈
}

func (ctx *Context) Get(name string) (Object, bool) {
	obj, ok := ctx.store[name]
	// 优先获取自身的上下文环境
	if !ok && ctx.stack != nil {
		obj, ok = ctx.stack.Get(name) // 然后在从全局的上下文环境栈中取数
	}
	return obj, ok
}
func (ctx *Context) Set(name string, val Object) Object {
	ctx.store[name] = val
	return val
}

// NewFunctionStackContext  给每个函数增加一个函数栈
func NewFunctionStackContext(stack *Context) *Context {
	s := make(map[string]Object)
	return &Context{store: s, stack: stack}
}
