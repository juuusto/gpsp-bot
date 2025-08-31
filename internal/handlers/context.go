package handlers

import (
	"time"

	dg "github.com/bwmarrin/discordgo"
	"github.com/napuu/gpsp-bot/pkg/utils"
	tele "gopkg.in/telebot.v4"
)

type Service int

const (
	Telegram Service = iota + 1
	Discord
	// Matrix // not supported, perhaps at one point
)

type ContextHandler interface {
	Execute(*Context)
	SetNext(ContextHandler)
}

type Action string
type ActionDescription string

const (
	Tuplilla      Action = "tuplilla"
	DownloadVideo Action = "dl"
	Ping          Action = "ping"
	Euribor       Action = "euribor"
	Reactions     Action = "reactions"
)

var ActionMap = map[Action]ActionDescription{
	Tuplilla:      "Tuplilla...",
	DownloadVideo: "Lataa video",
	Euribor:       "Tuoreet Euribor-korot",
	Ping:          "Ping",
	Reactions:     "Näytä emoji-reaktioiden tilastot",
}

type Context struct {
	Service Service
	// The original message without any parsing
	// (except on Telegram events, the possible "@<botname>"" is removed)
	rawText string
	// Message without action string and
	// possibly related prefixes or suffixes
	parsedText string
	id         string // Some services use string, some int, some int64. They're now strings at our context.
	replyToId  string
	// Must store separate from replyToId as
	// replyToId = 0 might refer to first message
	// or to no message at all
	shouldReplyToMessage bool
	isReply              bool
	chatId               string
	action               Action
	url                  string

	doneTyping         chan struct{}
	gotDubz            bool
	dubzNegation       chan string
	lastCubeThrownTime time.Time

	rates utils.EuriborRateEntry

	TelebotContext tele.Context
	Telebot        *tele.Bot

	DiscordSession *dg.Session
	DiscordMessage *dg.MessageCreate

	originalVideoPath string

	// Location of the video that is finally sent.
	// Different handlers might edit this during the processing.
	finalVideoPath                string
	finalImagePath                string
	textResponse                  string
	sendVideoSucceeded            bool
	startSeconds                  chan float64
	durationSeconds               chan float64
	cutVideoArgsParsed            chan bool
	shouldDeleteOriginalMessage   bool
	shouldNagAboutOriginalMessage bool

	// Reaction tracking fields
	reaction          bool   // true if this context represents a reaction event
	reactionMessageID string // ID of the message that was reacted to
	reactionUser      string // user who reacted
	reactionEmoji     string // emoji used in the reaction
}
