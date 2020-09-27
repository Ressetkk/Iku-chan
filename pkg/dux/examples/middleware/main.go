package main

import (
	"github.com/Ressetkk/Iku-chan/pkg/dux"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		logrus.Fatal("BOT_TOKEN cannot be empty. Exiting.")
	}
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		logrus.WithError(err).Fatal("Could not initialize Discord session")
	}
	session.LogLevel = discordgo.LogWarning
	if err := session.Open(); err != nil {
		logrus.WithError(err).Fatal("Could not open Discord Bot session")
	}
	logrus.Info("Discord session initialized. Initializing commands.")

	// Initialize root command
	r := &dux.Command{
		Name: "hello",
		Run: func(ctx *dux.Context, args []string) {
			ctx.SendTextf("```Hello World!. I got args: %v\n```", args)
		}}

	// add middleware for root command
	r.AddMiddleware(func(ctx *dux.Context) *dux.Context {
		if !ctx.Channel.NSFW {
			ctx.SendText("You can only use this command in NSFW channel!")
			return nil
		}
		ctx.SendText("Pervert...")
		return ctx
	})

	// add subcommand and it's own middleware
	subc := &dux.Command{
		Name: "funky",
		Run: func(ctx *dux.Context, args []string) {
			ctx.SendText("Now we're funky!")
		},
	}
	r.AddCommand(subc)
	subc.AddMiddleware(func(ctx *dux.Context) *dux.Context {
		ctx.SendText("We're not even funky there!")
		return ctx
	})

	// Initialize handler options
	opts := dux.Options{
		AllowMentions: true,
		Aliases:       []string{"world", "helloWorld"},
	}

	// add DG0 handler for a root command with options
	session.AddHandler(r.Handler(opts))

	defer func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Kill, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		s := <-sig
		logrus.Infof("Requested %v. Exiting.", s)
		_ = session.Close()
	}()
}
