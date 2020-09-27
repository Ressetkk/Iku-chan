package dux

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
)

type Commands map[string]*Command

// String returns all values from routeCol in an order
func (rc Commands) String() string {
	// TODO sorting
	var res string
	for k, v := range rc {
		res += fmt.Sprintf("%s	%s\n", k, v.Short)
	}
	return res
}

type Command struct {
	Name        string
	Run         func(ctx *Context, args []string)
	Description string
	Short       string

	commands    Commands
	parent      *Command
	helpFunc    func(ctx *Context)
	middlewares []middleware
}

type Options struct {
	Aliases       []string
	AllowMentions bool
}

func (c *Command) GetRoutes() Commands {
	return c.commands
}
func (c *Command) GetRoute(name string) (*Command, bool) {
	if found, ok := c.commands[name]; ok {
		return found, true
	}
	return c, false
}

func (c *Command) AddCommand(cmd *Command) {
	if c.commands == nil {
		c.commands = make(map[string]*Command)
	}
	cmd.parent = c
	c.commands[cmd.Name] = cmd
}

func (c *Command) AddCommands(cmds ...*Command) {
	for _, cmd := range cmds {
		c.AddCommand(cmd)
	}
}

func (c *Command) Execute(ctx *Context, args []string) {
	rt := c
	if len(args) != 0 {
		rt, args = c.DeepFind(args)
	}

	if rt.Run != nil {
		ctx.Route = rt
		for _, mw := range rt.buildMiddlewareChain() {
			ctx = mw.Middleware(ctx)
			if ctx == nil {
				return
			}
		}
		rt.Run(ctx, args)
	}
}

func (c *Command) DeepFind(args []string) (*Command, []string) {
	rt := c
	for len(args) > 0 {
		name := args[0]
		foundRt, ok := rt.GetRoute(name)
		if !ok {
			return rt, args
		}
		rt = foundRt
		args = args[1:]
	}
	return rt, args
}

func (c *Command) Handler(o Options) func(s *discordgo.Session, m *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		prefix, args := getPrefixAndArgs(m.Content)
	check:
		switch {
		case m.Author.ID == s.State.User.ID:
			return
		case o.AllowMentions && "<@!"+s.State.User.ID+">" == prefix:
			break
		case prefix != c.Name:
			for _, alias := range o.Aliases {
				if alias == prefix {
					break check
				}
			}
			return
		}

		ctx, err := NewContext(s, m)
		if err != nil {
			return
		}
		c.Execute(ctx, args)
	}
}

func getPrefixAndArgs(m string) (string, []string) {
	args := strings.Split(m, " ")
	return args[0], args[1:]
}
