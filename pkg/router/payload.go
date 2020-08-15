package router

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

// TODO maybe Interface???
type Payload struct {
	Session *discordgo.Session
	Event   *discordgo.MessageCreate
}

// SendText sends text message to the Discord API.
func (h Payload) SendText(text string) (err error) {
	_, err = h.Session.ChannelMessageSend(h.Event.ChannelID, text)
	return
}

// SendTextf sends text message to the Discord API with custom formatting.
func (h Payload) SendTextf(format string, elements ...interface{}) (err error) {
	_, err = h.Session.ChannelMessageSend(h.Event.ChannelID, fmt.Sprintf(format, elements...))
	return
}

// SendEmbed sends discordgo.MessageEmbed struct to the Discord API.
func (h Payload) SendEmbed(embed *discordgo.MessageEmbed) (err error) {
	_, err = h.Session.ChannelMessageSendEmbed(h.Event.ChannelID, embed)
	return
}

// Send sends complete discordgo.MessageSend struct to the Discord API.
func (h Payload) Send(msg *discordgo.MessageSend) (err error) {
	_, err = h.Session.ChannelMessageSendComplex(h.Event.ChannelID, msg)
	return
}
