package driver

import (
	"log"

	"github.com/snakesneaks/discord-channel-management-bot/adapter/controller"
	"github.com/snakesneaks/discord-channel-management-bot/adapter/gateway"
	"github.com/snakesneaks/discord-channel-management-bot/adapter/presenter"
	"github.com/snakesneaks/discord-channel-management-bot/driver/db"
	"github.com/snakesneaks/discord-channel-management-bot/usecase/interactor"
)

type DiscordChannelUserDriver interface {
	DeleteChannelUsersOfChannel(guildID, channelID string) error
	JoinOrLeaveChannel(guildID, channelID, userID string, isJoin bool) error
}

type discordChannelUserDriver struct {
	controller controller.DiscordChannelUserController
}

func newDiscordChannelUserController() controller.DiscordChannelUserController {
	db, err := db.GetDB()
	if err != nil {
		log.Fatal(err)
	}
	return controller.NewDiscordChannelUserController(
		db,
		interactor.NewDiscordChannelUserInputPort,
		presenter.NewDiscordChannelUserOutputPort,
		gateway.NewDiscordChannelUserRepository,
	)
}

func NewDiscordChannelUserDriver() DiscordChannelUserDriver {
	return discordChannelUserDriver{
		controller: newDiscordChannelUserController(),
	}
}

func (d discordChannelUserDriver) JoinOrLeaveChannel(guildID, channelID, userID string, isJoin bool) error {
	if isJoin {
		return d.controller.JoinChannel(guildID, userID, channelID)
	} else {
		return d.controller.LeaveChannel(guildID, userID, channelID)
	}
}

func (d discordChannelUserDriver) DeleteChannelUsersOfChannel(guildID, channelID string) error {
	return d.controller.DeleteChannel(guildID, channelID)
}

/*
func (d discordChannelUserDriver) GetChannelUsersOfGuild(guildID string) ([]*entity.DiscordChannelUser, error)
func (d discordChannelUserDriver) GetChannelUsersOfUser(guildID, userID string) ([]*entity.DiscordChannelUser, error)
func (d discordChannelUserDriver) GetChannelUsersInChannel(guildID, channelID string) ([]*entity.DiscordChannelUser, error)
*/
