package nhentai

import (
	"fmt"
	"github.com/Ressetkk/Iku-chan/pkg/dux"
	"github.com/Ressetkk/Iku-chan/pkg/nhapi"
	"github.com/bwmarrin/discordgo"
	"math/rand"
	"strconv"
	"time"
)

func GetCmd() *dux.Command {
	return &dux.Command{
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
					e := ctx.SendTextf("NHentai API returned an error: %s", err)
					ctx.Logger.WithError(e).Error("context returned an error")
					return
				}
				toSend := makeEmbed(res)
				err = ctx.SendEmbed(&toSend)
				if err != nil {
					ctx.Logger.WithError(err).Error("failed send embed")
				}
			}
		},
	}
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

func makeEmbed(result *nhapi.Result) discordgo.MessageEmbed {
	embed := discordgo.MessageEmbed{
		URL:         fmt.Sprintf("%s/g/%v", nhapi.NHentaiApiUrl, result.ID),
		Title:       result.Title["english"],
		Description: result.Title["japanese"],
		Timestamp:   time.Unix(result.UploadTimestamp, 0).Format(time.RFC3339),
		Color:       randomRGB(),
		Footer:      &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("No. %v", result.ID)},
		Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: result.GetThumbnailUrl()},
	}
	pages := discordgo.MessageEmbedField{
		Name:   "Pages",
		Inline: true,
		Value:  fmt.Sprintf("%v", result.NumOfPages),
	}
	languages := discordgo.MessageEmbedField{
		Name:   "Languages",
		Inline: true,
	}
	categories := discordgo.MessageEmbedField{
		Name:   "Categories",
		Inline: true,
	}
	tags := discordgo.MessageEmbedField{
		Name: "Tags",
	}
	for _, tag := range result.Tags {
		switch tag.Type {
		case "artist":
			if embed.Author == nil {
				embed.Author = &discordgo.MessageEmbedAuthor{
					URL:  tag.URL.Full(),
					Name: tag.Name,
				}
			}
			break
		case "tag":
			if len(tags.Value+tag.Name) < 1020 {
				if tags.Value != "" {
					tags.Value += ", "
				}
				tags.Value += tag.Name
			}
			break
		case "category":
			if len(categories.Value+tag.Name) < 1020 {
				if categories.Value != "" {
					categories.Value += ", "
				}
				categories.Value += fmt.Sprintf("[%s](%s)", tag.Name, tag.URL.Full())
			}
			break
		case "language":
			if len(languages.Value+tag.Name) < 1020 {
				if languages.Value != "" {
					languages.Value += ", "
				}
				languages.Value += fmt.Sprintf("[%s](%s)", tag.Name, tag.URL.Full())
			}
			break
		}
	}
	embed.Fields = []*discordgo.MessageEmbedField{&pages, &languages, &categories, &tags}
	return embed
}

func randomRGB() int {
	var rgb int
	r := rand.Intn(255)
	g := rand.Intn(255)
	b := rand.Intn(255)
	rgb = (r & 0xFF << 16) | (g & 0xFF << 8) | (b & 0xFF)
	return rgb
}
