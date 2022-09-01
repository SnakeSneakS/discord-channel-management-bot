package interactor

import (
	"github.com/snakesneaks/discord-channel-management-bot/entity"
	"github.com/snakesneaks/discord-channel-management-bot/usecase/port"
)

type Config struct {
	repository port.ConfigRepository
}

var _ port.ConfigInputPort = (*Config)(nil)

func NewConfigInputPort(repository port.ConfigRepository) port.ConfigInputPort {
	return &Config{
		repository: repository,
	}
}

func (c *Config) LoadEnvironment() (*entity.Environment, error) {
	return c.repository.LoadEnvironment()
}
