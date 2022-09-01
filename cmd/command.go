package cmd

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/snakesneaks/discord-channel-management-bot/driver"
	"github.com/snakesneaks/discord-channel-management-bot/driver/discord"
)

const (
	commandHelp          = "help"
	commandCreateChannel = "create-channel"
	commandJoinChannel   = "join-channel"
	commandLeaveChannel  = "leave-channel"
	commandDeleteChannel = "delete-channel"
	commandShowChannels  = "show-channels"
	commandTest          = "test"

	optionChannelID        = "channel-id"
	optionChannelName      = "channel-name"
	optionChannelTopic     = "channel-topic"
	optionChannelIsPrivate = "channel-is-private"
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
			Name:        commandLeaveChannel,
			Description: "leave channel",
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
			Name:                     commandDeleteChannel,
			Description:              "delete channel (only manager)",
			DefaultMemberPermissions: &permissionManageServer,
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
	discordChannelDriver := driver.NewDiscordChannelDriver()
	discordChannelUserDriver := driver.NewDiscordChannelUserDriver()

	commandHandlers := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		commandHelp: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			commands := "[command]: [description]\n"
			for _, v := range newDiscordCommands() {
				commands += fmt.Sprintf("`/%s`: %s\n", v.Name, v.Description)
			}
			discord.CreateMessage(s, i.Interaction, commands)
		},
		commandCreateChannel: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			//get input
			optMap := getOptionMap(i.Interaction)
			channelName, ok := optMap[optionChannelName]
			if !ok {
				discord.CreateMessage(s, i.Interaction, fmt.Sprintf("error: option {%s} is needed.", optionChannelName))
				return
			}
			channelTopic, ok := optMap[optionChannelTopic]
			if !ok {
				discord.CreateMessage(s, i.Interaction, fmt.Sprintf("error: option {%s} is needed.", optionChannelTopic))
				return
			}
			isPrivate, ok := optMap[optionChannelIsPrivate]
			if !ok {
				discord.CreateMessage(s, i.Interaction, fmt.Sprintf("error: option {%s} is needed.", optionChannelIsPrivate))
				return
			}

			//get or create setting
			setting, found, err := discordChannelDriver.GetSetting(i.GuildID)
			if err != nil {
				discord.CreateMessage(s, i.Interaction, fmt.Sprintf("error: failed to get setting. %v", err))
			}
			if !found {
				category, err := discord.CreateCategory(s, i.Interaction, "discord-channel-management-bot-category", "this is the category discord-channel-management-bot handle.", 99999)
				if err != nil {
					discord.CreateMessage(s, i.Interaction, fmt.Sprintf("error: failed to create category. %v", err))
					return
				}

				setting, err = discordChannelDriver.CreateSetting(category.GuildID, category.ID)
				if err != nil {
					discord.CreateMessage(s, i.Interaction, fmt.Sprintf("error: failed to create setting in database. %v", err))
				}
			}

			//create channel
			c, err := discord.CreateChannel(s, i.Interaction, channelName.StringValue(), channelTopic.StringValue(), setting.ParentCategoryID, isPrivate.BoolValue())
			if err != nil {
				discord.CreateMessage(s, i.Interaction, fmt.Sprintf("error: failed to create channel. %v", err))
				return
			}
			if err := discordChannelDriver.CreateChannel(i.GuildID, c.ID, c.Name, c.Topic, i.Member.User.ID, isPrivate.BoolValue()); err != nil {
				discord.CreateMessage(s, i.Interaction, fmt.Sprintf("error: failed to create channel in database. %v", err))
				return
			}

			//join user
			if err := discord.SetMemberPermissionToChannel(s, c.ID, i.Member.User.ID, true); err != nil {
				discord.CreateMessage(s, i.Interaction, fmt.Sprintf("error: failed to join channel. %v", err))
				return
			}
			if err := discordChannelUserDriver.JoinOrLeaveChannel(i.GuildID, c.ID, i.Member.User.ID, true); err != nil {
				discord.CreateMessage(s, i.Interaction, fmt.Sprintf("error: failed to join channel in database. %v", err))
				return
			}

			//success
			discord.CreateMessage(s, i.Interaction, fmt.Sprintf("channel created!\n`id`: %s\n`name`: %s\n`topic`: %s", c.ID, c.Name, c.Topic))
		},
		commandJoinChannel: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			optMap := getOptionMap(i.Interaction)
			channelID, ok := optMap[optionChannelID]
			if !ok {
				discord.CreateMessage(s, i.Interaction, fmt.Sprintf("error: option {%s} is needed.", optionChannelID))
				return
			}

			channel, err := discordChannelDriver.GetChannel(i.GuildID, channelID.StringValue())
			if err != nil {
				discord.CreateMessage(s, i.Interaction, fmt.Sprintf("error: %v", err))
				return
			}

			if channel.IsPrivate {
				discord.CreateMessage(s, i.Interaction, fmt.Sprintf("error: channel {%s} is a private channel. you must be invited by user in that channel.", channelID.StringValue()))
				return
			}

			if err := discord.SetMemberPermissionToChannel(s, channelID.StringValue(), i.Member.User.ID, true); err != nil {
				discord.CreateMessage(s, i.Interaction, fmt.Sprintf("error: failed to join channel. %v", err))
				return
			}

			if err := discordChannelUserDriver.JoinOrLeaveChannel(i.GuildID, channelID.StringValue(), i.Member.User.ID, true); err != nil {
				discord.CreateMessage(s, i.Interaction, fmt.Sprintf("error: failed to join channel in database. %v", err))
				return
			}

			discord.CreateMessage(s, i.Interaction, fmt.Sprintf("user %s joined into %s.", i.Member.User.Mention(), channel.ChannelID))
		},
		commandLeaveChannel: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			optMap := getOptionMap(i.Interaction)
			channelID, ok := optMap[optionChannelID]
			if !ok {
				discord.CreateMessage(s, i.Interaction, fmt.Sprintf("error: option {%s} is needed.", optionChannelID))
				return
			}

			if err := discord.SetMemberPermissionToChannel(s, channelID.StringValue(), i.Member.User.ID, false); err != nil {
				discord.CreateMessage(s, i.Interaction, fmt.Sprintf("error: failed to leave channel. %v", err))
				return
			}

			if err := discordChannelUserDriver.JoinOrLeaveChannel(i.GuildID, channelID.StringValue(), i.Member.User.ID, false); err != nil {
				discord.CreateMessage(s, i.Interaction, fmt.Sprintf("error: failed to leave channel in database. %v", err))
				return
			}

			//success
			discord.CreateMessage(s, i.Interaction, fmt.Sprintf("user %s leave from %s.", i.Member.User.Mention(), channelID.StringValue()))
		},
		commandDeleteChannel: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			optMap := getOptionMap(i.Interaction)
			channelID, ok := optMap[optionChannelID]
			if !ok {
				discord.CreateMessage(s, i.Interaction, fmt.Sprintf("error: option {%s} is needed.", optionChannelID))
				return
			}

			c, err := discord.DeleteChannel(s, channelID.StringValue())
			if err != nil {
				discord.CreateMessage(s, i.Interaction, fmt.Sprintf("error: failed to delete channel. %v", err))
				return
			}

			if err := discordChannelDriver.DeleteChannel(i.GuildID, channelID.StringValue()); err != nil {
				discord.CreateMessage(s, i.Interaction, fmt.Sprintf("error: failed to delete channel from database. %v", err))
				return
			}

			if err := discordChannelUserDriver.DeleteChannelUsersOfChannel(i.GuildID, channelID.StringValue()); err != nil {
				discord.CreateMessage(s, i.Interaction, fmt.Sprintf("error: failed to delete channel-user-relation from database. %v", err))
				return
			}

			discord.CreateMessage(s, i.Interaction, fmt.Sprintf("deleted channel. Name: %s", c.Name))
		},
		commandShowChannels: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			channels, err := discordChannelDriver.GetChannels(i.GuildID)
			if err != nil {
				discord.CreateMessage(s, i.Interaction, fmt.Sprintf("error: %v", err))
			}

			content := fmt.Sprintf("channels: %d\n", len(channels))
			for _, channel := range channels {
				isPrivateText := "public"
				if channel.IsPrivate {
					isPrivateText = "private"
				}
				content += fmt.Sprintf("[%s] `%s` %s: %s\n", isPrivateText, channel.ChannelID, channel.ChannelName, channel.ChannelTopic.String)
			}

			discord.CreateMessage(s, i.Interaction, content)
		},
		commandTest: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			discord.CreateMessage(s, i.Interaction, "test by "+i.Member.User.Username)
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
