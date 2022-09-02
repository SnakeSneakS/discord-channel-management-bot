package gateway

import (
	"fmt"
	"time"

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
	stored, err := r.GetChannel(channel.GuildID, channel.ChannelID)
	if err != nil {
		return err
	}

	channel.ID = stored.ID
	channel.CreatedAt = stored.CreatedAt
	channel.UpdatedAt = time.Now()
	channel.DeletedAt = stored.DeletedAt
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
