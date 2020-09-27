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
