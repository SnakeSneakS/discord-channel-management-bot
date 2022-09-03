package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/snakesneaks/discord-channel-management-bot/driver"
	"github.com/snakesneaks/discord-channel-management-bot/driver/discord"
	"github.com/snakesneaks/discord-channel-management-bot/driver/discord_channel"
	"github.com/snakesneaks/discord-channel-management-bot/driver/discord_channel_user"
	"github.com/snakesneaks/discord-channel-management-bot/driver/discord_setting"
)

const instantMessageDuration = time.Second * 20

const (
	commandHelp              = "help"
	commandCreateChannel     = "create"
	commandJoinChannel       = "join"
	commandInviteChannelUser = "invite"
	commandLeaveChannel      = "leave"
	commandDeleteChannel     = "delete"
	commandUpdateChannel     = "update"
	commandShowChannels      = "show-channels"
	commandTest              = "test"

	optionChannelName      = "channel-name"
	optionChannelTopic     = "channel-topic"
	optionChannelIsPrivate = "channel-is-private"

	optionChannelID = "channel-id"
	optionUser      = "user"
	optionChannel   = "channel"
)

func newDiscordCommands() []*discordgo.ApplicationCommand {
	defaultAllow := true
	var permissionManageServer int64 = discordgo.PermissionManageServer
	commands := []*discordgo.ApplicationCommand{
		{
			Name:              commandHelp,
			Description:       "show help message",
			DefaultPermission: &defaultAllow,
		},
		{
			Name:              commandCreateChannel,
			Description:       "create channel",
			DefaultPermission: &defaultAllow,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        optionChannelName,
					Description: "channel name",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        optionChannelTopic,
					Description: "channel topic. explain what this channel do.",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Name:        optionChannelIsPrivate,
					Description: "if true, channel is private",
					Required:    true,
				},
			},
		},
		{
			Name:        commandJoinChannel,
			Description: "join channel",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        optionChannelID,
					Description: "channel id",
					Required:    true,
				},
			},
		},
		{
			Name:        commandInviteChannelUser,
			Description: "invite user to channel",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        optionChannel,
					Description: "channel",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        optionUser,
					Description: "user",
					Required:    true,
				},
			},
		},
		{
			Name:        commandLeaveChannel,
			Description: "leave channel",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        optionChannel,
					Description: "channel",
					Required:    true,
				},
			},
		},
		{
			Name:                     commandDeleteChannel,
			Description:              "delete channel (only manager)",
			DefaultMemberPermissions: &permissionManageServer,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        optionChannel,
					Description: "channel",
					Required:    true,
				},
			},
		},
		{
			Name:        commandUpdateChannel,
			Description: "update channel setting",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        optionChannel,
					Description: "channel",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        optionChannelName,
					Description: "channel name",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        optionChannelTopic,
					Description: "explanation what this channel do",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Name:        optionChannelIsPrivate,
					Description: "is channel private? if private, invitation is needed to join the channel",
					Required:    false,
				},
			},
		},
		{
			Name:        commandShowChannels,
			Description: "show channels",
		},
		/*
			{
				Name:        commandTest,
				Description: "test by developer",
			},
		*/
	}
	return commands
}

func getOptionMap(i *discordgo.Interaction) map[string]*discordgo.ApplicationCommandInteractionDataOption {
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, option := range options {
		optionMap[option.Name] = option
	}
	return optionMap
}

