package bot

import (
	"github.com/bwmarrin/discordgo"
)

// Context : Bot context structure
type Context struct {
	Discord      *discordgo.Session
	Guild        *discordgo.Guild
	VoiceChannel *discordgo.Channel
	TextChannel  *discordgo.Channel
	User         *discordgo.User
	Message      *discordgo.MessageCreate
	Args         []string

	DB         *DBWorker
	Conf       *Config
	CmdHandler *CommandHandler
	Sessions   *SessionManager
	Youtube    *Youtube
	BotMsg     *BotMessages
	Data       *DataType
	Guilds     GuildsMap
}

// NewContext create new context
func NewContext(discord *discordgo.Session, guild *discordgo.Guild, textChannel *discordgo.Channel,
	user *discordgo.User, message *discordgo.MessageCreate, conf *Config, cmdHandler *CommandHandler,
	sessions *SessionManager, youtube *Youtube, botMsg *BotMessages, dataType *DataType, dbWorker *DBWorker, guilds GuildsMap) *Context {
	ctx := new(Context)
	ctx.Discord = discord
	ctx.Guild = guild
	ctx.TextChannel = textChannel
	ctx.User = user
	ctx.Message = message
	ctx.Conf = conf
	ctx.CmdHandler = cmdHandler
	ctx.Sessions = sessions
	ctx.Youtube = youtube
	ctx.BotMsg = botMsg
	ctx.Data = dataType
	ctx.DB = dbWorker
	ctx.Guilds = guilds
	return ctx
}

// Loc returns translated string by key
func (ctx *Context) Loc(key string) string {
	// Check if translation exist
	if len(ctx.Conf.Locales[ctx.GetGuild().Language][key]) == 0 {
		return ctx.Conf.Locales["en"][key]
	}
	return ctx.Conf.Locales[ctx.GetGuild().Language][key]
}

// WeatherCode returns unicode symbol of weather font icon
func (ctx *Context) WeatherCode(code string) string {
	return ctx.Conf.WeatherCodes[code]
}

// GetVoiceChannel returns user voice channel
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

func (ctx *Context) GetGuild() *GuildData {
	return ctx.Guilds[ctx.Guild.ID]
}
