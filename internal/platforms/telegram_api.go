package platforms

import (
	"log/slog"
	"time"

	"github.com/napuu/gpsp-bot/internal/chain"
	"github.com/napuu/gpsp-bot/internal/config"
	"github.com/napuu/gpsp-bot/internal/handlers"

	tele "gopkg.in/telebot.v4"
)

func wrapTeleHandler(bot *tele.Bot, chain *chain.HandlerChain) func(c tele.Context) error {
	return func(c tele.Context) error {
		chain.Process(&handlers.Context{TelebotContext: c, Telebot: bot, Service: handlers.Telegram})
		return nil
	}
}

func TelebotCompatibleVisibleCommands() []tele.Command {
	commands := make([]tele.Command, 0, len(config.EnabledFeatures()))
	for _, action := range config.EnabledFeatures() {
		if handlers.Action(action) == handlers.Ping {
			continue
		}
		commands = append(commands, tele.Command{
			Text:        string(action),
			Description: string(handlers.ActionMap[handlers.Action(action)]),
		})
	}
	return commands
}

func RunTelegramBot() {
	bot := getTelegramBot()
	chain := chain.NewChainOfResponsibility()

	err := bot.SetCommands(TelebotCompatibleVisibleCommands())
	if err != nil {
		slog.Error(err.Error())
	}

	bot.Handle(tele.OnText, wrapTeleHandler(bot, chain))
	
	// Handle reactions to track emoji responses to videos
	bot.Handle(tele.OnReaction, wrapReactionHandler(bot, chain))

	go bot.Start()
}

func wrapReactionHandler(bot *tele.Bot, chain *chain.HandlerChain) func(c tele.Context) error {
	return func(c tele.Context) error {
		// Create context for reaction events
		ctx := &handlers.Context{
			TelebotContext: c,
			Telebot:        bot,
			Service:        handlers.Telegram,
			reaction:       true,
		}
		
		// Extract reaction information
		if c.Message() != nil && c.Message().Reactions != nil {
			// Get the latest reaction
			reactions := c.Message().Reactions
			if len(reactions) > 0 {
				latestReaction := reactions[len(reactions)-1]
				ctx.reactionMessageID = c.Message().Text // Use message text as identifier
				ctx.reactionUser = latestReaction.User.Username
				ctx.reactionEmoji = latestReaction.Emoji
			}
		}
		
		chain.Process(ctx)
		return nil
	}
}

func getTelegramBot() *tele.Bot {
	pref := tele.Settings{
		Token:     config.FromEnv().TELEGRAM_TOKEN,
		ParseMode: tele.ModeHTML,
		Poller: &tele.LongPoller{
			Timeout: 10 * time.Second,
			AllowedUpdates: []string{
				"message",
				"message_reaction",
			},
		},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		panic(err)
	}

	return b
}
