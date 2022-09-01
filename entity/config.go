package entity

import (
	"github.com/Netflix/go-env"
)

type Environment struct {
	App struct {
		REMOVE_COMMANDS bool `env:"REMOVE_COMMANDS"`
	}

	DiscordBot struct {
		TOKEN     string `env:"DISCORD_BOT_TOKEN"`
		CLIENT_ID string `env:"DISCORD_CLIENT_ID"`
		//GUILD_ID   string `env:"DISCORD_GUILD_ID"`
		//CHANNEL_ID string `env:"DISCORD_CHANNEL_ID"`
	}

	Database struct {
		User         string `env:"MYSQL_USER"`
		Password     string `env:"MYSQL_PASSWORD"`
		RootPassword string `env:"MYSQL_ROOT_PASSWORD"`
		Host         string `env:"DB_HOST"`
		Port         string `env:"DB_PORT"`
		Database     string `env:"MYSQL_DATABASE"`
	}

	Extras env.EnvSet
}

/*
func GetEnvironment() Environment {
	var environment Environment
	es, err := env.UnmarshalFromEnviron(&environment)
	if err != nil {
		log.Fatal(err)
	}
	environment.Extras = es
	//log.Print(es)
	return environment
}
*/
