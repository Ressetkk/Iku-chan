package dux

// MiddlewareFunc is a type that defines middleware function.
type MiddlewareFunc func(ctx *Context) *Context

type middleware interface {
	Middleware(ctx *Context) *Context
}

// Middleware executes MiddlewareFunc on a provided Context.
func (mf MiddlewareFunc) Middleware(ctx *Context) *Context {
	return mf(ctx)
}

// AddMiddleware adds defined middlewares to the Command's middleware chain.
func (c *Command) AddMiddleware(mdf ...MiddlewareFunc) {
	for _, mf := range mdf {
		c.middlewares = append(c.middlewares, mf)
	}
}

// HasMiddleware returns true if the command has middlewares assigned.
func (c Command) HasMiddleware() bool {
	return len(c.middlewares) != 0
}

// buildMiddlewareChain returns a combined middleware chain for a Command with all parent Commands' middlewares.
func (c *Command) buildMiddlewareChain() []middleware {
	cmd := c
	mw := c.middlewares
	for cmd.parent != nil {
		cmd = cmd.parent
		if cmd.HasMiddleware() {
			mw = append(cmd.middlewares, mw...)
		}
	}
	return mw
}
