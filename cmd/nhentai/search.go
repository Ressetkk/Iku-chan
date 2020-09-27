package nhentai

import "github.com/Ressetkk/Iku-chan/pkg/dux"

func SearchCmd() *dux.Command {
	return &dux.Command{
		Name: "search",
		Description: `Search NHentai for the sauce.
Type whatever you like, and I will show you the results from the site.
The command accepts two parameters for search:
sort - TBA
TBA - TBA`,
		Short:   "Search NHentai for the sauce",
		Example: "search jojo sort=date",
		Run: func(ctx *dux.Context, args []string) {
			ctx.SendTextf("Not yet implemented... args %s", args)
		},
	}
}
