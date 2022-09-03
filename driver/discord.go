package driver

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/snakesneaks/discord-channel-management-bot/driver/discord"
	"github.com/snakesneaks/discord-channel-management-bot/driver/discord_channel"
	"github.com/snakesneaks/discord-channel-management-bot/driver/discord_channel_user"
	"github.com/snakesneaks/discord-channel-management-bot/driver/discord_setting"
	"github.com/snakesneaks/discord-channel-management-bot/entity"
	"gorm.io/gorm"
)

const (
	categoryName    = "discord-channel-management-bot-category"
	categoryTopic   = "this is the category discord-channel-management-bot handle."
	channelName     = "general"
	channelTopic    = "this is a channel where this bot show imformation."
	categoryPositon = 99999
)

// when message is automatically deleted in messageDescriptionChannel
var messageRemainDuration = time.Minute

// error for validate setting
var (
	errParentCategoryNotExist     = errors.New("parent category is not exist")
	errDescriptionChannelNotExist = errors.New("description channel is not exist")
	errDescriptionMessageNotExist = errors.New("description message is not exist")
)

// error for validate channels
var (
	errDBchannelNotExist        = errors.New("channel not exist in database")
	errdiscordChannelNotExist   = errors.New("channel not exist in discord")
	errParentCategoryIDmismatch = errors.New("channel found, but parent category id mismatch")
)

type DiscordDriver interface {
	CreateChannel(s *discordgo.Session, i *discordgo.Interaction, guildID, channelName, channelTopic string, isPrivate bool) (*discordgo.Channel, error)
	JoinChannel(s *discordgo.Session, guildID, channelID, userID string) error
	InviteUserToChannel(s *discordgo.Session, subjectUserID, guildID, channelID, userID string) error
	LeaveChannel(s *discordgo.Session, guildID, channelID, userID string) error
	DeleteChannel(s *discordgo.Session, guildID, channelID string) (*discordgo.Channel, error)
	UpdateChannel(s *discordgo.Session, subjectUserID, guildID, channelID, channelName, channelTopic string, isPrivate bool) error
	ShowInfo(s *discordgo.Session, guildID string) error
	GetChannel(guildID, channelID string) (*entity.DiscordChannel, error)

	GetChannelsInGuild(guildID string) ([]*entity.DiscordChannel, error)
	GetChannelUsersOfUser(guildID, userID string) ([]*entity.DiscordChannelUser, error)
}

type discordDriver struct {
	discordChannelDriver     discord_channel.DiscordChannelDriver
	discordChannelUserDriver discord_channel_user.DiscordChannelUserDriver
	discordSettingDriver     discord_setting.DiscordSettingDriver
}

func NewDiscordDriver(
	discordChannelDriver discord_channel.DiscordChannelDriver,
	discordChannelUserDriver discord_channel_user.DiscordChannelUserDriver,
	discordSettingDriver discord_setting.DiscordSettingDriver,
) DiscordDriver {
	return discordDriver{
		discordChannelDriver:     discordChannelDriver,
		discordChannelUserDriver: discordChannelUserDriver,
		discordSettingDriver:     discordSettingDriver,
	}
}

func (d discordDriver) InitSetting(s *discordgo.Session, guildID string) (*entity.DiscordChannelSetting, error) {
	category, err := discord.CreateCategory(s, guildID, categoryName, categoryTopic, categoryPositon)
	if err != nil {
		return nil, fmt.Errorf("error: failed to create category. %v", err)
	}

	c, err := discord.CreateChannel(s, guildID, channelName, channelTopic, category.ID, false, discordgo.ChannelTypeGuildText)
	if err != nil {
		return nil, fmt.Errorf("error: failed to create channel for bot. %v", err)
	}

	m, err := s.ChannelMessageSend(c.ID, "<<< DISCORD CHANNEL MANAGEMENT BOT IS HERE >>>")
	if err != nil {
		return nil, fmt.Errorf("error: failed to create message for bot. %v", err)
	}

	setting, err := d.discordSettingDriver.CreateOrUpdateSetting(category.GuildID, category.ID, c.ID, m.ID)
	if err != nil {
		return nil, fmt.Errorf("error: failed to create setting in database. %v", err)
	}

	return setting, nil
}

