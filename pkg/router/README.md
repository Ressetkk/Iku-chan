# router 
This packacge attempts to provide [Cobra-like](https://github.com/spf13/cobra/) commands hierarchy for Discord bots.

To create a new DiscordGo binding follow steps below:

1. Create Command struct and define needed values:
2. Each command can be used as a handler for DiscordGo Events. Command struct implements `Handler` function which returns
handler function for DiscordGo.
3. Register new Handler for DiscordGo session.
4. Add commands and middleware.

Example:
```go
func main() {
    session, _ := discordgo.New("Bot " + token)
    cmd := router.Command{
        Use:     "helloworld",
        Aliases: []string{"hello", "world"},
        Description: "prints hello world to the world!",
        Short:   "print hello world",
        Example: "helloworld",
        Run: func(h *router.Payload) {
            h.SendText("Hello World!")
        },
    }
    session.AddHandler(cmd.Handler(router.Options{AllowMentions: true, IgnoreCases: true}))
    session.Open()
}
```

More detailed documentation soon.