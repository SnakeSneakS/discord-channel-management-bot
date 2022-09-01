package discord

import "github.com/bwmarrin/discordgo"

// CreateChannel
func CreateChannel(s *discordgo.Session, i *discordgo.Interaction, channelName, channelTopic, categoryID string, isPrivate bool) (*discordgo.Channel, error) {
	c, err := s.GuildChannelCreate(i.GuildID, channelName, discordgo.ChannelTypeGuildText)
	if err != nil {
		return nil, err
	}

	c, err = s.ChannelEditComplex(c.ID, &discordgo.ChannelEdit{
		Name:                 channelName,
		Topic:                channelTopic,
		PermissionOverwrites: []*discordgo.PermissionOverwrite{},
		ParentID:             categoryID,
		//Archived:             &isArchive, //スレッドはアーカイブできるらしい?
	})
	if err != nil {
		return nil, err
	}

	//channel permission
	g, err := s.State.Guild(i.GuildID)
	if err != nil {
		return nil, err
	}

	for _, role := range g.Roles {
		if err := s.ChannelPermissionSet(
			c.ID,
			role.ID,
			discordgo.PermissionOverwriteTypeRole,
			0,
			discordgo.PermissionAll,
		); err != nil {
			return c, err
		}
	}

	/*
		if err := SetMemberPermissionToChannel(s, c.ID, i.Member.User.ID, true); err != nil {
			return c, err
		}
	*/

	return c, nil
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
			return err
		}
	} else {
		if err := s.ChannelPermissionSet(
			channelID,
			userID,
			discordgo.PermissionOverwriteTypeMember,
			0,
			discordgo.PermissionAll,
		); err != nil {
			return err
		}
	}

	return nil
}

func CreateCategory(s *discordgo.Session, i *discordgo.Interaction, categoryName, categoryTopic string, categoryPosition int) (*discordgo.Channel, error) {
	c, err := s.GuildChannelCreate(i.GuildID, categoryName, discordgo.ChannelTypeGuildCategory)
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
