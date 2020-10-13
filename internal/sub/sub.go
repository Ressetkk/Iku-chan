package sub

import (
	"github.com/Ressetkk/Iku-chan/internal/embed"
	"github.com/Ressetkk/Iku-chan/pkg/dux"
	"github.com/Ressetkk/Iku-chan/pkg/nhapi"
	"github.com/Ressetkk/dgwidgets"
)

type PagUpdater struct {
	Pag                   *dgwidgets.Paginator
	Ctx                   *dux.Context
	Query, Sort           string
	Client                *nhapi.Client
	MaxIndex, Pages, Page int
}

func (s *PagUpdater) OnNotify(index int) {
	if s.MaxIndex == index && s.Page < s.Pages {
		s.Page++

		res, err := s.Client.Search(s.Query, s.Sort, s.Page)
		if err != nil {
			s.Ctx.SendTextf("NHentai API returned an error: %v", err)
			return
		}
		if len(res.Results) == 0 {
			s.Ctx.SendText("NHentai returned no results for this query...")
			return
		}
		s.MaxIndex += res.NumPerPage
		s.Pag.Add(embed.GenerateEmbeds(res.Results)...)
	}
}
