package gateway

import (
	"fmt"

	"github.com/snakesneaks/discord-channel-management-bot/entity"
	"github.com/snakesneaks/discord-channel-management-bot/usecase/port"
	"gorm.io/gorm"
)

var _ port.DiscordChannelRepository = (*discordChannelRepository)(nil)

type discordChannelRepository struct {
	conn *gorm.DB
}

func NewDiscordChannelRepository(conn *gorm.DB) port.DiscordChannelRepository {
	return discordChannelRepository{
		conn: conn,
	}
}

func (r discordChannelRepository) Create(channel *entity.DiscordChannel) error {
	tx := r.conn.Create(channel)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected != 1 {
		return fmt.Errorf("tried to create channel %+#v, but no row affected", channel)
	}

	return nil
}

func (r discordChannelRepository) Update(channel *entity.DiscordChannel) error {
	if _, err := r.GetChannel(channel.GuildID, channel.ChannelID); err != nil {
		return err
	}

	tx := r.conn.Save(channel)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected != 1 {
		return fmt.Errorf("failed to update discordChannel, %+#v", channel)
	}

	return nil
}

func (r discordChannelRepository) Delete(guildID, channelID string) error {
	var channels []*entity.DiscordChannel
	tx := r.conn.Where("guild_id = ? and channel_id = ?", guildID, channelID).Delete(&channels)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return fmt.Errorf("failed to delete discordChannel, no row affected. guildID: %s, channelID: %s", guildID, channelID)
	}

	return nil
}

func (r discordChannelRepository) GetChannel(guildID, channelID string) (*entity.DiscordChannel, error) {
	var c_stored entity.DiscordChannel
	tx := r.conn.Where("guild_id = ? and channel_id = ?", guildID, channelID).First(&c_stored)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected != 1 {
		return nil, fmt.Errorf("no row affected when getting channel, guildID: %s, channelID: %s", guildID, channelID)
	}

	return &c_stored, nil
}

func (r discordChannelRepository) GetChannels(guildID string) ([]*entity.DiscordChannel, error) {
	var c_stored []*entity.DiscordChannel
	tx := r.conn.Where("guild_id = ?", guildID).Find(&c_stored)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return c_stored, nil
}

// if found, return true with no error
func (r discordChannelRepository) GetSetting(guildID string) (*entity.DiscordChannelSetting, bool, error) {
	var c_stored *entity.DiscordChannelSetting
	tx := r.conn.Where("guild_id = ?", guildID).First(&c_stored)
	if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
		return nil, false, tx.Error
	}
	if tx.RowsAffected != 1 || tx.Error == gorm.ErrRecordNotFound {
		return nil, false, nil
	}
	return c_stored, true, nil
}

func (r discordChannelRepository) CreateOrUpdateSetting(s *entity.DiscordChannelSetting) error {
	var stored entity.DiscordChannelSetting
	tx := r.conn.Where("guild_id = ?", s.GuildID).First(&stored)
	if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
		return tx.Error
	}
	if tx.RowsAffected == 1 || tx.Error == gorm.ErrRecordNotFound {
		tx := r.conn.Save(s)
		if tx.Error != nil {
			return tx.Error
		}
	} else {
		tx := r.conn.Create(s)
		if tx.Error != nil {
			return tx.Error
		}
		if tx.RowsAffected != 1 {
			return fmt.Errorf("failed to create %+#v", s)
		}
	}

	return nil
}
