package embed

import (
	"fmt"
	"github.com/Ressetkk/Iku-chan/pkg/nhapi"
	"github.com/bwmarrin/discordgo"
	"math/rand"
	"time"
)

func Make(result *nhapi.Result) discordgo.MessageEmbed {
	embed := discordgo.MessageEmbed{
		URL:         fmt.Sprintf("%s/g/%v", nhapi.NHentaiApiUrl, result.ID),
		Title:       result.Title["english"],
		Description: result.Title["japanese"],
		Timestamp:   time.Unix(result.UploadTimestamp, 0).Format(time.RFC3339),
		Color:       randomRGB(),
		Footer:      &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("No. %v", result.ID)},
		Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: result.GetCoverThumbnail()},
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

func GenerateEmbeds(results []nhapi.Result) []*discordgo.MessageEmbed {
	var embeds []*discordgo.MessageEmbed
	for _, r := range results {
		e := Make(&r)
		embeds = append(embeds, &e)
	}
	return embeds
}
