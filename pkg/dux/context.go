package dux

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"sync"
)

type Context struct {
	Session *discordgo.Session
	Event   *discordgo.MessageCreate
	User    *discordgo.User
	Channel *discordgo.Channel
	Message *discordgo.Message
	Route   *Command

	Logger *logrus.Entry
	sync.Mutex
	keys map[string]interface{}
}

func (c *Context) Set(key string, val interface{}) {
	c.Lock()
	c.keys[key] = val
	c.Unlock()
}

func (c *Context) Get(key string) (interface{}, error) {
	if val, ok := c.keys[key]; ok {
		return val, nil
	}
	return nil, fmt.Errorf("key not found: %s", key)
}

// SendText sends text message to the Discord API.
func (c *Context) SendText(text string) (err error) {
	_, err = c.Session.ChannelMessageSend(c.Event.ChannelID, text)
	return
}

// SendTextf sends text message to the Discord API with custom formatting.
func (c *Context) SendTextf(format string, elements ...interface{}) (err error) {
	_, err = c.Session.ChannelMessageSend(c.Event.ChannelID, fmt.Sprintf(format, elements...))
	return
}

// SendEmbed sends discordgo.MessageEmbed struct to the Discord API.
func (c *Context) SendEmbed(embed *discordgo.MessageEmbed) (err error) {
	_, err = c.Session.ChannelMessageSendEmbed(c.Event.ChannelID, embed)
	return
}

// Send sends complete discordgo.MessageSend struct to the Discord API.
func (c *Context) Send(msg *discordgo.MessageSend) (err error) {
	_, err = c.Session.ChannelMessageSendComplex(c.Event.ChannelID, msg)
	return
}

func NewContext(s *discordgo.Session, m *discordgo.MessageCreate) (*Context, error) {
	ch, err := s.Channel(m.ChannelID)
	if err != nil {
		return nil, err
	}
	return &Context{
		Session: s,
		Event:   m,
		User:    m.Author,
		Channel: ch,
		Message: m.Message,
		keys:    make(map[string]interface{}),
		Mutex:   sync.Mutex{},
	}, nil
}
