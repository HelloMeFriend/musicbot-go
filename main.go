package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var voiceConnection *discordgo.VoiceConnection

func main() {

	sess, err := discordgo.New("Bot ")
	if err != nil {
		log.Fatal(err)
		return
	}

	sess.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {

		if m.Author.ID == s.State.User.ID {
			return
		}

		switch m.Content {
		case "hello":
			s.ChannelMessageSend(m.ChannelID, "world!")

		case "!play":
			// Find the voice channel the user is in
			guild, err := s.State.Guild(m.GuildID)
			if err != nil {
				log.Println("Error finding guild:", err)
				return
			}

			var userVoiceState *discordgo.VoiceState
			for _, vs := range guild.VoiceStates {
				if vs.UserID == m.Author.ID {
					userVoiceState = vs
					break
				}
			}

			if userVoiceState == nil {
				s.ChannelMessageSend(m.ChannelID, "You must be in a voice channel first!")
				return
			}

			// Join the user's voice channel
			vc, err := s.ChannelVoiceJoin(m.GuildID, userVoiceState.ChannelID, false, true)
			if err != nil {
				log.Println("Failed to join voice channel:", err)
				s.ChannelMessageSend(m.ChannelID, "Failed to join voice channel.")
				return
			}

			voiceConnection = vc
			s.ChannelMessageSend(m.ChannelID, "Joined the voice channel!")

			link, err := SearchYoutube("3hunna", m)
			s.ChannelMessageSend(m.ChannelID, link.VideoURL)

			// Remember to disconnect from the voice channel when done

		case "!quit":
			if voiceConnection != nil {
				voiceConnection.Disconnect()
			}
		}
	})

	sess.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	err = sess.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer sess.Close()

	fmt.Println("Bot online")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
