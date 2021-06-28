package cmdmgr

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"sync"
)

type Context struct {
	context.Context

	Logger *logrus.Entry

	done chan struct{}
	mu   sync.Mutex
	s    *discordgo.Session
	i    *discordgo.InteractionCreate
}

func (c *Context) Value(key interface{}) interface{} {
	c.mu.Lock()
	defer c.mu.Unlock()
	ch := key.(string)
	switch ch {
	case "session":
		return c.s
	case "interaction":
		return c.i
	}
	return nil
}

func NewContext(s *discordgo.Session, i *discordgo.InteractionCreate) *Context {
	return &Context{
		s: s,
		i: i,
	}
}
