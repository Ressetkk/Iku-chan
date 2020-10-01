package dux

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
	"text/template"
)

const defaultUsageTmpl = `{{ .Name }}{{ if .Description }} - {{ .Description }}{{ end }}
{{ if .Example }}
Example: {{ .Example }}
{{ end }}
{{- if .HasSubcommands }}
Available commands:
{{ .GetCommands }}
{{- end }}
Type "{{ .Root.Name }} help [command]" to show help about provided command.`

type Commands map[string]*Command

// String returns all values from routeCol in an order
func (rc Commands) String() string {
	// TODO sorting
	var res string
	for k, v := range rc {
		res += fmt.Sprintf("%s - %s\n", k, v.Short)
	}
	return res
}

type Command struct {
	Name        string
	Run         func(ctx *Context, args []string)
	Description string
	Short       string
	Example     string
	UsageTmpl   string

	commands    Commands
	parent      *Command
	middlewares []middleware

	logger      *logrus.Entry
	helpFunc    func(ctx *Context, args []string)
	usageString string
}

func (c *Command) GetCommands() Commands {
	return c.commands
}
func (c *Command) GetCommand(name string) (*Command, bool) {
	if found, ok := c.commands[name]; ok {
		return found, true
	}
	return c, false
}

func (c *Command) AddCommand(cmd *Command) {
	if c.commands == nil {
		c.commands = make(map[string]*Command)
	}
	if c.Name == "" {
		logrus.WithField("cmd", cmd).Panic("Command name is empty!")
	}
	cmd.parent = c
	c.commands[cmd.Name] = cmd
	cmd.initHelpFunc()
	cmd.logger = logrus.WithField("command", cmd.Name)
}

func (c *Command) AddCommands(cmds ...*Command) {
	for _, cmd := range cmds {
		c.AddCommand(cmd)
	}
}

func (c *Command) Execute(ctx *Context, args []string) {
	cmd := c
	ctx.Logger = c.logger // use command's log entry to determine which command ran
	if len(args) != 0 {
		cmd, args = c.DeepFind(args)
	}
	c.logger.Debug(cmd)
	if cmd.Run != nil {
		ctx.Route = cmd
		// this might slow down the response from bot if we have a lot middleware added.
		// consider having static middleware list for each command
		for _, mw := range cmd.buildMiddlewareChain() {
			c.logger.Trace(mw)
			ctx = mw.Middleware(ctx)
			if ctx == nil {
				// we failed in one of the middlewares
				return
			}
		}
		c.logger.Trace(ctx, args)
		c.logger.Debug("exec: ", cmd)
		cmd.Run(ctx, args)
	} else if cmd.helpFunc != nil {
		cmd.helpFunc(ctx, args)
	}
}

func (c *Command) initHelpFunc() {
	if c.helpFunc == nil {
		c.helpFunc = func(ctx *Context, args []string) {
			if len(args) > 0 {
				name := args[0]
				cmd, ok := c.GetCommand(name)
				if ok {
					ctx.SendText("```\n" + cmd.usageString + "```")
					return
				}
				ctx.SendTextf("```\nCommand not found \"%s\"\n\n%s```", name, c.usageString)
				return
			}
			ctx.SendText("```\n" + c.usageString + "```")
		}
	}
}

func (c *Command) Usage() string {
	sb := &strings.Builder{}
	t := template.New("top")
	if c.UsageTmpl == "" {
		t, _ = t.Parse(defaultUsageTmpl)
	} else {
		t, _ = t.Parse(c.UsageTmpl)
	}
	t.Execute(sb, c)
	return sb.String()
}

func (c *Command) DeepFind(args []string) (*Command, []string) {
	rt := c
	for len(args) > 0 {
		name := args[0]
		foundRt, ok := rt.GetCommand(name)
		if !ok {
			return rt, args
		}
		rt = foundRt
		args = args[1:]
	}
	return rt, args
}

func (c Command) HasSubcommands() bool {
	return len(c.commands) > 0
}

func (c *Command) Root() *Command {
	if c.parent != nil {
		return c.parent.Root()
	}
	return c
}

func getPrefixAndArgs(m string) (string, []string) {
	args := strings.Split(m, " ")
	return args[0], args[1:]
}
