package dux

import (
	"github.com/bwmarrin/discordgo"
	"sync"
)

type Options struct {
	Aliases       []string
	AllowMentions bool
	HelpCmd       *Command
}

func (c *Command) renderUsages() {
	cmds := c.collectAllCommands([]*Command{c})
	cmdChan := make(chan *Command, 100)
	var wg sync.WaitGroup
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				cmd, ok := <-cmdChan
				if !ok {
					return
				}
				cmd.usageString = cmd.Usage()
			}
		}()
	}
	for len(cmds) > 0 {
		cmdChan <- cmds[0]
		cmds = cmds[1:]
	}
	close(cmdChan)
	wg.Wait()
}

func (c *Command) collectAllCommands(cmds []*Command) []*Command {
	if len(c.commands) > 0 {
		for _, cmd := range c.commands {
			cmds = append(cmds, cmd)
			if len(cmd.commands) > 0 {
				cmds = cmd.collectAllCommands(cmds)
			}
		}
	}
	return cmds
}

func (c *Command) Handler(o Options) func(s *discordgo.Session, m *discordgo.MessageCreate) {
	if o.HelpCmd == nil {
		c.AddCommand(&Command{
			Name: "help",
			Description: `Shows help for provided command.
The command accepts the path to a command starting from subcommand of root.
"Eg. root help sub1 sub2" will return help command for "sub2"`,
			Example: "help subcommand1 subcommand2",
			Short:   "Shows help for provided command.",
			Run: func(ctx *Context, args []string) {
				cmd, args := c.DeepFind(args)
				cmd.helpFunc(ctx, args)
			},
		})
	} else {
		c.AddCommand(o.HelpCmd)
	}
	c.initHelpFunc()
	c.renderUsages()
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
