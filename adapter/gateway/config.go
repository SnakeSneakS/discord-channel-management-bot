package gateway

import (
	"log"

	"github.com/Netflix/go-env"
	"github.com/snakesneaks/discord-channel-management-bot/entity"
	"github.com/snakesneaks/discord-channel-management-bot/usecase/port"
)

var _ port.ConfigRepository = (*ConfigRepository)(nil)

type ConfigRepository struct {
}

func NewConfigRepository() port.ConfigRepository {
	return ConfigRepository{}
}

func (c ConfigRepository) LoadEnvironment() (*entity.Environment, error) {
	var environment entity.Environment
	es, err := env.UnmarshalFromEnviron(&environment)
	if err != nil {
		log.Fatal(err)
	}
	environment.Extras = es
	//log.Print(es)
	return &environment, err
}
