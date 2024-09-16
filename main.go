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
			return
		case "!stop":
			if voiceConnection != nil {
				voiceConnection.Disconnect()
			}
		}

		srch := ""
		if len(m.Content) > 6 {
			srch = m.Content[0:6]
		}

		if srch == "!play " {
			searchString := m.Content[5:]

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

			link, err := SearchYoutube(searchString, m)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Failed to find request.")
				return
			}

			s.ChannelMessageSend(m.ChannelID, link.VideoURL)

			err = playAudio(vc, link.VideoURL)
			if err != nil {
				log.Println("Error playing audio:", err)
				s.ChannelMessageSend(m.ChannelID, "Error playing audio.")
			}

			return

			// Remember to disconnect from the voice channel when done
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
