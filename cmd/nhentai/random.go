package nhentai

import (
	"github.com/Ressetkk/Iku-chan/internal/embed"
	"github.com/Ressetkk/Iku-chan/pkg/dux"
	"github.com/Ressetkk/Iku-chan/pkg/nhapi"
	"time"
)

func RandomCmd() *dux.Command {
	cmd := &dux.Command{
		Name: "random",
		Run:  randomHandler,
		Description: `Come and grab the sauce from the sacred pool bestowed upon us from lewders all across the globe!
Disclaimer: You will literally get random sauce. We're not even generating it! Use at your own risk!`,
		Short:   "Get random sauce.",
		Example: "random",
	}

	cmd.AddMiddleware(dux.NSFWOnly)
	return cmd
}

func randomHandler(ctx *dux.Context, args []string) {
	client := nhapi.New(nhapi.Options{
		Url:     nhapi.NHentaiApiUrl,
		Timeout: time.Second * 30,
	})
	id, err := client.Random()
	if err != nil {
		ctx.SendTextf("NHentai API returned an error: %v", err)
		return
	}
	res, err := client.Get(id)
	if err != nil {
		_, e := ctx.SendTextf("NHentai API returned an error: %v", err)
		ctx.Logger.WithError(e).Debug("context returned an error")
		return
	}
	toSend := embed.Make(res)
	_, err = ctx.SendEmbed(&toSend)
	if err != nil {
		ctx.Logger.WithError(err).Error("failed send embed")
	}
}
