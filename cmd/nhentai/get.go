package nhentai

import (
	"github.com/Ressetkk/Iku-chan/internal/embed"
	"github.com/Ressetkk/Iku-chan/pkg/dux"
	"github.com/Ressetkk/Iku-chan/pkg/nhapi"
	"strconv"
	"time"
)

func GetCmd() *dux.Command {
	cmd := &dux.Command{
		Name:  "get",
		Short: "Use sacred numbers to get the sauce",
		Description: `Use sacred number to get the most spicy sauce out there.
Don't be shy, just type the six-digit number and have some fun!

You can also get *multiple* sauces! Just add more numbers to the query.`,
		Example: "get 177013",
		Run: func(ctx *dux.Context, args []string) {
			client := nhapi.New(nhapi.Options{
				Url:     nhapi.NHentaiApiUrl,
				Timeout: time.Second * 30,
			})
			if len(args) == 0 {
				ctx.SendText("Please provide some numbers!")
				return
			}
			for _, num := range parseNumbers(args) {
				res, err := client.Get(num)
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
		},
	}

	cmd.AddMiddleware(dux.NSFWOnly)
	return cmd
}

func parseNumbers(args []string) []int {
	var numbers []int
	for _, arg := range args {
		iarg, err := strconv.Atoi(arg)
		if err != nil {
			continue
		}
		numbers = append(numbers, iarg)
	}
	return numbers
}
