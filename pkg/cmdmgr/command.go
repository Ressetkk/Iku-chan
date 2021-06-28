package cmdmgr

import (
	"context"
	"github.com/bwmarrin/discordgo"
)

// Command struct
type Command struct {
	Name        string
	Description string
	Handler     func(s *discordgo.Session, i *discordgo.InteractionCreate)
	CommandFunc func(ctx context.Context, s *discordgo.Session)
	middlewares []MiddlewareFunc
}

// Group struct
type Group struct {
	Name        string
	Description string
}

// SubCommand struct
type SubCommand struct {
	Name        string
	Description string
}
