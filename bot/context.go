package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io"
)

// Context: Bot context
type Context struct {
	Discord      *discordgo.Session
	Guild        *discordgo.Guild
	VoiceChannel *discordgo.Channel
	TextChannel  *discordgo.Channel
	User         *discordgo.User
	Message      *discordgo.MessageCreate
	Args         []string

	Conf       *Config
	CmdHandler *CommandHandler
	Sessions   *SessionManager
}

// Create new context
func NewContext(discord *discordgo.Session, guild *discordgo.Guild, textChannel *discordgo.Channel,
	user *discordgo.User, message *discordgo.MessageCreate, conf *Config, cmdHandler *CommandHandler,
	sessions *SessionManager) *Context {
	ctx := new(Context)
	ctx.Discord = discord
	ctx.Guild = guild
	ctx.TextChannel = textChannel
	ctx.User = user
	ctx.Message = message
	ctx.Conf = conf
	ctx.CmdHandler = cmdHandler
	ctx.Sessions = sessions
	return ctx
}

// Reply on massege
func (ctx Context) Reply(content string) *discordgo.Message {
	msg, err := ctx.Discord.ChannelMessageSend(ctx.TextChannel.ID, content)
	if err != nil {
		fmt.Println("Error whilst sending message,", err)
		return nil
	}
	return msg
}

// Reply on massege with file
func (ctx Context) ReplyFile(name string, r io.Reader) *discordgo.Message {
	msg, err := ctx.Discord.ChannelFileSend(ctx.TextChannel.ID, name, r)
	if err != nil {
		fmt.Println("Error whilst sending file,", err)
		return nil
	}
	return msg
}

// Returns user voice channel
func (ctx *Context) GetVoiceChannel() *discordgo.Channel {
	if ctx.VoiceChannel != nil {
		return ctx.VoiceChannel
	}
	for _, state := range ctx.Guild.VoiceStates {
		if state.UserID == ctx.User.ID {
			channel, _ := ctx.Discord.State.Channel(state.ChannelID)
			ctx.VoiceChannel = channel
			return channel
		}
	}
	return nil
}
