package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func CreateMessage(s *discordgo.Session, i *discordgo.Interaction, content string) {
	err := s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})
	if err != nil {
		FollowupMessage(s, i, fmt.Sprintf("Something went wrong: %s", err))
		return
	}
}

func EditMessage(s *discordgo.Session, i *discordgo.Interaction, content *string) {
	_, err := s.InteractionResponseEdit(i, &discordgo.WebhookEdit{
		Content: content,
	})
	if err != nil {
		FollowupMessage(s, i, fmt.Sprintf("Something went wrong: %s", err))
		return
	}
}

func FollowupMessage(s *discordgo.Session, i *discordgo.Interaction, content string) {
	s.FollowupMessageCreate(i, true, &discordgo.WebhookParams{
		Content: content,
	})
}
