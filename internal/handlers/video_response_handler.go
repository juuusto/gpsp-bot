package handlers

import (
	"bytes"
	"log/slog"
	"os"
	"time"
	"crypto/sha256"
	"encoding/hex"

	"github.com/bwmarrin/discordgo"
	"github.com/napuu/gpsp-bot/pkg/utils"
	"github.com/napuu/gpsp-bot/internal/db"
	tele "gopkg.in/telebot.v4"
)

type VideoResponseHandler struct {
	next ContextHandler
}

func (r *VideoResponseHandler) Execute(m *Context) {
	slog.Debug("Entering VideoResponseHandler")

	if len(m.finalVideoPath) > 0 {
		switch m.Service {
		case Telegram:
			chatId := tele.ChatID(utils.S2I(m.chatId))

			m.Telebot.Send(chatId, &tele.Video{File: tele.FromDisk(m.finalVideoPath)})
			m.sendVideoSucceeded = true

			// Record sent video
			go recordSentVideo(m, "telegram")
		case Discord:
			file, err := os.Open(m.finalVideoPath)
			if err != nil {
				panic(err)
			}
			defer file.Close()

			buf := bytes.NewBuffer(nil)
			_, err = buf.ReadFrom(file)
			if err != nil {
				panic(err)
			}

			message := &discordgo.MessageSend{
				Content: "",
				Files: []*discordgo.File{
					{
						Name:        "video.mp4", // this apparently doesn't matter
						ContentType: "video/mp4",
						Reader:      buf,
					},
				},
			}

			_, err = m.DiscordSession.ChannelMessageSendComplex(m.chatId, message)
			if err != nil {
				slog.Debug(err.Error())
			} else {
				m.sendVideoSucceeded = true

				// Record sent video
				go recordSentVideo(m, "discord")
			}
		}
	}

	r.next.Execute(m)
}

func (u *VideoResponseHandler) SetNext(next ContextHandler) {
	u.next = next
}

// recordSentVideo saves info about the sent video to the DB
func recordSentVideo(m *Context, service string) {
	dbHandle := db.GetGlobalDB()
	if dbHandle == nil {
		return
	}
	// Compute hash of the video file
	f, err := os.Open(m.finalVideoPath)
	if err != nil {
		return
	}
	defer f.Close()
	h := sha256.New()
	if _, err := h.Write([]byte(m.finalVideoPath)); err == nil {
		// fallback: just hash path if can't read
		// but try to hash file content
		if _, err := f.Seek(0, 0); err == nil {
			buf := make([]byte, 1024*1024)
			for {
				n, err := f.Read(buf)
				if n > 0 {
					h.Write(buf[:n])
				}
				if err != nil {
					break
				}
			}
		}
	}
	videoHash := hex.EncodeToString(h.Sum(nil))

	user := m.chatId
	if service == "telegram" && m.TelebotContext != nil && m.TelebotContext.Sender() != nil {
		user = m.TelebotContext.Sender().Username
	}
	if service == "discord" && m.DiscordMessage != nil && m.DiscordMessage.Author != nil {
		user = m.DiscordMessage.Author.Username
	}

	db.InsertSentVideo(dbHandle, db.SentVideo{
		User:        user,
		VideoHash:   videoHash,
		MessageUUID: m.id,
		Service:     service,
		CreatedAt:   time.Now().Unix(),
	})
}
