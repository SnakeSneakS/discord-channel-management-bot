package port

import (
	"github.com/snakesneaks/discord-channel-management-bot/entity"
)

type DiscordSettingInputPort interface {
	ShowSetting(guildID string) string
	GetSetting(guildID string) (*entity.DiscordChannelSetting, bool, error)
	CreateOrUpdateSetting(*entity.DiscordChannelSetting) error
}

type DiscordSettingOutputPort interface {
	ShowSetting(*entity.DiscordChannelSetting) string
}

type DiscordSettingRepository interface {
	GetSetting(guildID string) (*entity.DiscordChannelSetting, bool, error)
	CreateOrUpdateSetting(*entity.DiscordChannelSetting) error
}
