package test

import "github.com/Ressetkk/Iku-chan/pkg/router"

func HelloWorldCmd() *router.Command {
	c := &router.Command{
		Use:         "helloworld",
		Aliases:     []string{"hello", "world"},
		Description: "prints hello world to the world!",
		Short:       "print hello world",
		Example:     "helloworld",
		Run: func(h router.Payload) {
			h.SendText("Za warudo!")
		},
	}
	c.AddCommands(&router.Command{
		Use: "asd",
		Run: func(h router.Payload) {
			h.SendText("P2")
		},
	})
	c.AddMiddleware(func(m router.Payload) router.Payload {
		m.SendText("helloworld middleware")
		return m
	})
	return c
}
