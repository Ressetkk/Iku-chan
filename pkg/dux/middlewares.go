package dux

func NSFWOnly(ctx *Context) *Context {
	if !ctx.Channel.NSFW {
		ctx.SendText("You can use this command only on NSFW Channel!")
		return nil
	}
	return ctx
}
