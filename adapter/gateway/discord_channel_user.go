package gateway

import (
	"fmt"

	"github.com/snakesneaks/discord-channel-management-bot/entity"
	"github.com/snakesneaks/discord-channel-management-bot/usecase/port"
	"gorm.io/gorm"
)

type discordChannelUserRepository struct {
	conn *gorm.DB
}

func NewDiscordChannelUserRepository(conn *gorm.DB) port.DiscordChannelUserRepository {
	return discordChannelUserRepository{
		conn: conn,
	}
}

func (r discordChannelUserRepository) JoinChannel(guildID, userID, channelID string) error {
	var stored entity.DiscordChannelUser
	tx := r.conn.Where("guild_id = ? and user_id = ? and channel_id = ?", guildID, userID, channelID).First(&stored)
	if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
		return tx.Error
	}
	if tx.RowsAffected != 1 || tx.Error == gorm.ErrRecordNotFound {
		newDiscordChannelUserData := entity.DiscordChannelUser{
			GuildID:   guildID,
			UserID:    userID,
			ChannelID: channelID,
			IsActive:  true,
		}
		tx := r.conn.Create(&newDiscordChannelUserData)
		if tx.Error != nil {
			return tx.Error
		}
		if tx.RowsAffected != 1 {
			return fmt.Errorf("failed to create %+#v", newDiscordChannelUserData)
		}
	} else {
		if !stored.IsActive {
			stored.IsActive = true
			tx := r.conn.Save(stored)
			if tx.Error != nil {
				return tx.Error
			}
		}
	}

	return nil
}

func (r discordChannelUserRepository) LeaveChannel(guildID, userID, channelID string) error {
	var stored entity.DiscordChannelUser
	tx := r.conn.Where("guild_id = ? and user_id = ? and channel_id = ?", guildID, userID, channelID).First(&stored)
	if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
		return tx.Error
	}
	if tx.RowsAffected != 1 || tx.Error == gorm.ErrRecordNotFound {
		newDiscordChannelUserData := entity.DiscordChannelUser{
			GuildID:   guildID,
			UserID:    userID,
			ChannelID: channelID,
			IsActive:  false,
		}
		tx := r.conn.Create(newDiscordChannelUserData)
		if tx.Error != nil {
			return tx.Error
		}
		if tx.RowsAffected != 1 {
			return fmt.Errorf("failed to create %+#v", newDiscordChannelUserData)
		}
	} else {
		if stored.IsActive {
			stored.IsActive = false
			tx := r.conn.Save(stored)
			if tx.Error != nil {
				return tx.Error
			}
		}
	}

	return nil
}

func (r discordChannelUserRepository) DeleteChannel(guildID, channelID string) error {
	var channelUsers []*entity.DiscordChannelUser
	tx := r.conn.Where("guild_id = ? and channel_id = ?", guildID, channelID).Delete(&channelUsers)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return fmt.Errorf("no record matched")
	}
	return nil
}

func (r discordChannelUserRepository) GetChannelUsersOfGuild(guildID string) ([]*entity.DiscordChannelUser, error) {
	var discordChannelUsers []*entity.DiscordChannelUser
	tx := r.conn.Where("guild_id = ?", guildID).Find(&discordChannelUsers)
	if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
		return nil, tx.Error
	}
	return discordChannelUsers, nil
}

func (r discordChannelUserRepository) GetChannelUsersOfUser(guildID, userID string) ([]*entity.DiscordChannelUser, error) {
	var discordChannelUsers []*entity.DiscordChannelUser
	tx := r.conn.Where("guild_id = ? and user_id = ?", guildID, userID).Find(&discordChannelUsers)
	if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
		return nil, tx.Error
	}
	return discordChannelUsers, nil

}

func (r discordChannelUserRepository) GetChannelUserInChannel(guildID, channelID, userID string) (*entity.DiscordChannelUser, error) {
	var discordChannelUsers *entity.DiscordChannelUser
	tx := r.conn.Where("guild_id = ? and channel_id = ? and user_id = ?", guildID, channelID, userID).First(&discordChannelUsers)
	if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
		return nil, tx.Error
	}
	return discordChannelUsers, nil
}

func (r discordChannelUserRepository) GetChannelUsersInChannel(guildID, channelID string) ([]*entity.DiscordChannelUser, error) {
	var discordChannelUsers []*entity.DiscordChannelUser
	tx := r.conn.Where("guild_id = ? and channel_id = ?", guildID, channelID).Find(&discordChannelUsers)
	if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
		return nil, tx.Error
	}
	return discordChannelUsers, nil
}
