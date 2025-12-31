package config

import (
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

// LoadEnv は環境変数を読み取りConfig(グローバル変数)に代入する関数
// 引数にファイルを指定することによって、特定のenvファイルを読み取る
func LoadEnv(envfile ...string) error {
	if len(envfile) > 0 {
		if err := godotenv.Load(envfile...); err != nil {
			return err
		}
	}

	config := config{}

	if err := env.Parse(&config.App); err != nil {
		return err
	}

	if err := env.Parse(&config.Database); err != nil {
		return err
	}

	if err := env.Parse(&config.OAuth); err != nil {
		return err
	}

	if err := env.Parse(&config.Session); err != nil {
		return err
	}

	Config = &config

	return nil
}
