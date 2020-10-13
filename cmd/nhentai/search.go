package nhentai

import (
	"github.com/Ressetkk/Iku-chan/internal/embed"
	"github.com/Ressetkk/Iku-chan/internal/sub"
	"github.com/Ressetkk/Iku-chan/pkg/dux"
	"github.com/Ressetkk/Iku-chan/pkg/nhapi"
	"github.com/Ressetkk/dgwidgets"
	"strings"
	"time"
)

func SearchCmd() *dux.Command {
	cmd := &dux.Command{
		Name: "search",
		Description: `Search NHentai for the sauce.
Type whatever you like, and I will show you the results from the site.
The command accepts one parameter for search:
sort - sort results. [date, popular]`,
		Short:   "Search NHentai for the sauce",
		Example: "search jojo sort=date",
		Run:     searchHandler,
	}

	cmd.AddMiddleware(dux.NSFWOnly)
	return cmd
}
func searchHandler(ctx *dux.Context, args []string) {
	client := nhapi.New(nhapi.Options{Url: nhapi.NHentaiApiUrl, Timeout: time.Second * 30})
	query, sort := parseQuery(args)
	if query == "" {
		ctx.SendText("Provide search query.")
		return
	}

	pag := dgwidgets.NewPaginator(ctx.Session, ctx.Message.ChannelID)
	setupPaginator(pag)
	pag.Widget.LockToUsers(ctx.Message.Author.ID)

	results, err := client.Search(query, sort, 1)
	if err != nil {
		ctx.SendTextf("NHentai API returned an error: %v", err)
		return
	}
	if len(results.Results) == 0 {
		ctx.SendText("NHentai returned no results for this query...")
		return
	}

	pag.Add(embed.GenerateEmbeds(results.Results)...) // add initial elements
	pag.Index.AddSub(&sub.PagUpdater{
		Pag:      pag,
		Client:   client,
		Ctx:      ctx,
		MaxIndex: results.NumPerPage - 1,
		Page:     1, Pages: results.NumOfPages,
		Query: query, Sort: sort})

	err = pag.Spawn()
	if err != nil {
		ctx.Logger.WithError(err).Error("paginator handling failed")
	}
}

func setupPaginator(pag *dgwidgets.Paginator) {
	pag.ColourWhenDone = 0x32e0c4
	pag.Loop = true
	pag.DeleteReactionsWhenDone = true
	pag.Widget.RefreshAfterAction = true
	pag.Widget.Timeout = time.Second * 30
}

func parseQuery(args []string) (query string, sort string) {
	for _, arg := range args {
		switch {
		case strings.HasPrefix(arg, "sort="):
			sort = strings.ReplaceAll(arg, "sort=", "")
		default:
			query += arg + " "
		}
	}
	return
}
