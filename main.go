package main

import (
	"github.com/Ressetkk/Iku-chan/cmd/nhentai"
	"github.com/Ressetkk/Iku-chan/pkg/dux"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	token := os.Getenv("IKU_BOT_TOKEN")
	if token == "" {
		logrus.Fatal("IKU_BOT_TOKEN cannot be empty. Exiting.")
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

	r := &dux.Command{
		Name: "iku",
		Description: `The bot for most perverted and thirsty degenerates.
Come and use me, senpai~`,
	}

	r.AddCommands(nhentai.GetCmd(), nhentai.SearchCmd())

	opts := dux.Options{AllowMentions: true}
	session.AddHandler(r.Handler(opts))

	defer func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Kill, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		s := <-sig
		logrus.Infof("Requested %v. Exiting.", s)
		_ = session.Close()
	}()
}
