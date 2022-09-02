package interactor

import (
	"github.com/snakesneaks/discord-channel-management-bot/entity"
	"github.com/snakesneaks/discord-channel-management-bot/usecase/port"
)

var _ port.DiscordChannelUserInputPort = (*discordChannelUserInteractor)(nil)

type discordChannelUserInteractor struct {
	outputPort port.DiscordChannelUserOutputPort
	repository port.DiscordChannelUserRepository
}

func NewDiscordChannelUserInputPort(outputPort port.DiscordChannelUserOutputPort, repository port.DiscordChannelUserRepository) port.DiscordChannelUserInputPort {
	return discordChannelUserInteractor{
		outputPort: outputPort,
		repository: repository,
	}
}

func (i discordChannelUserInteractor) JoinChannel(guildID, userID, channelID string) error {
	return i.repository.JoinChannel(guildID, userID, channelID)
}

func (i discordChannelUserInteractor) LeaveChannel(guildID, userID, channelID string) error {
	return i.repository.LeaveChannel(guildID, userID, channelID)
}

func (i discordChannelUserInteractor) DeleteChannel(guildID, channelID string) error {
	return i.repository.DeleteChannel(guildID, channelID)
}

func (i discordChannelUserInteractor) GetChannelUsersOfGuild(guildID string) ([]*entity.DiscordChannelUser, error) {
	return i.repository.GetChannelUsersOfGuild(guildID)
}

func (i discordChannelUserInteractor) GetChannelUsersOfUser(guildID, userID string) ([]*entity.DiscordChannelUser, error) {
	return i.repository.GetChannelUsersOfUser(guildID, userID)
}

func (i discordChannelUserInteractor) GetChannelUserInChannel(guildID, channelID, userID string) (*entity.DiscordChannelUser, error) {
	return i.repository.GetChannelUserInChannel(guildID, channelID, userID)
}
func (i discordChannelUserInteractor) GetChannelUsersInChannel(guildID, channelID string) ([]*entity.DiscordChannelUser, error) {
	return i.repository.GetChannelUsersInChannel(guildID, channelID)
}