func (d discordDriver) getOrInitSetting(s *discordgo.Session, guildID string) (*entity.DiscordChannelSetting, error) {
	setting, found, err := d.discordSettingDriver.GetSetting(guildID)
	if err != nil {
		return nil, fmt.Errorf("error: failed to get setting. %v", err)
	}

	if !found {
		setting, err = d.InitSetting(s, guildID)
		if err != nil {
			return nil, err
		}
		return setting, err
	} else {
		//普通に見つかった場合は、間違っていないかvalidateする。
		if err := d.validateSetting(s, guildID, setting); err != nil {
			log.Printf("setting validation failed. so recreate setting. %v\n", err)
			if err == errParentCategoryNotExist || err == errDescriptionChannelNotExist || err == errDescriptionMessageNotExist {
				setting, err = d.InitSetting(s, guildID)
				if err != nil {
					return nil, err
				}
				return setting, nil
			}
		}
	}

	return setting, nil
}

// getChannel get channel with validation
func (d discordDriver) getChannel(s *discordgo.Session, guildID, channelID string) (*entity.DiscordChannel, error) {
	//setting
	setting, err := d.getOrInitSetting(s, guildID)
	if err != nil {
		return nil, err
	}

	//validate channel
	if discordChannel, err := d.validateChannel(s, guildID, channelID, setting); err != nil {
		switch err {

		case errDBchannelNotExist: //discordにはあるがDBにはない場合、discordでカテゴリから除外する
			if discordChannel.ParentID == setting.ParentCategoryID {
				if _, err := s.ChannelEditComplex(discordChannel.ID, &discordgo.ChannelEdit{
					ParentID: "",
				}); err != nil {
					return nil, err
				}
				return nil, fmt.Errorf("error: channel %s not found in database. so excluded from category", discordChannel.Mention())
			} else {
				return nil, fmt.Errorf("error: channel %s not found in database. this channel may not be managed by this bot", discordChannel.Mention())
			}

		case errdiscordChannelNotExist: //discordにはないがDBにはある場合、DBから除外する
			if err := d.discordChannelDriver.DeleteChannel(guildID, channelID); err != nil {
				return nil, err
			}
			return nil, fmt.Errorf("DB has data of non-existent channel, so deleted it. please try again")

		case errParentCategoryIDmismatch: //チャンネルは存在するがparentIDが異なる場合
			if _, err := s.ChannelEditComplex(discordChannel.ID, &discordgo.ChannelEdit{
				ParentID: setting.ParentCategoryID,
			}); err != nil {
				return nil, err
			}
			return nil, fmt.Errorf("channel %s category mismatch. so moved. please try again", discordChannel.Mention())

		default: //default, return error
			return nil, err
		}
	}

	//get channel
	channel, err := d.discordChannelDriver.GetChannel(guildID, channelID)
	if err != nil {
		return nil, err
	}

	return channel, nil
}

// getChannels get channels with validation
func (d discordDriver) getChannels(s *discordgo.Session, guildID string) ([]*entity.DiscordChannel, error) {
	channels, err := d.discordChannelDriver.GetChannels(guildID)
	if err != nil {
		return nil, err
	}

	for _, channel := range channels {
		if _, err := d.getChannel(s, guildID, channel.ChannelID); err != nil {
			return nil, err
		}
	}

	return channels, nil
}

func (d discordDriver) CreateChannel(s *discordgo.Session, i *discordgo.Interaction, guildID, channelName, channelTopic string, isPrivate bool) (*discordgo.Channel, error) {
	//setting
	setting, err := d.getOrInitSetting(s, guildID)
	if err != nil {
		return nil, err
	}

	//create channel in discord
	c, err := discord.CreateChannel(s, guildID, channelName, channelTopic, setting.ParentCategoryID, isPrivate, discordgo.ChannelTypeGuildText)
	if err != nil {
		return nil, err
	}

	//create channel in database
	if err := d.discordChannelDriver.CreateChannel(i.GuildID, c.ID, c.Name, c.Topic, i.Member.User.ID, isPrivate); err != nil {
		return c, fmt.Errorf("error: failed to create channel in database. %v", err)
	}
	//deny access toward created channel in discord
	if err := discord.DenyAllRolesToChannel(s, i.GuildID, c.ID); err != nil {
		return c, err
	}

	//join user who created this channel in discord
	if err := discord.SetMemberPermissionToChannel(s, c.ID, i.Member.User.ID, true); err != nil {
		return c, fmt.Errorf("error: failed to join channel. %v", err)
	}

	//join user who created this channel in database
	if err := d.discordChannelUserDriver.JoinOrLeaveChannel(i.GuildID, c.ID, i.Member.User.ID, true); err != nil {
		return c, fmt.Errorf("error: failed to join channel in database. %v", err)
	}

	return c, nil
}

