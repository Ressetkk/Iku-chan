package dux

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"sync"
)

type Handler struct {
	Aliases       []string
	AllowMentions bool
	HelpCmd       *Command
	Root          *Command
}

func (h *Handler) Set() func(s *discordgo.Session, m *discordgo.MessageCreate) {
	if h.HelpCmd == nil {
		h.Root.AddCommand(&Command{
			Name: "help",
			Description: `Shows help for provided command.
The command accepts the path to a command starting from subcommand of root.
"Eg. root help sub1 sub2" will return help command for "sub2"`,
			Example: "help subcommand1 subcommand2",
			Short:   "Shows help for provided command.",
			Run: func(ctx *Context, args []string) {
				cmd, args := h.Root.DeepFind(args)
				cmd.helpFunc(ctx, args)
			},
		})
	} else {
		h.Root.AddCommand(h.HelpCmd)
	}
	h.Root.initHelpFunc()
	h.Root.logger = logrus.WithField("command", h.Root.Name)
	renderUsages(h.Root)

	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		prefix, args := getPrefixAndArgs(m.Content)
		logrus.Trace(prefix, args)
	check:
		switch {
		case m.Author.ID == s.State.User.ID:
			return
		case h.AllowMentions && "<@!"+s.State.User.ID+">" == prefix:
			break
		case prefix != h.Root.Name:
			for _, alias := range h.Aliases {
				if alias == prefix {
					break check
				}
			}
			return
		}

		ctx, err := NewContext(s, m)
		logrus.Trace(ctx, err)
		if err != nil {
			logrus.WithError(err).Error("could not create dux Context")
			return
		}
		h.Root.Execute(ctx, args)
	}
}

func renderUsages(c *Command) {
	cmds := collectAllCommands(c, []*Command{c})
	cmdChan := make(chan *Command, 100)
	var wg sync.WaitGroup
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				cmd, ok := <-cmdChan
				if !ok {
					logrus.Trace("job done for worker")
					return
				}
				logrus.Trace("picked", cmd)
				cmd.usageString = cmd.Usage()
				logrus.Debug("done rendering", cmd.Name)
				logrus.Trace(cmd.usageString)
			}
		}()
	}
	for len(cmds) > 0 {
		cmdChan <- cmds[0]
		cmds = cmds[1:]
	}
	close(cmdChan)
	wg.Wait()
	logrus.Info("Finished generating Usage strings.")
}

func collectAllCommands(prev *Command, cmds []*Command) []*Command {
	if len(prev.commands) > 0 {
		for _, cmd := range prev.commands {
			cmds = append(cmds, cmd)
			if len(cmd.commands) > 0 {
				cmds = collectAllCommands(cmd, cmds)
			}
		}
	}
	return cmds
}
