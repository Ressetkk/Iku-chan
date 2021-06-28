package cmdmgr

import "github.com/bwmarrin/discordgo"

type Router struct {
	Commands    []*Command
	middlewares []MiddlewareFunc
}

func (r Router) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	name := i.ApplicationCommandData().Name
	ctx := NewContext(s, i)

	for _, c := range r.Commands {
		if c.Name == name {
			c.CommandFunc(ctx, s)
		}
	}
}
