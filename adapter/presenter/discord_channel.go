package presenter

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/snakesneaks/discord-channel-management-bot/entity"
	"github.com/snakesneaks/discord-channel-management-bot/usecase/port"
)

var _ port.DiscordChannelOutputPort = (*discordChannel)(nil)

type discordChannel struct {
}

func NewDiscordChannelOutputPort() port.DiscordChannelOutputPort {
	return discordChannel{}
}

func (c discordChannel) ShowChannels(s *discordgo.Session, i *discordgo.Interaction, channels []*entity.DiscordChannel) string {
	content := "channels: \n"
	for _, c := range channels {
		isPrivateText := ""
		if c.IsPrivate {
			isPrivateText = "🔒"
		}
		content += fmt.Sprintf("[%s] [%s] %s\n", isPrivateText, c.ChannelName, c.ChannelTopic.String)
	}
	return content
}
