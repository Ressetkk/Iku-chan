package nhentai

import (
	"github.com/Ressetkk/Iku-chan/pkg/dux"
)

func GetCmd() *dux.Command {
	return &dux.Command{
		Name:  "get",
		Short: "Use sacred numbers to get the sauce",
		Description: `Use sacred number to get the most spicy sauce out there.
Don't be shy, just type the six-digit number and have some fun!`,
		Example: "get 177013",
		Run: func(ctx *dux.Context, args []string) {
			ctx.SendTextf("Not yet implemented... args %s", args)
		},
	}
}
