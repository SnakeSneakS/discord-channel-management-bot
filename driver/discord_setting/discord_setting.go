package discord_setting

import (
	"log"

	"github.com/snakesneaks/discord-channel-management-bot/adapter/controller"
	"github.com/snakesneaks/discord-channel-management-bot/adapter/gateway"
	"github.com/snakesneaks/discord-channel-management-bot/adapter/presenter"
	"github.com/snakesneaks/discord-channel-management-bot/driver/db"
	"github.com/snakesneaks/discord-channel-management-bot/entity"
	"github.com/snakesneaks/discord-channel-management-bot/usecase/interactor"
)

type DiscordSettingDriver interface {
	ShowSetting(guildID string) string
	GetSetting(guildID string) (*entity.DiscordChannelSetting, bool, error)
	CreateOrUpdateSetting(guildID, categoryID, descriptionChannelID, descriptionChannelMessageID string) (*entity.DiscordChannelSetting, error)
}

type discordSettingDriver struct {
	controller controller.DiscordSettingController
}

func newDiscordSettingController() controller.DiscordSettingController {
	db, err := db.GetDB()
	if err != nil {
		log.Fatal(err)
	}
	return controller.NewDiscordSettingController(
		db,
		interactor.NewDiscordSettingInputPort,
		presenter.NewDiscordSettingOutputPort,
		gateway.NewDiscordSettingRepository,
	)
}

func NewDiscordSettingDriver() DiscordSettingDriver {
	return discordSettingDriver{
		controller: newDiscordSettingController(),
	}
}

func (d discordSettingDriver) ShowSetting(guildID string) string {
	return d.controller.ShowSetting(guildID)
}

func (d discordSettingDriver) GetSetting(guildID string) (*entity.DiscordChannelSetting, bool, error) {
	return d.controller.GetSetting(guildID)
}
func (d discordSettingDriver) CreateOrUpdateSetting(guildID, categoryID, descriptionChannelID, descriptionChannelMessageID string) (*entity.DiscordChannelSetting, error) {
	setting := &entity.DiscordChannelSetting{
		GuildID:                     guildID,
		ParentCategoryID:            categoryID,
		DescriptionChannelID:        descriptionChannelID,
		DescriptionChannelMessageID: descriptionChannelMessageID,
	}

	if err := d.controller.CreateOrUpdateSetting(setting); err != nil {
		return setting, err
	}
	return setting, nil
}
