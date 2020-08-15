package router

// MiddlewareFunc is a type that defines middleware function.
type MiddlewareFunc func(m Payload) Payload

type middleware interface {
	Middleware(m Payload) Payload
}

// Middleware executes MiddlewareFunc on a provided payload.
func (mf MiddlewareFunc) Middleware(m Payload) Payload {
	return mf(m)
}

// AddMiddleware adds defined middlewares to the command's middleware chain.
func (c *Command) AddMiddleware(mdf ...MiddlewareFunc) {
	for _, mf := range mdf {
		c.middlewares = append(c.middlewares, mf)
	}
}

// HasMiddleware returns true if the command has middlewares assigned.
func (c Command) HasMiddleware() bool {
	return len(c.middlewares) != 0
}

// BuildMiddlewareChain returns a combined middleware chain for a command with all parent commands' middlewares.
func (c *Command) BuildMiddlewareChain() []middleware {
	cmd := c
	mw := c.middlewares
	for cmd.HasParent() {
		cmd = c.parent
		if cmd.HasMiddleware() {
			mw = append(cmd.middlewares, mw...)
		}
	}
	return mw
}
