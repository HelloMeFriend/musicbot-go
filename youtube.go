package main

import (
	"context"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type Song struct {
	ChannelID string
	User      string
	VidID     string
	Title     string
	Duration  string
	VideoURL  string
}

func SearchYoutube(searchString string, m *discordgo.MessageCreate) (song_struct Song, err error) { //(url, title, time string, err error)

	ctx := context.Background()

	service, err := youtube.NewService(ctx, option.WithAPIKey(""))
	if err != nil {
		log.Fatalf("Error creating new YouTube client: %v", err)
	}

	call := service.Search.List([]string{"snippet"}).Q(searchString).MaxResults(1).Type("video")

	response, err := call.Do()
	if err != nil {
		return song_struct, fmt.Errorf("error making YouTube search API call: %v", err)
	}

	item := response.Items[0]
	song_struct = Song{
		ChannelID: item.Id.ChannelId,
		User:      item.Snippet.ChannelTitle,
		VidID:     item.Id.VideoId,
		Title:     item.Snippet.Title,
		Duration:  "0:00",
		VideoURL:  fmt.Sprintf("https://www.youtube.com/watch?v=%s", item.Id.VideoId),
	}

	return song_struct, nil

}
