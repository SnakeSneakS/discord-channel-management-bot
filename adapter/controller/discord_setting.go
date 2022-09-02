package controller

import (
	"github.com/snakesneaks/discord-channel-management-bot/entity"
	"github.com/snakesneaks/discord-channel-management-bot/usecase/port"
	"gorm.io/gorm"
)

type DiscordSettingController interface {
	ShowSetting(guildID string) string
	GetSetting(guildID string) (*entity.DiscordChannelSetting, bool, error)
	CreateOrUpdateSetting(*entity.DiscordChannelSetting) error
}

type discordSettingController struct {
	conn          *gorm.DB
	InputFactory  func(port.DiscordSettingOutputPort, port.DiscordSettingRepository) port.DiscordSettingInputPort
	OutputFactory func() port.DiscordSettingOutputPort
	RepoFactory   func(conn *gorm.DB) port.DiscordSettingRepository
}

func NewDiscordSettingController(
	conn *gorm.DB,
	inputFactory func(port.DiscordSettingOutputPort, port.DiscordSettingRepository) port.DiscordSettingInputPort,
	outputFactory func() port.DiscordSettingOutputPort,
	repoFactory func(conn *gorm.DB) port.DiscordSettingRepository,
) DiscordSettingController {
	return discordSettingController{
		conn:          conn,
		InputFactory:  inputFactory,
		OutputFactory: outputFactory,
		RepoFactory:   repoFactory,
	}
}

func (c discordSettingController) ShowSetting(guildID string) string {
	repo := c.RepoFactory(c.conn)
	outputPort := c.OutputFactory()
	inputPort := c.InputFactory(outputPort, repo)
	return inputPort.ShowSetting(guildID)
}

func (c discordSettingController) GetSetting(guildID string) (*entity.DiscordChannelSetting, bool, error) {
	repo := c.RepoFactory(c.conn)
	outputPort := c.OutputFactory()
	inputPort := c.InputFactory(outputPort, repo)
	return inputPort.GetSetting(guildID)
}

func (c discordSettingController) CreateOrUpdateSetting(s *entity.DiscordChannelSetting) error {
	repo := c.RepoFactory(c.conn)
	outputPort := c.OutputFactory()
	inputPort := c.InputFactory(outputPort, repo)
	return inputPort.CreateOrUpdateSetting(s)
}
