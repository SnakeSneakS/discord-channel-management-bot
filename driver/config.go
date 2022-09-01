package driver

import (
	"github.com/snakesneaks/discord-channel-management-bot/adapter/controller"
	"github.com/snakesneaks/discord-channel-management-bot/adapter/gateway"
	"github.com/snakesneaks/discord-channel-management-bot/entity"
	"github.com/snakesneaks/discord-channel-management-bot/usecase/interactor"
)

//環境変数をロードする。失敗した場合は実行終了する。
func LoadEnv() (*entity.Environment, error) {
	c := controller.NewConfig(interactor.NewConfigInputPort, gateway.NewConfigRepository)
	return c.GetEnvironment()
}
