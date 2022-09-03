package cmd

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/snakesneaks/discord-channel-management-bot/entity"
)

func startDiscordSession(env *entity.Environment) *discordgo.Session {
	s, err := discordgo.New("Bot " + env.DiscordBot.TOKEN)
	if err != nil {
		log.Fatalf("error creating Discord session: %v", err)
	}

	s.Identify.Intents = discordgo.IntentsAll

	err = s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	for _, guild := range s.State.Guilds {
		startGuildSession(s, guild.ID)
	}

	return s
}

func endDiscordSession(s *discordgo.Session) {
	log.Println("Removing commands...")
	for _, guild := range s.State.Guilds {
		if err := removeCommands(s, guild.ID); err != nil {
			log.Println(err)
		}
	}
	s.Close()
}

func startGuildSession(s *discordgo.Session, guildID string) {
	log.Println("guildID: " + guildID)
	//s.ChannelMessageSend("707961532457156652", "i'm here!")

	//remove commands
	/*
		log.Println("Removing commands...")
		if err := removeCommands(s, guildID); err != nil {
			log.Println(err)
		}
	*/

	//add commands
	log.Println("Adding commands...")
	commands := newDiscordCommands()
	if err := addCommands(s, guildID, commands); err != nil {
		log.Println(err)
	}

	//add handlers
	if err := addHandlers(s, guildID); err != nil {
		log.Println(err)
	}
}
