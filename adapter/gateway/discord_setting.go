package gateway

import (
	"fmt"
	"time"

	"github.com/snakesneaks/discord-channel-management-bot/entity"
	"github.com/snakesneaks/discord-channel-management-bot/usecase/port"
	"gorm.io/gorm"
)

type DiscordSettingRepository interface {
}

type discordSettingRepository struct {
	conn *gorm.DB
}

func NewDiscordSettingRepository(conn *gorm.DB) port.DiscordSettingRepository {
	return discordSettingRepository{
		conn: conn,
	}
}

// if found, return true with no error
func (r discordSettingRepository) GetSetting(guildID string) (*entity.DiscordChannelSetting, bool, error) {
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

func (r discordSettingRepository) CreateOrUpdateSetting(s *entity.DiscordChannelSetting) error {
	var stored entity.DiscordChannelSetting
	tx := r.conn.Where("guild_id = ?", s.GuildID).First(&stored)
	if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
		return tx.Error
	}
	if tx.RowsAffected != 1 || tx.Error == gorm.ErrRecordNotFound {
		tx := r.conn.Create(s)
		if tx.Error != nil {
			return tx.Error
		}
		if tx.RowsAffected != 1 {
			return fmt.Errorf("failed to create %+#v", s)
		}
	} else {
		s.ID = stored.ID
		s.CreatedAt = stored.CreatedAt
		s.UpdatedAt = time.Now()
		s.DeletedAt = stored.DeletedAt
		tx := r.conn.Save(s)
		if tx.Error != nil {
			return tx.Error
		}
	}

	return nil
}
