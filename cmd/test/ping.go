package test

import (
	"fmt"
	"github.com/Ressetkk/Iku-chan/pkg/router"
	"github.com/sparrc/go-ping"
)

func PingCmd() *router.Command {
	cmd := &router.Command{
		Use:         "ping",
		Description: "Check bot's ping to the Discord API. Bot will answer with the response time in milliseconds.",
		Short:       "Respond with \"pong!\" And give response time.",
		Example:     "ping",
		Run: func(h router.Payload) {
			//FIXME resolving the DNS records and set URL to Discord API.
			p, err := ping.NewPinger("8.8.8.8")
			if err != nil {
				h.SendText(err.Error())
				return
			}
			p.Count = 1
			p.SetPrivileged(true)
			p.OnFinish = func(s *ping.Statistics) {
				h.SendText(fmt.Sprintf("Pong! Avg. Response time **%s**.", s.AvgRtt))
			}
			p.Run()
		},
	}
	return cmd
}
