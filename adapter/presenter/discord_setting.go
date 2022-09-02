package presenter

import (
	"fmt"

	"github.com/snakesneaks/discord-channel-management-bot/entity"
	"github.com/snakesneaks/discord-channel-management-bot/usecase/port"
)

type discordSettingPresenter struct {
}

func NewDiscordSettingOutputPort() port.DiscordSettingOutputPort {
	return discordSettingPresenter{}
}

func (p discordSettingPresenter) ShowSetting(setting *entity.DiscordChannelSetting) string {
	return fmt.Sprintf("guildID: %s\ncategoryID: %s\ndescriptionChannelID: %s\ndescriptionMessageID: %s\n", setting.GuildID, setting.ParentCategoryID, setting.DescriptionChannelID, setting.DescriptionChannelMessageID)
}
