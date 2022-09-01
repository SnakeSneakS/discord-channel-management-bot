package presenter

import (
	"github.com/snakesneaks/discord-channel-management-bot/usecase/port"
)

type discordChannelUserPresenter struct {
}

func NewDiscordChannelUserOutputPort() port.DiscordChannelUserOutputPort {
	return discordChannelUserPresenter{}
}
