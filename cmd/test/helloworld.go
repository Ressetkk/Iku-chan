package test

import (
	"github.com/Ressetkk/Iku-chan/pkg/dux"
)

func HelloWorldCmd() *dux.Command {
	c := &dux.Command{
		Name:        "helloworld",
		Description: "prints hello world to the world!",
		Short:       "print hello world",
		Run: func(ctx *dux.Context, args []string) {
			ctx.SendText("aaaaaa")
		},
	}
	c.AddCommand(&dux.Command{
		Name: "asd",
		Run: func(ctx *dux.Context, args []string) {
			ctx.SendText("hentai")
		},
	})
	c.AddMiddleware(func(ctx *dux.Context) *dux.Context {
		if !ctx.Channel.NSFW {
			ctx.SendText("You can only use this command on NSFW channel!")
			return nil
		}
		return ctx
	})
	return c
}