// JoinChannel join channel
func (d discordDriver) JoinChannel(s *discordgo.Session, guildID, channelID, userID string) error {
	//get channel
	channel, err := d.getChannel(s, guildID, channelID)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}

	//authorization check
	if channel.IsPrivate {
		return fmt.Errorf("error: channel {%s} is a private channel. you must be invited by user in that channel", channelID)
	}

	//join channel in discord
	if err := discord.SetMemberPermissionToChannel(s, channelID, userID, true); err != nil {
		return fmt.Errorf("error: failed to join channel. %v", err)
	}

	//join channel in database
	if err := d.discordChannelUserDriver.JoinOrLeaveChannel(guildID, channelID, userID, true); err != nil {
		return fmt.Errorf("error: failed to join channel in database. %v", err)
	}

	return nil
}

// InviteUserToChannel
func (d discordDriver) InviteUserToChannel(s *discordgo.Session, subjectUserID, guildID, channelID, userID string) error {
	// authorization check
	_, err := d.discordChannelUserDriver.GetChannelUserInChannel(guildID, channelID, subjectUserID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return fmt.Errorf("error: %v", err)
	}
	if err == gorm.ErrRecordNotFound {
		return fmt.Errorf("error: you must be in channel {%s} to invite other member", channelID)
	}

	// invete uset to channel in discord
	if err := discord.SetMemberPermissionToChannel(s, channelID, userID, true); err != nil {
		return fmt.Errorf("error: failed to join channel. %v", err)
	}

	// invite user to channel in database
	if err := d.discordChannelUserDriver.JoinOrLeaveChannel(guildID, channelID, userID, true); err != nil {
		return fmt.Errorf("error: failed to join channel in database. %v", err)
	}

	return nil
}

// LeaveChannel
func (d discordDriver) LeaveChannel(s *discordgo.Session, guildID, channelID, userID string) error {
	//join in discord
	if err := discord.SetMemberPermissionToChannel(s, channelID, userID, false); err != nil {
		return err
	}

	//join in database
	if err := d.discordChannelUserDriver.JoinOrLeaveChannel(guildID, channelID, userID, false); err != nil {
		return fmt.Errorf("error: failed to leave channel in database. %v", err)
	}

	//success
	return nil
}

func (d discordDriver) DeleteChannel(s *discordgo.Session, guildID, channelID string) (*discordgo.Channel, error) {
	//delete channel in discord
	c, err := discord.DeleteChannel(s, channelID)
	if err != nil {
		return nil, fmt.Errorf("error: failed to delete channel. %v", err)
	}

	//delete channel in database
	if err := d.discordChannelDriver.DeleteChannel(guildID, channelID); err != nil {
		return c, fmt.Errorf("error: failed to delete channel from database. %v", err)
	}
	if err := d.discordChannelUserDriver.DeleteChannelUsersOfChannel(guildID, channelID); err != nil {
		return c, fmt.Errorf("error: failed to delete channel-user-relation from database. %v", err)
	}

	return c, nil
}

func (d discordDriver) UpdateChannel(s *discordgo.Session, subjectUserID, guildID, channelID, channelName, channelTopic string, isPrivate bool) error {
	//user authorization
	if _, err := d.discordChannelUserDriver.GetChannelUserInChannel(guildID, channelID, subjectUserID); err != nil {
		return fmt.Errorf("user is not in the channel. access denied")
	}

	//update channel in discord
	setting, err := d.getOrInitSetting(s, guildID)
	if err != nil {
		return err
	}
	channel, err := d.validateChannel(s, guildID, channelID, setting)
	if err != nil {
		return err
	}
	if _, err := s.ChannelEditComplex(channel.ID, &discordgo.ChannelEdit{
		Name:  channelName,
		Topic: channelTopic,
	}); err != nil {
		return err
	}

	//update channel in database
	if err := d.discordChannelDriver.UpdateChannel(guildID, channelID, channelName, channelTopic, isPrivate); err != nil {
		return err
	}

	return nil
}

