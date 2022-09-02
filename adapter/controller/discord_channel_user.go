package controller

import (
	"github.com/snakesneaks/discord-channel-management-bot/entity"
	"github.com/snakesneaks/discord-channel-management-bot/usecase/port"
	"gorm.io/gorm"
)

type DiscordChannelUserController interface {
	JoinChannel(guildID, userID, channelID string) error
	LeaveChannel(guildID, userID, channelID string) error
	DeleteChannel(guildID, channelID string) error

	GetChannelUsersOfGuild(guildID string) ([]*entity.DiscordChannelUser, error)
	GetChannelUsersOfUser(guildID, userID string) ([]*entity.DiscordChannelUser, error)
	GetChannelUserInChannel(guildID, channelID, userID string) (*entity.DiscordChannelUser, error)
	GetChannelUsersInChannel(guildID, channelID string) ([]*entity.DiscordChannelUser, error)
}

type discordChannelUserController struct {
	conn          *gorm.DB
	InputFactory  func(port.DiscordChannelUserOutputPort, port.DiscordChannelUserRepository) port.DiscordChannelUserInputPort
	OutputFactory func() port.DiscordChannelUserOutputPort
	RepoFactory   func(*gorm.DB) port.DiscordChannelUserRepository
}

func NewDiscordChannelUserController(
	conn *gorm.DB,
	inputFactory func(port.DiscordChannelUserOutputPort, port.DiscordChannelUserRepository) port.DiscordChannelUserInputPort,
	outputFactory func() port.DiscordChannelUserOutputPort,
	repoFactory func(conn *gorm.DB) port.DiscordChannelUserRepository,
) DiscordChannelUserController {
	return discordChannelUserController{
		conn:          conn,
		InputFactory:  inputFactory,
		OutputFactory: outputFactory,
		RepoFactory:   repoFactory,
	}
}

func (c discordChannelUserController) JoinChannel(guildID, userID, channelID string) error {
	return c.InputFactory(c.OutputFactory(), c.RepoFactory(c.conn)).JoinChannel(guildID, userID, channelID)
}

func (c discordChannelUserController) LeaveChannel(guildID, userID, channelID string) error {
	return c.InputFactory(c.OutputFactory(), c.RepoFactory(c.conn)).LeaveChannel(guildID, userID, channelID)
}

func (c discordChannelUserController) DeleteChannel(guildID, channelID string) error {
	return c.InputFactory(c.OutputFactory(), c.RepoFactory(c.conn)).DeleteChannel(guildID, channelID)
}

func (c discordChannelUserController) GetChannelUsersOfGuild(guildID string) ([]*entity.DiscordChannelUser, error) {
	return c.InputFactory(c.OutputFactory(), c.RepoFactory(c.conn)).GetChannelUsersOfGuild(guildID)
}

func (c discordChannelUserController) GetChannelUsersOfUser(guildID, userID string) ([]*entity.DiscordChannelUser, error) {
	return c.InputFactory(c.OutputFactory(), c.RepoFactory(c.conn)).GetChannelUsersOfUser(guildID, userID)
}

func (c discordChannelUserController) GetChannelUserInChannel(guildID, channelID, userID string) (*entity.DiscordChannelUser, error) {
	return c.InputFactory(c.OutputFactory(), c.RepoFactory(c.conn)).GetChannelUserInChannel(guildID, channelID, userID)
}

func (c discordChannelUserController) GetChannelUsersInChannel(guildID, channelID string) ([]*entity.DiscordChannelUser, error) {
	return c.InputFactory(c.OutputFactory(), c.RepoFactory(c.conn)).GetChannelUsersInChannel(guildID, channelID)
}
