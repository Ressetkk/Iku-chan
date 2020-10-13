package activity

import (
	"fmt"
	"github.com/Ressetkk/Iku-chan/internal/embed"
	"github.com/Ressetkk/Iku-chan/pkg/dux"
	"github.com/Ressetkk/Iku-chan/pkg/nhapi"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"time"
)

type HentaiActivity struct {
	Client  *nhapi.Client
	Session *discordgo.Session
	Logger  *logrus.Entry

	numberOfTheDay int
	ticker         *time.Ticker
}

func (ha *HentaiActivity) TodayCmd() *dux.Command {
	cmd := &dux.Command{
		Name: "today",
		Run: func(ctx *dux.Context, args []string) {

			res, err := ha.Client.Get(ha.numberOfTheDay)
			if err != nil {
				ctx.SendText("Could not fetch today's Hentai...")
				ctx.Logger.WithError(err).Error("Could not fetch today's Hentai...")
				return
			}
			eb := embed.Make(res)
			ms := discordgo.MessageSend{
				Content: "Here's your today's dose of spice.",
				Embed:   &eb,
			}
			ctx.Send(&ms)
		},
		Description: `Returns today's hentai. This changes on daily basis so check it out whenever you can!`,
		Short:       "Returns today's hentai",
		Example:     "today",
	}

	cmd.AddMiddleware(dux.NSFWOnly)
	return cmd
}

func (ha *HentaiActivity) Update() {
	err := ha.updateStatus()
	if err != nil {
		ha.Logger.WithError(err).Error("could not update activity")
	}
}

func (ha *HentaiActivity) updateStatus() error {
	num, err := ha.Client.Random()
	if err != nil {
		return fmt.Errorf("nhapi error: %w", err)
	}
	ha.numberOfTheDay = num
	ud := discordgo.UpdateStatusData{
		Game: &discordgo.Game{
			Name: fmt.Sprintf("%v", ha.numberOfTheDay),
			Type: discordgo.GameTypeWatching,
			URL:  fmt.Sprintf("https://nhentai.net/g/%v", ha.numberOfTheDay),
		},
		Status: "online",
	}
	return ha.Session.UpdateStatusComplex(ud)
}

func (ha *HentaiActivity) Run() {
	ha.Update()
	if ha.ticker == nil {
		ha.ticker = time.NewTicker(time.Hour * 24)
	}
	for {
		select {
		case <-ha.ticker.C:
			ha.Update()
		}
	}
}

func NewHActivity(session *discordgo.Session) *HentaiActivity {
	return &HentaiActivity{
		Client:  nhapi.New(nhapi.Options{Url: nhapi.NHentaiApiUrl, Timeout: time.Second * 30}),
		Session: session,
		Logger:  logrus.WithField("component", "HentaiActivity"),
	}
}
