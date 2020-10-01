package main

import (
	"flag"
	"github.com/Ressetkk/Iku-chan/cmd/nhentai"
	"github.com/Ressetkk/Iku-chan/pkg/dux"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

var (
	logLevel = flag.String("logLevel", "warn", "Set log level.")
)

func main() {
	flag.Parse()
	logrus.SetFormatter(&logrus.JSONFormatter{})
	token := os.Getenv("IKU_BOT_TOKEN")
	if token == "" {
		logrus.Fatal("IKU_BOT_TOKEN cannot be empty. Exiting.")
	}
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		logrus.WithError(err).Fatal("Could not initialize Discord session")
	}
	switch *logLevel {
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
		session.LogLevel = discordgo.LogError
		break
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
		session.LogLevel = discordgo.LogWarning
		break
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
		session.LogLevel = discordgo.LogWarning
		break
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
		session.LogLevel = discordgo.LogInformational
		break
	case "trace":
		logrus.SetLevel(logrus.TraceLevel)
		session.LogLevel = discordgo.LogDebug
		break
	default:
		logrus.Fatal("wrong \"logLevel\" value: %s", *logLevel)
	}

	if err := session.Open(); err != nil {
		logrus.WithError(err).Fatal("Could not open Discord Bot session")
	}

	logrus.Info("Discord session initialized. Initializing commands.")

	r := &dux.Command{
		Name: "iku",
		Description: `The bot for most perverted and thirsty degenerates.
Come and use me, senpai~`,
	}

	r.AddCommands(nhentai.GetCmd(), nhentai.SearchCmd(), nhentai.RandomCmd())
	handler := dux.Handler{
		AllowMentions: true,
		Root:          r,
	}
	session.AddHandler(handler.Set())

	defer func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Kill, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		s := <-sig
		logrus.Infof("Requested %v. Exiting.", s)
		_ = session.Close()
	}()
}
