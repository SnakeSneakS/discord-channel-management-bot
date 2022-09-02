package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// CreateChannel
func CreateChannel(s *discordgo.Session, guildID, channelName, channelTopic, categoryID string, isPrivate bool, channelType discordgo.ChannelType) (*discordgo.Channel, error) {
	c, err := s.GuildChannelCreate(guildID, channelName, channelType)
	if err != nil {
		return nil, fmt.Errorf("error: failed to create channel. %v", err)
	}

	c, err = s.ChannelEditComplex(c.ID, &discordgo.ChannelEdit{
		Name:                 channelName,
		Topic:                channelTopic,
		PermissionOverwrites: []*discordgo.PermissionOverwrite{},
		ParentID:             categoryID,
		//Archived:             &isArchive, //スレッドはアーカイブできるらしい?
	})
	if err != nil {
		return nil, fmt.Errorf("error: failed to create channel. %v", err)
	}
	/*
		if err := SetMemberPermissionToChannel(s, c.ID, i.Member.User.ID, true); err != nil {
			return c, err
		}
	*/

	return c, nil
}

// Channel閲覧権限を無に返す
func DenyAllRolesToChannel(s *discordgo.Session, guildID, channelID string) error {
	g, err := s.State.Guild(guildID)
	if err != nil {
		return fmt.Errorf("error: failed to get guild. %v", err)
	}

	for _, role := range g.Roles {
		if err := s.ChannelPermissionSet(
			channelID,
			role.ID,
			discordgo.PermissionOverwriteTypeRole,
			0,
			discordgo.PermissionAll,
		); err != nil {
			return fmt.Errorf("error: failed to change permission. %v", err)
		}
	}

	return nil
}

func DeleteChannel(s *discordgo.Session, channelID string) (*discordgo.Channel, error) {
	return s.ChannelDelete(channelID)
}

func SetMemberPermissionToChannel(s *discordgo.Session, channelID, userID string, isAllow bool) error {
	if isAllow {
		if err := s.ChannelPermissionSet(
			channelID,
			userID,
			discordgo.PermissionOverwriteTypeMember,
			discordgo.PermissionAllText,
			0,
		); err != nil {
			return fmt.Errorf("error: failed to change permission. %v", err)
		}
	} else {
		if err := s.ChannelPermissionSet(
			channelID,
			userID,
			discordgo.PermissionOverwriteTypeMember,
			0,
			discordgo.PermissionAll,
		); err != nil {
			return fmt.Errorf("error: failed to change permission. %v", err)
		}
	}

	return nil
}

func CreateCategory(s *discordgo.Session, guildID, categoryName, categoryTopic string, categoryPosition int) (*discordgo.Channel, error) {
	c, err := s.GuildChannelCreate(guildID, categoryName, discordgo.ChannelTypeGuildCategory)
	if err != nil {
		return nil, err
	}

	c, err = s.ChannelEditComplex(c.ID, &discordgo.ChannelEdit{
		Name:     c.Name,
		Topic:    c.Topic,
		Position: categoryPosition,
	})
	if err != nil {
		return nil, err
	}

	return c, nil
}

// Get Channel By ID
func GetChannel(s *discordgo.Session, guildID, channelID string) (*discordgo.Channel, error) {

	channel, err := s.State.Channel(channelID)
	if err != nil {
		return nil, fmt.Errorf("channel not found. guildID: %s, channelID: %s", guildID, channelID)
	}
	return channel, nil
}
