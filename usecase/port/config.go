package port

import "github.com/snakesneaks/discord-channel-management-bot/entity"

//configを入力する.
type ConfigInputPort interface {
	LoadEnvironment() (*entity.Environment, error)
}

//configを出力する.何も表示しないため必要なし
type ConfigOutputPort interface {
}

//configをロードする
type ConfigRepository interface {
	LoadEnvironment() (*entity.Environment, error)
}

//inputとoutputを繋ぐ
/*
type ConfigRepository interface {
	GetEnvironment() (*entity.Environment, error)
}
*/
