package interactor

import (
	"fmt"

	"github.com/snakesneaks/discord-channel-management-bot/entity"
	"github.com/snakesneaks/discord-channel-management-bot/usecase/port"
)

type discordSettingInteractor struct {
	outputPort port.DiscordSettingOutputPort
	repository port.DiscordSettingRepository
}

func NewDiscordSettingInputPort(outputPort port.DiscordSettingOutputPort, repository port.DiscordSettingRepository) port.DiscordSettingInputPort {
	return discordSettingInteractor{
		outputPort: outputPort,
		repository: repository,
	}
}

func (i discordSettingInteractor) GetSetting(guildID string) (*entity.DiscordChannelSetting, bool, error) {
	return i.repository.GetSetting(guildID)
}
func (i discordSettingInteractor) CreateOrUpdateSetting(setting *entity.DiscordChannelSetting) error {
	return i.repository.CreateOrUpdateSetting(setting)
}

func (i discordSettingInteractor) ShowSetting(guildID string) string {
	content := "setting: \n"
	setting, b, err := i.GetSetting(guildID)
	if err != nil {
		content += fmt.Sprintf("error: %v\n", err)
	} else if !b {
		content += "error: setting not found!"
	} else {
		content += i.outputPort.ShowSetting(setting)
	}
	return content
}
