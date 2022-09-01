package interactor

import (
	"github.com/bwmarrin/discordgo"
	"github.com/snakesneaks/discord-channel-management-bot/entity"
	"github.com/snakesneaks/discord-channel-management-bot/usecase/port"
)

var _ port.DiscordChannelInputPort = (*discordChannelInteractor)(nil)

type DiscordChannelInteractor interface {
	CreateChannel(channel *entity.DiscordChannel) error
	UpdateChannel(channel *entity.DiscordChannel) error
	DeleteChannel(guildID, channelID string) error
	GetChannel(guildID, channelID string) (*entity.DiscordChannel, error)
	GetChannels(guildID string) ([]*entity.DiscordChannel, error)
	GetSetting(s *discordgo.Session) (*entity.DiscordChannelSetting, error)
	CreateOrUpdateSetting(*entity.DiscordChannelSetting) error
}

type discordChannelInteractor struct {
	outputPort port.DiscordChannelOutputPort
	repository port.DiscordChannelRepository
}

func NewDiscordChannelInputPort(outputPort port.DiscordChannelOutputPort, repository port.DiscordChannelRepository) port.DiscordChannelInputPort {
	return &discordChannelInteractor{
		outputPort: outputPort,
		repository: repository,
	}
}

func (interactor discordChannelInteractor) CreateChannel(channel *entity.DiscordChannel) error {
	return interactor.repository.Create(channel)
}

func (interactor discordChannelInteractor) UpdateChannel(channel *entity.DiscordChannel) error {
	return interactor.repository.Update(channel)
}

func (interactor discordChannelInteractor) DeleteChannel(guildID, channelID string) error {
	return interactor.repository.Delete(guildID, channelID)
}

func (interactor discordChannelInteractor) GetChannel(guildID, channelID string) (*entity.DiscordChannel, error) {
	return interactor.repository.GetChannel(guildID, channelID)
}

func (interactor discordChannelInteractor) GetChannels(guildID string) ([]*entity.DiscordChannel, error) {
	return interactor.repository.GetChannels(guildID)
}

func (interactor discordChannelInteractor) GetSetting(guildID string) (*entity.DiscordChannelSetting, bool, error) {
	return interactor.repository.GetSetting(guildID)
}

func (interactor discordChannelInteractor) CreateOrUpdateSetting(s *entity.DiscordChannelSetting) error {
	return interactor.repository.CreateOrUpdateSetting(s)
}
