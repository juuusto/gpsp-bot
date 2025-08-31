package handlers

import (
	"fmt"
	"log/slog"
	"sort"
	"strings"

	"github.com/napuu/gpsp-bot/internal/db"
)

type ReactionsHandler struct {
	next ContextHandler
}

func (r *ReactionsHandler) Execute(m *Context) {
	slog.Debug("Entering ReactionsHandler")
	
	if m.action == Reactions {
		// Get reaction statistics from database
		dbHandle := db.GetGlobalDB()
		if dbHandle == nil {
			m.textResponse = "Tietokanta ei ole saatavilla."
			r.next.Execute(m)
			return
		}

		stats, err := db.GetReactionStats(dbHandle)
		if err != nil {
			slog.Error("Failed to get reaction stats", "error", err)
			m.textResponse = "Virhe haettaessa reaktioita."
			r.next.Execute(m)
			return
		}

		if len(stats) == 0 {
			m.textResponse = "Ei vielä emoji-reaktioita tallennettu."
			r.next.Execute(m)
			return
		}

		// Sort emojis by count (descending)
		type emojiCount struct {
			emoji string
			count int
		}
		var sortedStats []emojiCount
		for emoji, count := range stats {
			sortedStats = append(sortedStats, emojiCount{emoji, count})
		}
		sort.Slice(sortedStats, func(i, j int) bool {
			return sortedStats[i].count > sortedStats[j].count
		})

		// Build response message
		var response strings.Builder
		response.WriteString("📊 Emoji-reaktioiden tilastot:\n\n")
		
		for _, stat := range sortedStats {
			response.WriteString(fmt.Sprintf("%s: %d kertaa\n", stat.emoji, stat.count))
		}

		m.textResponse = response.String()
	}

	r.next.Execute(m)
}

func (r *ReactionsHandler) SetNext(next ContextHandler) {
	r.next = next
}