func newDiscordCommandHandler() map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	discordDriver := driver.NewDiscordDriver(
		discord_channel.NewDiscordChannelDriver(),
		discord_channel_user.NewDiscordChannelUserDriver(),
		discord_setting.NewDiscordSettingDriver(),
	)

	commandHandlers := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		//Command Help
		commandHelp: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			commands := "[command]: [description]\n"
			for _, v := range newDiscordCommands() {
				commands += fmt.Sprintf("`/%s`: %s\n", v.Name, v.Description)
			}
			discord.CreateMessageInstant(s, i.Interaction, commands, instantMessageDuration)
		},

		//Command Create Channel
		commandCreateChannel: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			discord.CreateMessageInstant(s, i.Interaction, "trying to create channel", instantMessageDuration)

			//get option
			optMap := getOptionMap(i.Interaction)
			channelName, ok := optMap[optionChannelName]
			if !ok {
				discord.FollowUpMessageInstant(s, i.Interaction, fmt.Sprintf("error: option %s is needed.", optionChannelName), instantMessageDuration)
				return
			}
			channelTopic, ok := optMap[optionChannelTopic]
			if !ok {
				discord.FollowUpMessageInstant(s, i.Interaction, fmt.Sprintf("error: option %s is needed.", optionChannelTopic), instantMessageDuration)
				return
			}
			isPrivate, ok := optMap[optionChannelIsPrivate]
			if !ok {
				discord.FollowUpMessageInstant(s, i.Interaction, fmt.Sprintf("error: option %s is needed.", optionChannelIsPrivate), instantMessageDuration)
				return
			}

			//create channel
			c, err := discordDriver.CreateChannel(s, i.Interaction, i.GuildID, channelName.StringValue(), channelTopic.StringValue(), isPrivate.BoolValue())
			if err != nil {
				discord.FollowUpMessageInstant(s, i.Interaction, err.Error(), instantMessageDuration)
				return
			}

			//success
			discord.FollowUpMessageInstant(s, i.Interaction, fmt.Sprintf("channel created!\n`id`: %s\n`name`: %s\n`topic`: %s", c.ID, c.Name, c.Topic), instantMessageDuration)
			if _, err = s.ChannelMessageSend(c.ID, fmt.Sprintf("Channel created!!\nID: %s\nName: %s\nTopic: %s\nIsPrivate: %t\n", c.ID, c.Name, c.Topic, isPrivate.BoolValue())); err != nil {
				discord.FollowUpMessageInstant(s, i.Interaction, err.Error(), instantMessageDuration)
				return
			}

			//show info
			if err := discordDriver.ShowInfo(s, i.GuildID); err != nil {
				discord.FollowUpMessageInstant(s, i.Interaction, err.Error(), instantMessageDuration)
				return
			}
		},

		//Command Join Channel
		commandJoinChannel: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			discord.CreateMessageInstant(s, i.Interaction, "trying to join channel", instantMessageDuration)

			//get option
			optMap := getOptionMap(i.Interaction)
			channelID, ok := optMap[optionChannelID]
			if !ok {
				discord.FollowUpMessageInstant(s, i.Interaction, fmt.Sprintf("error: option %s is needed.", optionChannelID), instantMessageDuration)
				return
			}

			//join
			if err := discordDriver.JoinChannel(s, i.GuildID, channelID.StringValue(), i.Member.User.ID); err != nil {
				discord.FollowUpMessageInstant(s, i.Interaction, err.Error(), instantMessageDuration)
				return
			}

			//success
			discord.FollowUpMessageInstant(s, i.Interaction, fmt.Sprintf("user %s joined into %s.", i.Member.User.Mention(), channelID.StringValue()), instantMessageDuration)
			if _, err := s.ChannelMessageSend(channelID.StringValue(), fmt.Sprintf("user joined! %s", i.Member.User.Mention())); err != nil {
				discord.FollowUpMessageInstant(s, i.Interaction, err.Error(), instantMessageDuration)
				return
			}
		},

		//Command Invite User
		commandInviteChannelUser: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			discord.CreateMessageInstant(s, i.Interaction, "trying to invite user to channel", instantMessageDuration)

			//get option
			optMap := getOptionMap(i.Interaction)
			channelOpt, ok := optMap[optionChannel]
			if !ok {
				discord.FollowUpMessageInstant(s, i.Interaction, fmt.Sprintf("error: option %s is needed.", optionChannel), instantMessageDuration)
				return
			}
			userOpt, ok := optMap[optionUser]
			if !ok {
				discord.FollowUpMessageInstant(s, i.Interaction, fmt.Sprintf("error: option %s is needed.", optionUser), instantMessageDuration)
				return
			}
			channel := channelOpt.ChannelValue(s)
			user := userOpt.UserValue(s)

			//invite
			if err := discordDriver.InviteUserToChannel(s, i.Member.User.ID, i.GuildID, channel.ID, user.ID); err != nil {
				discord.FollowUpMessageInstant(s, i.Interaction, err.Error(), instantMessageDuration)
				return
			}

			//success
			discord.FollowUpMessageInstant(s, i.Interaction, fmt.Sprintf("user %s joined into %s.", user.Mention(), channel.Mention()), instantMessageDuration)
			if _, err := s.ChannelMessageSend(channel.ID, fmt.Sprintf("user joined! %s", user.Mention())); err != nil {
				discord.FollowUpMessageInstant(s, i.Interaction, err.Error(), instantMessageDuration)
				return
			}
		},

		//Command Leave Channel
		commandLeaveChannel: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			discord.CreateMessageInstant(s, i.Interaction, "trying to leave channel", instantMessageDuration)

			optMap := getOptionMap(i.Interaction)
			channelOpt, ok := optMap[optionChannel]
			if !ok {
				discord.FollowUpMessageInstant(s, i.Interaction, fmt.Sprintf("error: option %s is needed.", optionChannel), instantMessageDuration)
				return
			}
			channel := channelOpt.ChannelValue(s)

			//leave channel
			if err := discordDriver.LeaveChannel(s, i.GuildID, channel.ID, i.Member.User.ID); err != nil {
				discord.FollowUpMessageInstant(s, i.Interaction, err.Error(), instantMessageDuration)
				return
			}

			//success
			discord.FollowUpMessageInstant(s, i.Interaction, fmt.Sprintf("user %s leave from %s.", i.Member.User.Mention(), channel.Mention()), instantMessageDuration)
			if _, err := s.ChannelMessageSend(channel.ID, fmt.Sprintf("user left! %s", i.Member.User.Mention())); err != nil {
				discord.FollowUpMessageInstant(s, i.Interaction, err.Error(), instantMessageDuration)
				return
			}
		},

		//Command Delete Channel
		commandDeleteChannel: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			discord.CreateMessageInstant(s, i.Interaction, "trying to delete channel", instantMessageDuration)

			//option
			optMap := getOptionMap(i.Interaction)
			channelOpt, ok := optMap[optionChannel]
			if !ok {
				discord.FollowUpMessageInstant(s, i.Interaction, fmt.Sprintf("error: option %s is needed.", optionChannel), instantMessageDuration)
				return
			}
			channel := channelOpt.ChannelValue(s)

			//delete channel
			c, err := discordDriver.DeleteChannel(s, i.GuildID, channel.ID)
			if err != nil {
				discord.FollowUpMessageInstant(s, i.Interaction, err.Error(), instantMessageDuration)
				return
			}

			//success
			discord.FollowUpMessageInstant(s, i.Interaction, fmt.Sprintf("deleted channel. Name: %s", c.Name), instantMessageDuration)

			//show info
			if err := discordDriver.ShowInfo(s, i.GuildID); err != nil {
				discord.FollowUpMessageInstant(s, i.Interaction, err.Error(), instantMessageDuration)
				return
			}
		},

		//Command Update Channel
		commandUpdateChannel: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			discord.CreateMessageInstant(s, i.Interaction, "trying to update channel setting", instantMessageDuration)

			//option
			optMap := getOptionMap(i.Interaction)
			channelOpt, ok := optMap[optionChannel]
			if !ok {
				discord.FollowUpMessageInstant(s, i.Interaction, fmt.Sprintf("error: option %s is needed.", optionChannel), instantMessageDuration)
				return
			}
			channel := channelOpt.ChannelValue(s)
			discordChannel, err := discordDriver.GetChannel(i.GuildID, channel.ID)
			if err != nil {
				discord.FollowUpMessageInstant(s, i.Interaction, err.Error(), instantMessageDuration)
				return
			}
			var channelName string
			channelNameOpt, ok := optMap[optionChannelName]
			if !ok {
				channelName = discordChannel.ChannelName
			} else {
				channelName = channelNameOpt.StringValue()
			}
			var channelTopic string
			channelTopicOpt, ok := optMap[optionChannelTopic]
			if !ok {
				channelTopic = discordChannel.ChannelTopic.String
			} else {
				channelTopicOpt.StringValue()
			}
			var isPrivate bool
			isPrivateOpt, ok := optMap[optionChannelIsPrivate]
			if !ok {
				isPrivate = discordChannel.IsPrivate
			} else {
				isPrivate = isPrivateOpt.BoolValue()
			}

			//update
			if err := discordDriver.UpdateChannel(s, i.Member.User.ID, i.GuildID, channel.ID, channelName, channelTopic, isPrivate); err != nil {
				discord.FollowUpMessageInstant(s, i.Interaction, err.Error(), instantMessageDuration)
				return
			}

			//success
			discord.FollowUpMessageInstant(s, i.Interaction, "channel setting updated.", instantMessageDuration)
			if _, err = s.ChannelMessageSend(channel.ID, fmt.Sprintf("Channel setting changed!!\nID: %s\nName: %s\nTopic: %s\nIsPrivate: %t\n", channel.ID, channel.Name, channel.Topic, isPrivate)); err != nil {
				discord.FollowUpMessageInstant(s, i.Interaction, err.Error(), instantMessageDuration)
				return
			}

			//show info
			if err := discordDriver.ShowInfo(s, i.GuildID); err != nil {
				discord.FollowUpMessageInstant(s, i.Interaction, err.Error(), instantMessageDuration)
				return
			}
		},

		//Command Update Info
		commandShowChannels: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			discord.CreateMessageInstant(s, i.Interaction, "trying to show channels", instantMessageDuration)

			//existing channels
			discordChannels, err := discordDriver.GetChannelsInGuild(i.GuildID)
			if err != nil {
				discord.FollowUpMessageInstant(s, i.Interaction, err.Error(), instantMessageDuration)
				return
			}

			//joining channels
			discordChannelUsers, err := discordDriver.GetChannelUsersOfUser(i.GuildID, i.Member.User.ID)
			if err != nil {
				discord.FollowUpMessageInstant(s, i.Interaction, err.Error(), instantMessageDuration)
				return
			}
			joiningChannels := make(map[string]interface{}, len(discordChannelUsers))
			for _, v := range discordChannelUsers {
				joiningChannels[v.ChannelID] = nil
			}

			content := "**[show channels]**: \nalready joined channels are __underlined__\n"
			publicChannels := "\n[public] \n"
			privateChannels := "\n[private] \n"
			for _, channel := range discordChannels {
				var addContent string
				if _, ok := joiningChannels[channel.ChannelID]; ok {
					addContent = fmt.Sprintf("`%s` __%s__\n", channel.ChannelID, channel.ChannelName)
				} else {

					addContent = fmt.Sprintf("`%s` %s\n", channel.ChannelID, channel.ChannelName)
				}
				if channel.IsPrivate {
					privateChannels += addContent
				} else {
					publicChannels += addContent
				}
			}
			content += publicChannels
			content += privateChannels

			discord.FollowUpMessageInstant(s, i.Interaction, content, instantMessageDuration)
		},

		//Command Test
		commandTest: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			discord.CreateMessageInstant(s, i.Interaction, "test by "+i.Member.User.Username, instantMessageDuration)
		},
	}
	return commandHandlers
}

func addCommands(s *discordgo.Session, guildID string, commands []*discordgo.ApplicationCommand) error {
	for _, v := range commands {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, guildID, v)
		if err != nil {
			return fmt.Errorf("cannot create '%v' command: %v", v.Name, err)
		}
	}
	return nil
}

func removeCommands(s *discordgo.Session, guildID string) error {
	registeredCommands, err := s.ApplicationCommands(s.State.User.ID, guildID)

	if err != nil {
		return fmt.Errorf("could not fetch registered commands: %v", err)
	}

	for _, v := range registeredCommands {
		err := s.ApplicationCommandDelete(s.State.User.ID, guildID, v.ID)
		if err != nil {
			return fmt.Errorf("cannot delete '%v' command: %v", v.Name, err)

		}
	}
	return nil
}

func addHandlers(s *discordgo.Session, guildID string) error {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		commandHandlers := newDiscordCommandHandler()

		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}
		default:
			log.Println(i.Data)
		}

	})
	return nil
}
