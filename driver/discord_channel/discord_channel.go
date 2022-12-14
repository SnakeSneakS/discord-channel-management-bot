package discord_channel

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/snakesneaks/discord-channel-management-bot/adapter/controller"
	"github.com/snakesneaks/discord-channel-management-bot/adapter/gateway"
	"github.com/snakesneaks/discord-channel-management-bot/adapter/presenter"
	"github.com/snakesneaks/discord-channel-management-bot/driver/db"
	"github.com/snakesneaks/discord-channel-management-bot/entity"
	"github.com/snakesneaks/discord-channel-management-bot/usecase/interactor"
)

type DiscordChannelDriver interface {
	GetChannel(guildID, channelID string) (*entity.DiscordChannel, error)
	GetChannels(guildID string) ([]*entity.DiscordChannel, error)

	DeleteChannel(guildID, channelID string) error
	CreateChannel(guildID, channelID, channelName, channelTopic, userID string, isPrivate bool) error
	UpdateChannel(guildID, channelID, channelName, channelTopic string, isPrivate bool) error
}

type discordChannelDriver struct {
	controller controller.DiscordChannelController
}

func newDiscordChannelController() controller.DiscordChannelController {
	db, err := db.GetDB()
	if err != nil {
		log.Fatal(err)
	}
	return controller.NewDiscordChannelController(
		db,
		interactor.NewDiscordChannelInputPort,
		presenter.NewDiscordChannelOutputPort,
		gateway.NewDiscordChannelRepository,
	)
}

func NewDiscordChannelDriver() DiscordChannelDriver {
	return discordChannelDriver{
		controller: newDiscordChannelController(),
	}
}

func (d discordChannelDriver) GetChannel(guildID, channelID string) (*entity.DiscordChannel, error) {
	return d.controller.GetChannel(guildID, channelID)
}
func (d discordChannelDriver) GetChannels(guildID string) ([]*entity.DiscordChannel, error) {
	return d.controller.GetChannels(guildID)
}

func (d discordChannelDriver) DeleteChannel(guildID, channelID string) error {
	return d.controller.DeleteChannel(guildID, channelID)
}

func (d discordChannelDriver) CreateChannel(guildID, channelID, channelName, channelTopic, userID string, isPrivate bool) error {
	return d.controller.CreateChannel(guildID, channelID, channelName, channelTopic, userID, isPrivate, discordgo.ChannelTypeGuildText)
}

func (d discordChannelDriver) UpdateChannel(guildID, channelID, channelName, channelTopic string, isPrivate bool) error {
	return d.controller.UpdateChannel(guildID, channelID, channelName, channelTopic, isPrivate)
}
