package controller

import (
	"github.com/snakesneaks/discord-channel-management-bot/entity"
	"github.com/snakesneaks/discord-channel-management-bot/usecase/port"
)

type Config struct {
	InputFactory func(port.ConfigRepository) port.ConfigInputPort
	RepoFactory  func() port.ConfigRepository
}

func NewConfig(
	inputFactory func(port.ConfigRepository) port.ConfigInputPort,
	repoFactory func() port.ConfigRepository,
) Config {
	return Config{
		InputFactory: inputFactory,
		RepoFactory:  repoFactory,
	}
}

//portを組み立ててinputPortを呼び出す
func (c *Config) GetEnvironment() (*entity.Environment, error) {
	repository := c.RepoFactory()
	inputPort := c.InputFactory(repository)
	env, err := inputPort.LoadEnvironment()
	/*
		if err != nil {
			log.Fatal(err)
		}
	*/
	return env, err
}