// ShowInfo show information
func (d discordDriver) ShowInfo(s *discordgo.Session, guildID string) error {

	channels, err := d.getChannels(s, guildID)
	if err != nil {
		return err
	}

	setting, err := d.getOrInitSetting(s, guildID)
	if err != nil {
		return err
	}

	content := "DISCORD CHANNEL MANAGEMENT BOT\n"

	//info
	//content += "\nsetting: \n"
	//content += fmt.Sprintf("GuildID: %s\nCategoryID: %s\nChannelID: %s\nMessageID: %s\n", setting.GuildID, setting.ParentCategoryID, setting.DescriptionChannelID, setting.DescriptionChannelMessageID)

	//channels
	content += fmt.Sprintf("\nchannels: %d\n", len(channels))
	publicChannels := "**[public]** join yourself. `/join`\n"
	privateChannels := "**[private]** invitation needed. `/invite`\n"
	for _, channel := range channels {
		num, err := d.discordChannelUserDriver.GetChannelUserNumInChannel(guildID, channel.ChannelID)
		if err != nil {
			return err
		}

		if channel.IsPrivate {
			//privateChannels += fmt.Sprintf("`%s` **#%s** (num: %d)\n", channel.ChannelID, channel.ChannelName, num)
			privateChannels += fmt.Sprintf("Name: **%s**\nTopic:   %s\nID:        %s\nnum:     %d\n\n", channel.ChannelName, channel.ChannelID, channel.ChannelTopic.String, num)
		} else {
			//publicChannels += fmt.Sprintf("`%s` **#%s** (num: %d)\n", channel.ChannelID, channel.ChannelName, num)
			publicChannels += fmt.Sprintf("Name: **%s**\nTopic:   %s\nID:        %s\nnum:     %d\n\n", channel.ChannelName, channel.ChannelID, channel.ChannelTopic.String, num)
		}
	}
	content += fmt.Sprintf("\n%s\n%s\n", publicChannels, privateChannels)

	//help
	content += "\nhelp: \n"
	content += "run command `/help` to show help."

	_, err = s.ChannelMessageEdit(setting.DescriptionChannelID, setting.DescriptionChannelMessageID, content)
	if err != nil {
		return err
	}

	return nil
}

// validateSetting
func (d discordDriver) validateSetting(s *discordgo.Session, guildID string, setting *entity.DiscordChannelSetting) error {
	//parent check
	if _, err := discord.GetChannel(s, guildID, setting.ParentCategoryID); err != nil {
		log.Println(err)
		return errParentCategoryNotExist
	}
	if _, err := discord.GetChannel(s, guildID, setting.DescriptionChannelID); err != nil {
		log.Println(err)
		return errDescriptionChannelNotExist
	}

	//message check
	isDescriptionMessageExist := false
	messages, err := s.ChannelMessages(setting.DescriptionChannelID, 10, "", "", "")
	if err != nil {
		return err
	}
	for _, message := range messages {
		if message.ID == setting.DescriptionChannelMessageID {
			isDescriptionMessageExist = true
		} else {
			if time.Since(message.Timestamp) > messageRemainDuration {
				if err := s.ChannelMessageDelete(message.ChannelID, message.ID); err != nil {
					return err
				}
			}
		}
	}
	if !isDescriptionMessageExist {
		return errDescriptionMessageNotExist
	}

	return nil
}

// validateChannel
func (d discordDriver) validateChannel(s *discordgo.Session, guildID, channelID string, setting *entity.DiscordChannelSetting) (*discordgo.Channel, error) {
	//discord channel
	channel, err := discord.GetChannel(s, guildID, channelID)
	if err != nil {
		log.Println(err)
		return nil, errdiscordChannelNotExist
	}

	//db channels
	_, err = d.discordChannelDriver.GetChannel(guildID, channelID)
	if err != nil {
		log.Println(err)
		return channel, errDBchannelNotExist
	}

	//parent category check
	if channel.ParentID != setting.ParentCategoryID {
		log.Println(err)
		return nil, errParentCategoryIDmismatch
	}

	return channel, nil
}

func (d discordDriver) GetChannel(guildID, channelID string) (*entity.DiscordChannel, error) {
	return d.discordChannelDriver.GetChannel(guildID, channelID)
}

func (d discordDriver) GetChannelsInGuild(guildID string) ([]*entity.DiscordChannel, error) {
	return d.discordChannelDriver.GetChannels(guildID)
}

func (d discordDriver) GetChannelUsersOfUser(guildID, userID string) ([]*entity.DiscordChannelUser, error) {
	return d.discordChannelUserDriver.GetChannelUserOfUser(guildID, userID)
}
