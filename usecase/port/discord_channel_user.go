package port

import "github.com/snakesneaks/discord-channel-management-bot/entity"

type DiscordChannelUserInputPort interface {
	JoinChannel(guildID, userID, channelID string) error
	LeaveChannel(guildID, userID, channelID string) error
	DeleteChannel(guildID, channelID string) error

	GetChannelUsersOfGuild(guildID string) ([]*entity.DiscordChannelUser, error)
	GetChannelUsersOfUser(guildID, userID string) ([]*entity.DiscordChannelUser, error)
	GetChannelUsersInChannel(guildID, channelID string) ([]*entity.DiscordChannelUser, error)
}

type DiscordChannelUserOutputPort interface {
}

type DiscordChannelUserRepository interface {
	JoinChannel(guildID, userID, channelID string) error
	LeaveChannel(guildID, userID, channelID string) error
	DeleteChannel(guildID, channelID string) error

	GetChannelUsersOfGuild(guildID string) ([]*entity.DiscordChannelUser, error)
	GetChannelUsersOfUser(guildID, userID string) ([]*entity.DiscordChannelUser, error)
	GetChannelUsersInChannel(guildID, channelID string) ([]*entity.DiscordChannelUser, error)
}
