package port

import (
	"github.com/bwmarrin/discordgo"
	"github.com/snakesneaks/discord-channel-management-bot/entity"
)

type DiscordChannelInputPort interface {
	CreateChannel(channel *entity.DiscordChannel) error
	UpdateChannel(channel *entity.DiscordChannel) error

	DeleteChannel(guildID, channelID string) error

	GetChannel(guildID, channelID string) (*entity.DiscordChannel, error)
	GetChannels(guildID string) ([]*entity.DiscordChannel, error)
}

type DiscordChannelOutputPort interface {
	ShowChannels(s *discordgo.Session, i *discordgo.Interaction, channels []*entity.DiscordChannel) string
}

type DiscordChannelRepository interface {
	Create(channel *entity.DiscordChannel) error
	Update(channel *entity.DiscordChannel) error

	Delete(guildID, channelID string) error

	GetChannel(guildID, channelID string) (*entity.DiscordChannel, error)
	GetChannels(guildID string) ([]*entity.DiscordChannel, error)
}
