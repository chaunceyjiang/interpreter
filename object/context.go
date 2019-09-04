package object

func NewContext() *Context {
	s := make(map[string]Object)
	return &Context{store: s}
}

// Context 上下文环境，用来保存上下文环境变量
type Context struct {
	store map[string]Object // 这里采用map来保存
}

func (ctx *Context) Get(name string) (Object, bool) {
	obj, ok := ctx.store[name]
	return obj, ok
}
func (ctx *Context) Set(name string, val Object) Object {
	ctx.store[name] = val
	return val
}
