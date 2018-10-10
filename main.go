package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"./bot"
	"./cmd"
	"github.com/bwmarrin/discordgo"
)

var (
	conf *bot.Config
	// CmdHandler bot command handler
	CmdHandler *bot.CommandHandler
	// Sessions bot session manager
	Sessions *bot.SessionManager
	botId    string
)

func main() {
	conf = bot.LoadConfig()
	CmdHandler = bot.NewCommandHandler()
	registerCommands()
	Sessions = bot.NewSessionManager()
	discord, err := discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	if err != nil {
		fmt.Println("Create session error, ", err)
		return
	}
	usr, err := discord.User("@me")
	if err != nil {
		fmt.Println("Error obtaining account details,", err)
		return
	}
	botId = usr.ID
	discord.AddHandler(commandHandler)
	discord.AddHandler(func(discord *discordgo.Session, ready *discordgo.Ready) {
		discord.UpdateStatus(0, "ящик")
		guilds := discord.State.Guilds
		fmt.Println("Ready with", len(guilds), "guilds.")
	})
	err = discord.Open()
	if err != nil {
		fmt.Printf("Connection open error: %v", err)
		return
	}
	defer discord.Close()
	fmt.Println("Bot is now running.")

	sc := make(chan os.Signal, 1)

	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

// Handle discord messages
func commandHandler(discord *discordgo.Session, message *discordgo.MessageCreate) {
	user := message.Author
	if user.ID == botId || user.Bot {
		return
	}
	args := strings.Split(message.Content, " ")
	name := strings.ToLower(args[0])
	command, found := CmdHandler.Get(name)
	if !found {
		return
	}
	channel, err := discord.State.Channel(message.ChannelID)
	if err != nil {
		fmt.Println("Error getting channel,", err)
		return
	}
	guild, err := discord.State.Guild(channel.GuildID)
	if err != nil {
		fmt.Println("Error getting guild,", err)
		return
	}
	ctx := bot.NewContext(discord, guild, channel, user, message, conf, CmdHandler, Sessions)
	ctx.Args = args[1:]
	c := *command
	c(*ctx)
}

// Adds bot commands
func registerCommands() {
	CmdHandler.Register("!r", cmd.PlayerCommand)
	CmdHandler.Register("!w", cmd.WeatherCommand)
	CmdHandler.Register("!t", cmd.TranslateCommand)
	CmdHandler.Register("!n", cmd.NewsCommand)
	CmdHandler.Register("!c", cmd.CurrencyCommand)
}
