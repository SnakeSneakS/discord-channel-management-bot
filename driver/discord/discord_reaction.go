package discord

import (
	"fmt"
	"time"

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

func CreateMessageInstant(s *discordgo.Session, i *discordgo.Interaction, content string, durationSeconds time.Duration) {
	err := s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("%s (this message is deleted in  %s)", content, durationSeconds),
		},
	})
	if err != nil {
		FollowupMessage(s, i, fmt.Sprintf("Something went wrong: %s", err))
	}

	time.AfterFunc(durationSeconds, func() {
		if err := s.InteractionResponseDelete(i); err != nil {
			FollowupMessage(s, i, fmt.Sprintf("Something went wrong: %s", err))
			return
		}
	})
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

func FollowUpMessageInstant(s *discordgo.Session, i *discordgo.Interaction, content string, durationSeconds time.Duration) {
	m, err := s.FollowupMessageCreate(i, true, &discordgo.WebhookParams{
		Content: fmt.Sprintf("%s (this message is deleted in  %s)", content, durationSeconds),
	})
	if err != nil {
		FollowupMessage(s, i, err.Error())
	}

	time.AfterFunc(durationSeconds, func() {
		if err := s.ChannelMessageDelete(i.ChannelID, m.ID); err != nil {
			FollowupMessage(s, i, err.Error())
		}
	})
}
