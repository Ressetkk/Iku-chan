package router

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
)

type CommandHandler func(h Payload)

type Command struct {
	Use         string
	Aliases     []string
	Description string
	Short       string
	Example     string
	Version     string
	Run         CommandHandler

	commands    []*Command
	parent      *Command
	middlewares []middleware
	// TODO pflag kinda works but we need the default values to be permanent.
}

type Options struct {
	AllowMentions bool
	IgnoreCases   bool
}

// Handler returns handler function that implements MessageCreate discordgo.EventHandler.
// It accepts Options struct with Handler options.
func (c Command) Handler(o Options) func(s *discordgo.Session, m *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		message := m.Message
		splitMessage := strings.Split(message.Content, " ")
		prefix := splitMessage[0]
		if m.Author.ID == s.State.User.ID {
			return
		}
		// TODO I don't like it. Creating new map on every event seems inefficient.
		availablePrefixes := append(c.Aliases, c.Use)
		shouldRespond := ifShouldRespond(prefix, s.State.User.ID, availablePrefixes, o.AllowMentions)
		if !shouldRespond {
			return
		}
		msgHandler := Payload{Session: s, Event: m}
		c.Execute(splitMessage[1:], msgHandler)
	}
}

// Execute executes command with argument and a Discord message Payload.
func (c *Command) Execute(args []string, h Payload) {
	if len(args) != 0 {
		// TODO flag parsing
		cmd := c.find(args)
		for _, mw := range cmd.BuildMiddlewareChain() {
			h = mw.Middleware(h)
		}
		cmd.Run(h)
	}
}

// FindSubcommand returns command based on a provided name.
// It checks all subcommands' names and aliases. If first result matches the criteria it returns the command and exits.
// If no commands found then returns an error.
func (c Command) FindSubcommand(name string) (*Command, error) {
	for _, cmd := range c.commands {
		if cmd.Use == name {
			return cmd, nil
		}
		if cmd.HasAliases() {
			for _, v := range cmd.Aliases {
				if v == name {
					return cmd, nil
				}
			}
		}
	}
	return nil, fmt.Errorf("command not found %v", name)
}

func (c *Command) find(commands []string) *Command {
	if len(commands) == 0 {
		return c
	}
	cmd, err := c.FindSubcommand(commands[0])
	if err != nil {
		return c
	}
	if cmd.HasSubcommands() {
		cmd = cmd.find(commands[1:])
	}
	return cmd
}

// AddCommands adds provided commands to the Command tree. After adding it performs sorting based on commands' name.
func (c *Command) AddCommands(cmds ...*Command) {
	for i, cmd := range cmds {
		if cmds[i] == c {
			// TODO maybe try to handle it safely with proper error or warning and skip the command
			panic("command cannot be a child of itself")
		}
		cmd.parent = c
		c.commands = append(c.commands, cmd)
		//TODO sort commands
	}
}

// HasAliases returns true if the command has aliases.
func (c Command) HasAliases() bool {
	return c.Aliases != nil
}

// HasSubcommands returns true if the command has subcommands.
func (c Command) HasSubcommands() bool {
	return c.commands != nil
}

// HasParent returns true if the command has parent command.
func (c Command) HasParent() bool {
	return c.parent != nil
}

// IsRunnable returns true if the command has defined Run function, e.g. is runnable.
func (c Command) IsRunnable() bool {
	return c.Run != nil
}

// TODO rewrite it
func ifShouldRespond(prefix, userID string, availablePrefixes []string, allowMentions bool) bool {
	shouldRespond := false
	if "<@!"+userID+">" == prefix && allowMentions {
		shouldRespond = true
	}
	for _, v := range availablePrefixes {
		if v == prefix {
			shouldRespond = true
		}
	}
	return shouldRespond
}
