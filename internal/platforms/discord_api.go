package platforms

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/napuu/gpsp-bot/internal/chain"
	"github.com/napuu/gpsp-bot/internal/config"
	"github.com/napuu/gpsp-bot/internal/handlers"
)

func wrapDiscoHandler(chain *chain.HandlerChain) func(s *discordgo.Session, m *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// Ignore all messages created by the bot itself
		if m.Author.ID == s.State.User.ID {
			return
		}

		// Wrap the context
		chain.Process(&handlers.Context{
			DiscordSession: s,
			DiscordMessage: m,
			Service:        handlers.Discord,
		})
	}
}

func RunDiscordBot() {
	token := config.FromEnv().DISCORD_TOKEN
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		slog.Error("Error creating Discord session", "error", err)
		return
	}

	// Create the chain of responsibility
	chain := chain.NewChainOfResponsibility()

	// Add a handler for messages
	dg.AddHandler(wrapDiscoHandler(chain))
	
	// Add a handler for reactions
	dg.AddHandler(wrapDiscordReactionHandler(chain))

	// Specify intents
	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages

	// Open the connection
	err = dg.Open()
	if err != nil {
		slog.Error("Error opening Discord connection", "error", err)
		return
	}
}

func wrapDiscordReactionHandler(chain *chain.HandlerChain) func(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	return func(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
		// Create context for reaction events
		ctx := &handlers.Context{
			DiscordSession: s,
			Service:        handlers.Discord,
			reaction:       true,
			reactionMessageID: r.MessageID,
			reactionUser:      r.User.Username,
			reactionEmoji:     r.Emoji.Name,
		}
		
		chain.Process(ctx)
	}
}
