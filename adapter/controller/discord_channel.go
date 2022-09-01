package controller

import (
	"database/sql"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/snakesneaks/discord-channel-management-bot/entity"
	"github.com/snakesneaks/discord-channel-management-bot/usecase/port"
	"gorm.io/gorm"
)

type DiscordChannelController interface {
	GetChannel(guildID, channelID string) (*entity.DiscordChannel, error)
	GetChannels(guildID string) ([]*entity.DiscordChannel, error)

	CreateChannel(guildID, channelID, channelName, channelTopic, userID string, isPrivate bool, t discordgo.ChannelType) error
	UpdateChannelName(guildID, channelID string, newChannelName string) error
	UpdateChannelTopic(guildID, channelID string, newChannelTopic sql.NullString) error
	UpdateChannelPrivate(guildID, channelID string, newIsPrivate bool) error
	UpdateChannelArchive(guildID, channelID string, newIsArchive bool) error
	UpdateLastMessageTime(guildID, channelID string, newLastMessageTime time.Time) error
	DeleteChannel(guildID, channelID string) error

	GetSetting(guildID string) (*entity.DiscordChannelSetting, bool, error)
	CreateOrUpdateSetting(*entity.DiscordChannelSetting) error
}

type discordChannelController struct {
	conn          *gorm.DB
	InputFactory  func(port.DiscordChannelOutputPort, port.DiscordChannelRepository) port.DiscordChannelInputPort
	OutputFactory func() port.DiscordChannelOutputPort
	RepoFactory   func(conn *gorm.DB) port.DiscordChannelRepository
}

func NewDiscordChannelController(
	conn *gorm.DB,
	inputFactory func(port.DiscordChannelOutputPort, port.DiscordChannelRepository) port.DiscordChannelInputPort,
	outputFactory func() port.DiscordChannelOutputPort,
	repoFactory func(conn *gorm.DB) port.DiscordChannelRepository,
) DiscordChannelController {
	return discordChannelController{
		conn:          conn,
		InputFactory:  inputFactory,
		OutputFactory: outputFactory,
		RepoFactory:   repoFactory,
	}
}

func (c discordChannelController) GetChannel(guildID, channelID string) (*entity.DiscordChannel, error) {
	repo := c.RepoFactory(c.conn)
	channel, err := repo.GetChannel(guildID, channelID)
	if err != nil {
		return nil, err
	}
	return channel, nil
}

func (c discordChannelController) GetChannels(guildID string) ([]*entity.DiscordChannel, error) {
	repo := c.RepoFactory(c.conn)
	channels, err := repo.GetChannels(guildID)
	if err != nil {
		return nil, err
	}
	return channels, nil
}

func (c discordChannelController) CreateChannel(guildID, channelID, channelName, channelTopic, userID string, isPrivate bool, channelType discordgo.ChannelType) error {
	channel := entity.DiscordChannel{
		GuildID:         guildID,
		ChannelID:       channelID,
		ChannelName:     channelName,
		ChannelTopic:    sql.NullString{String: channelTopic, Valid: true},
		IsPrivate:       isPrivate,
		IsArchive:       false,
		CreatedByUserID: userID,
		ChannelType:     channelType,
	}
	repo := c.RepoFactory(c.conn)
	outputPort := c.OutputFactory()
	inputPort := c.InputFactory(outputPort, repo)
	return inputPort.CreateChannel(&channel)
}

func (c discordChannelController) UpdateChannelName(guildID, channelID string, newChannelName string) error {
	repoFactory := c.RepoFactory(c.conn)
	channel, err := repoFactory.GetChannel(guildID, channelID)
	if err != nil {
		return err
	}
	channel.ChannelName = newChannelName
	repo := c.RepoFactory(c.conn)
	outputPort := c.OutputFactory()
	inputPort := c.InputFactory(outputPort, repo)
	return inputPort.UpdateChannel(channel)
}

func (c discordChannelController) UpdateChannelTopic(guildID, channelID string, newChannelTopic sql.NullString) error {
	repoFactory := c.RepoFactory(c.conn)
	channel, err := repoFactory.GetChannel(guildID, channelID)
	if err != nil {
		return err
	}
	channel.ChannelTopic = newChannelTopic
	repo := c.RepoFactory(c.conn)
	outputPort := c.OutputFactory()
	inputPort := c.InputFactory(outputPort, repo)
	return inputPort.UpdateChannel(channel)
}

func (c discordChannelController) UpdateChannelPrivate(guildID, channelID string, newIsPrivate bool) error {
	repoFactory := c.RepoFactory(c.conn)
	channel, err := repoFactory.GetChannel(guildID, channelID)
	if err != nil {
		return err
	}
	channel.IsPrivate = newIsPrivate
	repo := c.RepoFactory(c.conn)
	outputPort := c.OutputFactory()
	inputPort := c.InputFactory(outputPort, repo)
	return inputPort.UpdateChannel(channel)
}

func (c discordChannelController) UpdateChannelArchive(guildID, channelID string, newIsArchive bool) error {
	repoFactory := c.RepoFactory(c.conn)
	channel, err := repoFactory.GetChannel(guildID, channelID)
	if err != nil {
		return err
	}
	channel.IsArchive = newIsArchive
	repo := c.RepoFactory(c.conn)
	outputPort := c.OutputFactory()
	inputPort := c.InputFactory(outputPort, repo)
	return inputPort.UpdateChannel(channel)
}

func (c discordChannelController) UpdateLastMessageTime(guildID, channelID string, newLastMessageTime time.Time) error {
	repoFactory := c.RepoFactory(c.conn)
	channel, err := repoFactory.GetChannel(guildID, channelID)
	if err != nil {
		return err
	}
	channel.LastMessageTime = newLastMessageTime
	repo := c.RepoFactory(c.conn)
	outputPort := c.OutputFactory()
	inputPort := c.InputFactory(outputPort, repo)
	return inputPort.UpdateChannel(channel)
}

func (c discordChannelController) DeleteChannel(guildID, channelID string) error {
	repo := c.RepoFactory(c.conn)
	outputPort := c.OutputFactory()
	inputPort := c.InputFactory(outputPort, repo)
	return inputPort.DeleteChannel(guildID, channelID)
}

func (c discordChannelController) GetSetting(guildID string) (*entity.DiscordChannelSetting, bool, error) {
	repo := c.RepoFactory(c.conn)
	outputPort := c.OutputFactory()
	inputPort := c.InputFactory(outputPort, repo)
	return inputPort.GetSetting(guildID)
}

func (c discordChannelController) CreateOrUpdateSetting(s *entity.DiscordChannelSetting) error {
	repo := c.RepoFactory(c.conn)
	outputPort := c.OutputFactory()
	inputPort := c.InputFactory(outputPort, repo)
	return inputPort.CreateOrUpdateSetting(s)
}
