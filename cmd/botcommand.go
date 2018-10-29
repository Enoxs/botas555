package cmd

import (
	"strconv"

	"github.com/FlameInTheDark/dtbot/bot"
)

// BotCommand special bot commands handler
func BotCommand(ctx bot.Context) {
	if len(ctx.Args) == 0 {
		return
	}
	switch ctx.Args[0] {
	case "clear":
		if len(ctx.Args) < 2 {
			ctx.BotMsg.Clear(&ctx, 0)
			return
		}
		from, err := strconv.Atoi(ctx.Args[1])
		if err != nil {
			return
		}
		ctx.BotMsg.Clear(&ctx, from)
	case "help":
		ctx.ReplyEmbed(ctx.Loc("help"), ctx.Loc("help_reply"))
	}
}
