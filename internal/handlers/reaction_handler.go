package handlers

import (
	"log/slog"
	"time"

	"github.com/napuu/gpsp-bot/internal/db"
)

type ReactionHandler struct {
	next ContextHandler
}

func (r *ReactionHandler) Execute(m *Context) {
	slog.Debug("Entering ReactionHandler")
	
	// This handler processes reactions to video messages
	// It will be called when reaction events are received
	if m.reaction != nil {
		// Record the reaction in the database
		go recordVideoReaction(m)
	}

	r.next.Execute(m)
}

func (r *ReactionHandler) SetNext(next ContextHandler) {
	r.next = next
}

// recordVideoReaction saves the reaction to the database
func recordVideoReaction(m *Context) {
	dbHandle := db.GetGlobalDB()
	if dbHandle == nil {
		return
	}

	reaction := db.VideoReaction{
		MessageUUID: m.reactionMessageID,
		User:        m.reactionUser,
		Emoji:       m.reactionEmoji,
		CreatedAt:   time.Now().Unix(),
	}

	err := db.InsertVideoReaction(dbHandle, reaction)
	if err != nil {
		slog.Error("Failed to record video reaction", "error", err)
	} else {
		slog.Info("Recorded video reaction", 
			"message_id", m.reactionMessageID,
			"user", m.reactionUser,
			"emoji", m.reactionEmoji)
	}
}
