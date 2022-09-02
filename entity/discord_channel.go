package entity

import (
	"database/sql"
	"time"

	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

// このbotが管理するChannel
type DiscordChannel struct {
	gorm.Model

	GuildID      string
	ChannelID    string
	ChannelName  string
	ChannelTopic sql.NullString

	IsPrivate bool
	IsArchive bool

	CreatedByUserID string
	LastMessageTime time.Time `gorm:"autoCreateTime"`

	ChannelType discordgo.ChannelType
}

type DiscordChannelSetting struct {
	gorm.Model

	GuildID                     string
	ParentCategoryID            string
	DescriptionChannelID        string
	DescriptionChannelMessageID string
}
