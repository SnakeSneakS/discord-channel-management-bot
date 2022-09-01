package entity

import "gorm.io/gorm"

//Channelとユーザの管理
type DiscordChannelUser struct {
	gorm.Model

	GuildID   string
	ChannelID string
	UserID    string
	IsActive  bool
}
